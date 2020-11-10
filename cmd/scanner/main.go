/*
   --------------------------------------------------------------------------
	LICENSE HERE
   --------------------------------------------------------------------------
*/

package main

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

	result := struct {
		Target       string          `json:"target"`
		When         time.Time       `json:"when"`
		ProtoVersion int             `json:"proto_version"`
		Handshake    mysql.Handshake `json:"handshake"`
	}{
		Target:       target,
		When:         time.Now(),
		ProtoVersion: int(hs.GetProtoVersion()),
		Handshake:    hs,
	}

	serialized, _ := json.MarshalIndent(result, "", "  ")
	fmt.Println(string(serialized))
}
