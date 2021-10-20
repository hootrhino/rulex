// Copyright 2012 The GoSNMP Authors. All rights reserved.  Use of this
// source code is governed by a BSD-style license that can be found in the
// LICENSE file.

package test

import (
	"fmt"
	"log"
	"testing"

	g "github.com/gosnmp/gosnmp"
)
//
// https://www.alvestrand.no/objectid/top.html
//
func TestSnmp(t *testing.T) {

	// Default is a pointer to a GoSNMP struct that contains sensible defaults
	// eg port 161, community public, etc
	g.Default.Target = "127.0.0.1"
	g.Default.Community = "public"
	err := g.Default.Connect()
	if err != nil {
		log.Fatalf("Connect() err: %v", err)
	}
	defer g.Default.Conn.Close()

	oids := []string{".1.3.6.1.2.1.1.1.0"}
	result, err2 := g.Default.Get(oids)
	if err2 != nil {
		log.Fatalf("Get() err: %v", err2)
	}

	for i, variable := range result.Variables {
		fmt.Printf("%d: oid: %s ", i, variable.Name)

		switch variable.Type {
		case g.OctetString:
			fmt.Printf("string: %s\n", string(variable.Value.([]byte)))
		default:

			fmt.Printf("number: %d\n", g.ToBigInt(variable.Value))
		}
	}
}
