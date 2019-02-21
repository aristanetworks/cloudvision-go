// Copyright (c) 2019 Arista Networks, Inc.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the COPYING file.

package main

import (
	"compress/gzip"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
)

var (
	community = flag.String("c", "", "SNMP community string")
	dev       = flag.String("d", "", "Device hostname/IP")
	dumpfile  = flag.String("o", "", "Name of file to write SNMP dump to")
	oids      = oidFlags{}
	polls     = flag.Int("n", 2, "Number of polls to perform")
)

type oidFlags []string

func (o *oidFlags) String() string {
	return strings.Join(*o, "")
}

func (o *oidFlags) Set(v string) error {
	*o = append(*o, v)
	return nil
}

func init() {
	flag.Var(&oids, "oid", "OID to walk - may be repeated to specify multiple")
}

func snmpWalkCmd(oid string) (string, []string) {
	cmd := []string{"-O", "ne", "-Cc"}
	if *community != "" {
		cmd = append(cmd, "-c", *community)
	}
	cmd = append(cmd, *dev, oid)
	return "snmpbulkwalk", cmd
}

func snmpWalk(f io.Writer) {
	for _, o := range oids {
		c, args := snmpWalkCmd(o)
		cmd := exec.Command(c, args...)
		cmd.Stdout = f
		fmt.Printf("Walking OID '%s'...\n", o)
		if err := cmd.Run(); err != nil {
			log.Fatalf("Walk command failed: %v", err)
		}
	}
}

func main() {
	flag.Parse()

	if *dev == "" {
		fmt.Println("-d must be specified")
		os.Exit(1)
	}
	if *dumpfile == "" {
		fmt.Println("-o must be specified")
		os.Exit(1)
	}
	if len(oids) == 0 {
		oids = []string{"."}
	}

	f, err := os.Create(*dumpfile)
	if err != nil {
		log.Fatalf("Failed to open dumpfile: %v", err)
	}
	defer f.Close()

	gf := gzip.NewWriter(f)
	gf.Header.Name = *dumpfile + ".gz"
	defer gf.Close()

	for i := 0; i < *polls; i++ {
		snmpWalk(gf)
	}
}
