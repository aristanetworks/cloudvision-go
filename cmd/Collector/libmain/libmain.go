// Copyright (c) 2019 Arista Networks, Inc.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the COPYING file.

package libmain

import (
	"bytes"
	"context"
	"errors"
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
	"github.com/aristanetworks/cloudvision-go/provider"
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
	deviceType       *string
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

	// ingest server config
	ingestServerAddr *string

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

	// Sensor settings
	standalone              *bool // sensor running standalone or not
	sensorName              *string
	sensorHeartbeat         *time.Duration
	sensorFailureMaxBackoff *time.Duration

	// Datasource monitor settings
	logRate *float64
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
	deviceName = flag.String("deviceName", "cmd-device", "Device name")
	deviceType = flag.String("device", "", "Device type (available devices: "+deviceList()+")")
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
	gnmiServerAddr = flag.String("gnmiServerAddr", "",
		"Address of gNMI server(deprecated; use grpcServerAddr")

	// gRPC server config
	grpcServerAddr = flag.String("grpcServerAddr", "",
		"Address of gRPC server")

	// ingest server config
	ingestServerAddr = flag.String("ingestServerAddr", "",
		"Address of ingest server")

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

	protoVersion = flag.String("protoversion", "v2",
		"Protocol version to use for communicating with CV (must be v1 or v2.")

	sensorName = flag.String("sensor", "", "Sensor's identifying name. "+
		"Will not execute sensor if not set.")
	// sensor running standalone or not
	standalone = flag.Bool("standalone", true, "Run sensor in standalone mode")
	sensorHeartbeat = flag.Duration("heartbeatInterval", 10*time.Second,
		"Defines interval of sensor heartbeats")
	sensorFailureMaxBackoff = flag.Duration("failureBackoffMax", 1*time.Hour,
		"Defines maximum backoff for datasource failure retries")

	logRate = flag.Float64("logRate", 100.0, "Log rate limit (times per minute)"+
		" for datasource monitor")

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
		if *deviceType != "" {
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
	logrus.SetReportCaller(true)
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableColors: true,
		CallerPrettyfier: func(f *runtime.Frame) (function string, file string) {
			idx := strings.LastIndex(f.File, "/")
			fname := f.File
			if idx != -1 {
				// grab only file name
				fname = f.File[idx+1:]
				// grab file package if available
				idx := strings.LastIndex(f.File[:idx], "/")
				if idx != 1 {
					fname = f.File[idx+1:]
				}
			}
			return "", fmt.Sprintf("%s:%d", fname, f.Line)
		},
	})
}

