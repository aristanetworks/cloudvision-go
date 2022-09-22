// Copyright (c) 2019 Arista Networks, Inc.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the COPYING file.

package libmain

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"net"
	"net/http"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/aristanetworks/cloudvision-go/device"
	"github.com/aristanetworks/cloudvision-go/device/cvclient"
	v1client "github.com/aristanetworks/cloudvision-go/device/cvclient/v1"
	v2client "github.com/aristanetworks/cloudvision-go/device/cvclient/v2"
	_ "github.com/aristanetworks/cloudvision-go/device/devices" // import all registered devices
	"github.com/aristanetworks/cloudvision-go/device/gen"
	agrpc "github.com/aristanetworks/cloudvision-go/grpc"
	"github.com/aristanetworks/cloudvision-go/log"
	"github.com/aristanetworks/cloudvision-go/version"

	"github.com/aristanetworks/fsnotify"
	aflag "github.com/aristanetworks/goarista/flag"
	agnmi "github.com/aristanetworks/goarista/gnmi"
	"github.com/openconfig/gnmi/proto/gnmi"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

var (
	v    *bool
	help *bool

	// Logging config
	logLevel *string
	logDir   *string

	// Device config
	deviceName       *string
	deviceOptions    = aflag.Map{}
	deviceConfigFile *string
	noStream         *bool

	// MockCollector config
	mock        *bool
	mockFeature = aflag.Map{}
	mockTimeout *time.Duration

	// Dump Collector config
	dump        *bool
	dumpFile    *string
	dumpTimeout *time.Duration

	// gNMI server config
	gnmiServerAddr *string

	// gRPC server config
	grpcServerAddr *string

	// local gRPC server config for inventory service
	grpcAddr *string

	// local http monitor server addr
	monitorAddr *string

	// Auth config
	caFile   *string
	certFile *string
	keyFile  *string
	tlsFlag  *bool
	authInfo = agrpc.AuthFlag()

	// protocol version
	protoVersion *string
)

// Main is the "real" main.
func Main(sc device.SensorConfig) {
	v = flag.Bool("version", false, "Print the version number")
	help = flag.Bool("help", false, "Print program options")

	// Logging config
	logLevel = flag.String("logLevel", "info", "Log level verbosity "+
		"(available levels: trace, debug, info, warning, error, fatal, panic)")
	logDir = flag.String("logDir", "", "If specified, one log file per device will be created"+
		" and written in the directory. Otherwise logs will be written to stderr.")

	// Device config
	deviceName = flag.String("device", "",
		"Device type (available devices: "+deviceList()+")")
	deviceOptions = aflag.Map{}
	deviceConfigFile = flag.String("configFile", "", "Path to the config file for devices")
	noStream = flag.Bool("nostream", false,
		"If set, updates aren't streamed for specified device")

	// MockCollector config
	mock = flag.Bool("mock", false, "Run Collector in mock mode")
	mockFeature = aflag.Map{}
	mockTimeout = flag.Duration("mockTimeout", 60*time.Second,
		"Timeout for checking notifications in mock mode")

	// Dump Collector config
	dump = flag.Bool("dump", false, "Run Collector in dump mode")
	dumpFile = flag.String("dumpFile", "", "Path to output file used to dump gNMI SetRequests")
	dumpTimeout = flag.Duration("dumpTimeout", 20*time.Second,
		"Timeout for dumping gNMI SetRequests")

	// gNMI server config
	gnmiServerAddr = flag.String("gnmiServerAddr", "localhost:6030",
		"Address of gNMI server")

	// gRPC server config
	grpcServerAddr = flag.String("grpcServerAddr", "",
		"Address of gRPC server")

	// local gRPC server config for inventory service
	grpcAddr = flag.String("grpcAddr", "",
		"Collector gRPC server address (if unspecified, server will not run)")

	// local http monitor server addr
	monitorAddr = flag.String("monitorAddr", "",
		"The address for the monitor server. If empty, monitor is not started. "+
			"Example: 0.0.0.0:0 or localhost:6060. Port 0 will select one automatically.")

	// Auth config
	caFile = flag.String("cafile", "", "Path to CA file")
	certFile = flag.String("certfile", "", "Path to client TLS certificate file")
	keyFile = flag.String("keyfile", "", "Path to client TLS private key file")
	tlsFlag = flag.Bool("tls", false, "Enable TLS")

	protoVersion = flag.String("protoversion", "v1",
		"Protocol version to use for communicating with CV (must be v1 or v2.")

	flag.Var(mockFeature, "mockFeature",
		"<feature>=<path> option for mock mode, where <path> is a path that, "+
			"if present in the Collector output, signifies that the target device supports "+
			"the feature described in <feature>")
	flag.Var(deviceOptions, "deviceoption", "<key>=<value> option for the Device. "+
		"May be repeated to set multiple Device options.")
	flag.BoolVar(help, "h", false, "Print program options")

	flag.Parse()

	// Print version.
	if *v {
		vs := []string{version.CollectorVersion, runtime.Version()}

		// If the package version is different than the Collector version then we should
		// print the package version as well.
		if version.CollectorVersion != version.Version {
			vs = append(vs, version.Version)
		}
		fmt.Println(strings.Join(vs, " "))
		return
	}

	// Print help, including device-specific help,
	// if requested.
	if *help {
		if *deviceName != "" {
			if err := addHelp(); err != nil {
				logrus.Fatal(err)
			}
		}
		flag.Usage()
		return
	}

	// We're running for real at this point. Check that the config
	// is sane.
	validateConfig()

	initLogging()

	runMonitor()

	if *mock {
		runMock(context.Background())
		return
	}
	if *dump {
		runDump(context.Background())
		return
	}
	runMain(context.Background(), sc)
}

