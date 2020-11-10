/*
 * Copyright © 2020 Matthew Ellison <seglberg+oss@gmail.com>
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

// SCANNER
// 	Scanner connects to the given target IP and port and creates a TCP connection.
// 	It then waits for a response from server (or times out in the event the server doesn't send anything, e.g. a web server).
//
//	If the server responds, the message is decoded and is attempted to be parsed as a well known type (for example a MySQL packet).
//  The results are then printed as JSON to the STDOUT.

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/seglberg/protoscan/pkg/mysql"
)

var args = struct {
	target      **net.TCPAddr
	initTimeout *time.Duration
	readTimeout *time.Duration
}{
	kingpin.Arg("target", "Target host and port to scan").
		Default("localhost:3306").
		TCP(),

	kingpin.Flag("init-timeout", "Maximum amount of time to wait for a connection to be made").
		Default("10s").
		Duration(),

	kingpin.Flag("read-timeout", "Maximum amount of time to wait for server to respond once a connection is made. Set to 0 to wait indefinitely.").
		Default("5s").
		Duration(),
}

func init() {
	kingpin.Parse()
}

func main() {
	ctx := context.Background()

	// (1) Make the Initial TCP Connection

	target := (*args.target).String()

	dialer := &net.Dialer{
		Timeout: *args.initTimeout,
	}

	conn, err := dialer.DialContext(ctx, "tcp", target)
	if err != nil {
		log.Fatal(err)
	}
	// Best effort close of connection.
	defer func() {
		_ = conn.Close()
	}()

	// (2) Read First Packet

	err = conn.SetReadDeadline(time.Now().Add(*args.readTimeout))
	if err != nil {
		log.Fatal(err)
	}

	when := time.Now()

	packet, err := mysql.ReadPacket(conn)
	if err != nil {
		if errors.Is(err, os.ErrDeadlineExceeded) {
			log.Fatal("timed out waiting for server")
		} else {
			log.Fatal(err)
		}
	}

	// (3) Decode Packet Payload

	hs, err := mysql.DecodeHandshake(packet.Payload)
	if err != nil {
		log.Fatal(err)
	}

	// (4) Print Results

	result := struct {
		Target       string          `json:"target"`
		When         time.Time       `json:"when"`
		ProtoVersion int             `json:"proto_version"`
		Handshake    mysql.Handshake `json:"handshake"`
	}{
		Target:       target,
		When:         when,
		ProtoVersion: int(hs.GetProtoVersion()),
		Handshake:    hs,
	}

	serialized, _ := json.MarshalIndent(result, "", "  ")
	fmt.Println(string(serialized))
}