// newCVClient returns a CVClient object for a given device.
func newCVClient(gc gnmi.GNMIClient, info *device.Info) cvclient.CVClient {
	var isManager bool
	if _, ok := info.Device.(device.Manager); ok {
		isManager = true
	}
	if protoVersion != nil && *protoVersion == "v2" {
		return v2client.NewV2Client(gc, info)
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

	if *gnmiServerAddr == "" {
		gnmiCfg.Addr = *grpcServerAddr
	}

	if *ingestServerAddr == "" {
		ingestServerAddr = grpcServerAddr
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

	var grpcConn *grpc.ClientConn
	if *grpcServerAddr != "" {
		if *authInfo != noAuth {
			logrus.Infof("Connecting to gRPC server %+v", authInfo)
			var err error
			grpcConn, err = agrpc.DialWithAuth(ctx, *grpcServerAddr, authInfo)
			if err != nil {
				logrus.Fatalf("DialWithAuth error: %v", err)
			}
			logrus.Info("Connected")
			opts = append(opts,
				device.WithGRPCConn(grpcConn),
			)
		}
		opts = append(opts,
			device.WithGRPCServerAddr(*ingestServerAddr),
			device.WithGRPCConnector(sc.Connector),
			device.WithStandaloneStatus(*standalone),
		)
	}

	logrus.Infof("Connecting to gNMI server %+v", gnmiCfg)
	conn, err := agnmi.DialContextConn(ctx, gnmiCfg)
	if err != nil {
		logrus.Fatal(err)
	}
	defer conn.Close()
	gnmiClient := gnmi.NewGNMIClient(conn)
	waitForGNMIConnectivity(gnmiClient)
	logrus.Info("Connected to gNMI service")

	opts = append(opts, device.WithGNMIClient(gnmiClient))

	// Create inventory.
	inventory := device.NewInventoryWithOptions(ctx, opts...)

	group, ctx := errgroup.WithContext(ctx)

	var cmdDevice *device.Config
	if *deviceType != "" {
		// Keep the group alive when running with cmd line config
		cmdDevice = &device.Config{
			Name:     *deviceName,
			Device:   *deviceType,
			NoStream: *noStream,
			Options:  deviceOptions,
		}
		group.Go(func() error {
			<-ctx.Done()
			return nil
		})
	}

	// Watch config for old collector flow
	configs, err := createDeviceConfigs(cmdDevice, *deviceConfigFile)
	if err != nil {
		logrus.Fatal(err)
	}

	configCh := make(chan *device.Config, len(configs))
	for _, config := range configs {
		configCh <- config
	}

	group.Go(func() error {
		defer close(configCh)
		return watchConfig(ctx, cmdDevice, *deviceConfigFile, configCh, time.Second)
	})

	// Start sensor state machine.
	if len(*sensorName) > 0 {
		opts := []device.SensorOption{
			device.WithSensorHeartbeatInterval(*sensorHeartbeat),
			device.WithSensorGNMIClient(gnmiClient),
			device.WithSensorClientFactory(newCVClient),
			device.WithSensorConfigChan(configCh),
		}
		if sc.CredResolverCreator != nil {
			resolver, err := sc.CredResolverCreator(conn)
			if err != nil {
				logrus.Fatal(err)
			}
			opts = append(opts, device.WithSensorCredentialResolver(resolver))
		}
		if *grpcServerAddr != "" {
			opts = append(opts,
				device.WithSensorFailureRetryBackoffMax(*sensorFailureMaxBackoff),
				device.WithSensorGRPCConn(grpcConn),
				device.WithSensorConnector(sc.Connector),
				device.WithSensorConnectorAddress(*ingestServerAddr),
				device.WithSensorStandaloneStatus(*standalone))
		}
		group.Go(func() error {
			logrus.Infof("Starting sensor %v", *sensorName)

			//added retry logic in case of error due to GNMI service stop/restart
			backoffTimer := provider.NewBackoffTimer()
			for {
				select {
				case <-ctx.Done():
					return ctx.Err()
				case <-backoffTimer.Wait():
					waitForGNMIConnectivity(gnmiClient)
					sensor := device.NewSensor(*sensorName, *logRate, opts...)
					err := sensor.Run(ctx)
					// Sensor failed, schedule retry with backoff.
					// This is done before logging the error so we can log a precise retry delay.
					curBackoff := backoffTimer.Backoff()
					if !errors.Is(err, context.Canceled) {
						logrus.Infof("sensor run failed, retrying in %v. Err: %v",
							curBackoff, err)
					}
				}
			}
		})
	} else {
		// Push configs from the file to the inventory.
		group.Go(func() error {
			for {
				select {
				case <-ctx.Done():
					return nil
				case cfg, ok := <-configCh:
					if !ok {
						return nil
					}
					noop := NewBaseMonitor(logrus.WithField("Device", cfg.Name))
					info, err := device.NewDeviceInfo(ctx, cfg, noop)
					if err != nil {
						logrus.Infof("Error in device.NewDeviceInfo: %v. Dropping device %s.",
							err, cfg.Name)
						continue
					}

					if cfg.IsDeleted() {
						if err := inventory.Delete(info.ID); err != nil {
							logrus.Infof("Error in inventory.Add: %v. Dropping device %s.",
								err, info.ID)
						}
						continue
					}

					if err := inventory.Add(info); err != nil {
						logrus.Infof("Error in inventory.Add: %v. Dropping device %s.",
							err, info.ID)
						continue
					}
				}
			}
		})
	}

	if *grpcAddr != "" {
		grpcServer, listener, err := newGRPCServer(*grpcAddr, inventory)
		if err != nil {
			logrus.Fatalf("Failed to start gRPC server: %v", err)
		}
		group.Go(func() error { return grpcServer.Serve(listener) })
		group.Go(func() error {
			// A sensor failure should cause the program to end, so stop the grpc server.
			<-ctx.Done()
			logrus.Info("context canceled, closing grpc server...")
			grpcServer.Stop()
			_ = listener.Close()
			return nil
		})
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

func createDeviceConfigs(cmdDevice *device.Config,
	deviceConfigFile string) ([]*device.Config, error) {
	configs := []*device.Config{}
	if cmdDevice != nil {
		copy := *cmdDevice
		configs = append(configs, &copy)
	}

	if deviceConfigFile != "" {
		readConfigs, err := device.ReadConfigs(deviceConfigFile)
		if err != nil {
			return nil, err
		}
		configs = append(configs, readConfigs...)
	}

	// Make sure all configs have a name
	for i, config := range configs {
		if len(config.Name) == 0 { // force a name if not set
			config.Name = fmt.Sprintf("auto-datasource-%03d", i)
		}
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
	oh, err := device.OptionHelp(*deviceType)
	if err != nil {
		return fmt.Errorf("addHelp: %v", err)
	}

	var formattedOptions string
	if len(oh) > 0 {
		b := new(bytes.Buffer)
		aflag.FormatOptions(b, "Help options for device type '"+*deviceType+"':", oh)
		formattedOptions = b.String()
	}

	aflag.AddHelp("", formattedOptions)
	return nil
}

func validateConfig() {
	if *deviceConfigFile != "" && *deviceType != "" {
		logrus.Fatal("-config and -device should not be both specified.")
	}

	if *noStream && *deviceType == "" {
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

	if !*standalone && *ingestServerAddr == "" {
		logrus.Fatal("-ingestServerAddr must be specified in case of sensor not running standalone")
	}
}

func watchConfig(ctx context.Context, cmdDevice *device.Config, configPath string,
	pushToCh chan *device.Config, redeployDelay time.Duration) error {
	if configPath == "" {
		return nil
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()

	logrus.Infof("Monitoring %s config file for changes", configPath)

	configDir, _ := filepath.Split(configPath)
	group, ctx := errgroup.WithContext(ctx)
	existingConfigs := map[string]struct{}{}

	// Timer to delay config refresh on file events to prevent events firing in quick succession.
	timer := time.NewTimer(0)
	if !timer.Stop() { // Clear timer for now
		<-timer.C
	}
	defer timer.Stop()

	group.Go(func() error {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok { // 'Events' channel is closed
					return nil
				}
				if event.Name != configPath { // ignore events on other files
					continue
				}
				logrus.Debugf("Watcher event: %v", event)

				// File changes can come with many events at once, so we give it a second to settle
				// before trying to read the file, or we may read an empty file or fail to read.
				// Reading no data is undesirable as it would cause a config flap.
				_ = timer.Reset(redeployDelay)
			case <-timer.C:
				configs, err := createDeviceConfigs(cmdDevice, configPath)
				if err != nil {
					logrus.Errorf("Error creating device configs from watched config: %v", err)
					continue
				}
				logrus.Debugf("Watcher read %d configs", len(configs))

				newSet := map[string]struct{}{}
				for _, config := range configs {
					select {
					case pushToCh <- config:
					case <-ctx.Done():
						return nil
					}
					newSet[config.Name] = struct{}{}
				}
				for name := range existingConfigs {
					if _, ok := newSet[name]; !ok {
						select {
						case pushToCh <- device.NewDeletedConfig(name):
						case <-ctx.Done():
							return nil
						}
					}
				}
				existingConfigs = newSet
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