func runMonitor() {
	if len(*monitorAddr) <= 0 {
		return
	}

	monitorListener, err := net.Listen("tcp", *monitorAddr)
	if err != nil {
		logrus.Fatalf("Failed to listen on monitor address %v: %v", *monitorAddr, err)
	}
	logrus.Infof("Monitor listening on %s", monitorListener.Addr())
	go func() {
		err := http.Serve(monitorListener, nil)
		if err != nil {
			logrus.Infof("Monitor failed, no longer serving monitor endpoints: %v", err)
		}
	}()
}

func initLogging() {
	log.SetLogDir(*logDir)
	if lv, err := logrus.ParseLevel(*logLevel); err != nil {
		logrus.Fatal(err)
	} else {
		logrus.SetLevel(lv)
	}
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableColors: true,
	})
}

// newCVClient returns a CVClient object for a given device.
func newCVClient(gc gnmi.GNMIClient, info *device.Info) cvclient.CVClient {
	var isManager bool
	if _, ok := info.Device.(device.Manager); ok {
		isManager = true
	}
	if protoVersion != nil && *protoVersion == "v2" {
		return v2client.NewV2Client(gc, info.Device)
	}
	return v1client.NewV1Client(gc, info.ID, isManager)
}

func runMain(ctx context.Context, sc device.SensorConfig) {
	gnmiCfg := &agnmi.Config{
		Addr:        *gnmiServerAddr,
		CAFile:      *caFile,
		CertFile:    *certFile,
		KeyFile:     *keyFile,
		TLS:         *tlsFlag,
		DialOptions: []grpc.DialOption{grpc.WithBlock()},
	}

	var noAuth agrpc.Auth
	opts := []device.InventoryOption{
		device.WithClientFactory(newCVClient),
	}
	if *authInfo != noAuth {
		gnmiCfg.TLS = true
		gnmiCfg.CAFile = authInfo.CAFile()
		// XXX: agnmi.Dial does not use system cert pool if this is "",
		// rather it disables the certificate validation option - this
		// should be fixed.
		clientCreds, err := authInfo.ClientCredentials()
		if err != nil {
			logrus.Fatal(err)
		}
		gnmiCfg.DialOptions = append(gnmiCfg.DialOptions, clientCreds...)
	}

	if *grpcServerAddr != "" {
		logrus.Infof("Connecting to gRPC server %+v", authInfo)
		conn, err := agrpc.DialWithAuth(ctx, *grpcServerAddr, authInfo)
		if err != nil {
			logrus.Fatalf("DialWithAuth error: %v", err)
		}
		logrus.Info("Connected")

		opts = append(opts,
			device.WithGRPCConn(conn),
			device.WithGRPCServerAddr(*grpcServerAddr),
			device.WithGRPCConnector(sc.Connector),
		)
	}

	logrus.Infof("Connecting to gNMI server %+v", gnmiCfg)
	gnmiClient, err := agnmi.Dial(gnmiCfg)
	if err != nil {
		logrus.Fatal(err)
	}
	waitForGNMIConnectivity(gnmiClient)
	logrus.Info("Connected to gNMI service")

	opts = append(opts, device.WithGNMIClient(gnmiClient))

	// Create inventory.
	inventory := device.NewInventoryWithOptions(ctx, opts...)

	group, ctx := errgroup.WithContext(ctx)
	configs, err := createDeviceConfigs()
	if err != nil {
		logrus.Fatal(err)
	}
	for _, config := range configs {
		info, err := device.NewDeviceInfo(config)
		if err != nil {
			logrus.Infof("Error in device.NewDeviceInfo(): %v. Dropping device.", err)
			continue
		}
		err = inventory.Add(info)
		if err != nil {
			logrus.Infof("Error in inventory.Add(): %v. Dropping device %s.",
				err, info.ID)
			continue
		}
		group.Go(func() error {
			<-ctx.Done()
			return nil
		})
	}
	group.Go(func() error {
		return watchConfig(*deviceConfigFile, inventory)
	})

	if *grpcAddr != "" {
		grpcServer, listener, err := newGRPCServer(*grpcAddr, inventory)
		if err != nil {
			logrus.Fatalf("Failed to start gRPC server: %v", err)
		}
		group.Go(func() error { return grpcServer.Serve(listener) })
	}

	logrus.Info("Collector is running")
	if err := group.Wait(); err != nil {
		logrus.Fatalf("group returned with error: %v", err)
	}
	logrus.Infof("Collector is finished")
}

func waitForGNMIConnectivity(gnmiClient gnmi.GNMIClient) {
	logEvery := 5
	maxWait := 30 * time.Second
	for retry := 0; ; {
		_, err := gnmiClient.Capabilities(context.Background(), &gnmi.CapabilityRequest{},
			grpc.WaitForReady(true))
		if err == nil {
			return
		}
		retry++
		if retry%logEvery == 0 {
			logrus.Errorf("Unable to reach gNMI service: %v. retrying (attempt %d)...", err, retry)
			// reduce frequency of logs to once every 20 retries
			logEvery = 20
		}

		wait := time.Second * time.Duration(retry)
		if wait > maxWait {
			wait = maxWait
		}
		time.Sleep(wait)
	}
}

func createDeviceConfigs() ([]*device.Config, error) {
	configs := []*device.Config{}
	if *deviceName != "" {
		configs = append(configs, &device.Config{
			Device:   *deviceName,
			NoStream: *noStream,
			Options:  deviceOptions,
		})
	}

	if *deviceConfigFile != "" {
		readConfigs, err := device.ReadConfigs(*deviceConfigFile)
		if err != nil {
			return nil, err
		}
		configs = append(configs, readConfigs...)
	}
	return configs, nil
}

// Return a formatted list of available devices.
func deviceList() string {
	dl := device.Registered()
	if len(dl) > 0 {
		return strings.Join(dl, ", ")
	}
	return "none"
}

func addHelp() error {
	oh, err := device.OptionHelp(*deviceName)
	if err != nil {
		return fmt.Errorf("addHelp: %v", err)
	}

	var formattedOptions string
	if len(oh) > 0 {
		b := new(bytes.Buffer)
		aflag.FormatOptions(b, "Help options for device '"+*deviceName+"':", oh)
		formattedOptions = b.String()
	}

	aflag.AddHelp("", formattedOptions)
	return nil
}

func validateConfig() {
	if *deviceConfigFile != "" && *deviceName != "" {
		logrus.Fatal("-config and -device should not be both specified.")
	}

	if *noStream && *deviceName == "" {
		logrus.Fatal("device name must be specified if -nostream is set true")
	}

	if !*mock && len(mockFeature) > 0 {
		logrus.Fatal("-mockFeature is only valid in mock mode")
	}

	if *mock && *dump {
		logrus.Fatal("-mock and -dump should not be both specified")
	}

	if *dump && *dumpFile == "" {
		logrus.Fatal("-dumpFile must be specified in dump mode")
	}

	if *protoVersion != "v1" && *protoVersion != "v2" {
		logrus.Fatal("Protocol version must be either 'v1' or 'v2'")
	}
}

func watchConfig(configPath string, inventory device.Inventory) error {
	if configPath == "" {
		return nil
	}
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()
	if err != nil {
		return err
	}
	configDir, _ := filepath.Split(configPath)
	group, ctx := errgroup.WithContext(context.Background())
	group.Go(func() error {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok { // 'Events' channel is closed
					return nil
				}
				if event.Name == configPath {
					configs, err := createDeviceConfigs()
					if err != nil {
						logrus.Errorf("Error creating device configs from watched config: %v", err)
						continue
					}
					for _, config := range configs {
						info, err := device.NewDeviceInfo(config)
						if err != nil {
							logrus.Errorf(
								"Error creating device info from device config: %v", err)
							continue
						}
						err = inventory.Add(info)
						if err != nil {
							logrus.Errorf("Error adding device to inventory: %v", err)
							continue
						}
					}
				}
			case err, ok := <-watcher.Errors:
				if ok { // 'Errors' channel is not closed
					return fmt.Errorf("Watcher error: %v", err)
				}
				return nil
			case <-ctx.Done():
				return nil
			}
		}
	})
	// we have to watch the entire directory to pick up changes to symlinks
	if err = watcher.Add(configDir); err != nil {
		return err
	}
	return group.Wait()
}

func newGRPCServer(address string,
	inventory device.Inventory) (*grpc.Server, net.Listener, error) {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return nil, nil, err
	}
	grpcServer := grpc.NewServer()
	gen.RegisterDeviceInventoryServer(grpcServer, device.NewInventoryService(inventory))
	reflection.Register(grpcServer)
	healthpb.RegisterHealthServer(grpcServer, health.NewServer())
	return grpcServer, listener, nil
}
