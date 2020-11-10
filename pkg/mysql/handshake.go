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

// Package mysql provides facilities for decoding and inspecting MySQL wire protocol information.
package mysql

import (
	"encoding/binary"
	"fmt"
)

var ErrHandshakeDecode = fmt.Errorf("handshake decode")
var ErrHandshakeTruncated = fmt.Errorf("%w: truncated payload or not a mysql handshake", ErrHandshakeDecode)

type Handshake interface {
	// GetProtoVersion returns the MySQL Protocol Version.
	GetProtoVersion() uint8
}

// HandshakeV10 represents the MySQL initial handshake packet for protocol version 10.
// This payload is sent from the server to the client at the start of the conversation.
//
// See https://dev.mysql.com/doc/internals/en/connection-phase-packets.html#packet-Protocol::Handshake
type HandshakeV10 struct {
	// ServerVersion contains the human readable server version.
	ServerVersion string `json:"server_version"`

	// ThreadID is the connection ID.
	ThreadID uint32 `json:"thread_id"`

	// AuthPluginData contains scramble data used by authentication plugins.
	AuthPluginData []byte `json:"auth_plugin_data"`

	// CharacterSet contains the server's charset id.
	CharacterSet uint8 `json:"character_set"`

	// CapabilityFlags is a composite flag field used by the client and server
	// to communicate supported functions and features.
	CapabilityFlags Capability `json:"capability_flags"`

	// ServerStatusFlags is a composite flag field used to communicate the current
	// status of the server.
	ServerStatusFlags ServerStatus `json:"server_status_flags"`

	// AuthPluginName is the name (if any) of the auth plugin which the authentication scramble data belongs to.
	AuthPluginName string `json:"auth_plugin_name,omitempty"`
}

// GetProtoVersion returns the MySQL Protocol Version this Handshake implements.
// Always 10.
func (*HandshakeV10) GetProtoVersion() uint8 {
	return 10
}

// DecodeHandshake attempts to read and decode the given series of bytes as a MySQL Handshake payload.
// If decoding is successful, a Handshake representing the actual underlying handshake message will be returned.
//
// See https://dev.mysql.com/doc/internals/en/connection-phase.html
func DecodeHandshake(payload []byte) (Handshake, error) {
	// (1) Parse Handshake Version
	//		The first byte contains the protocol wire version.

	sub, pos, err := readBuffer(payload, 0, 1)
	if err != nil {
		return nil, ErrHandshakeTruncated
	}

	version := sub[0]

	// (2) Continue Decoding the Remaining Payload
	//		according to the version field.

	switch version {
	case 10:
		hs := &HandshakeV10{}

		// Variable: Server Version (NULL-Terminated)

		sub, pos, err = readBuffer(payload, pos, nullTermStringLen(payload[pos:]))
		if err != nil {
			return nil, ErrHandshakeTruncated
		}
		hs.ServerVersion = string(sub)
		pos++ // Add additional offset for null byte

		// 4 Bytes: Thread ID

		sub, pos, err = readBuffer(payload, pos, 4)
		if err != nil {
			return nil, ErrHandshakeTruncated
		}
		hs.ThreadID = binary.LittleEndian.Uint32(sub)

		// 8 Bytes: Part 1 of Auth Scramble

		sub, pos, err = readBuffer(payload, pos, 8)
		if err != nil {
			return nil, ErrHandshakeTruncated
		}
		hs.AuthPluginData = sub
		pos++ // Add additional offset for null byte

		// 2 Bytes: Lower Bytes of Capability Flags

		sub, pos, err = readBuffer(payload, pos, 2)
		if err != nil {
			return nil, ErrHandshakeTruncated
		}
		lowerCap := sub

		// 1 Byte: Character Set ID

		sub, pos, err = readBuffer(payload, pos, 1)
		if err != nil {
			return nil, ErrHandshakeTruncated
		}
		hs.CharacterSet = sub[0]

		// 2 Bytes: Server Status Flags

		sub, pos, err = readBuffer(payload, pos, 2)
		if err != nil {
			return nil, ErrHandshakeTruncated
		}
		hs.ServerStatusFlags = ServerStatus(binary.LittleEndian.Uint16(sub))

		// 2 Bytes: Upper Byes of Capability Flags

		sub, pos, err = readBuffer(payload, pos, 2)
		if err != nil {
			return nil, ErrHandshakeTruncated
		}
		upperCap := sub
		hs.CapabilityFlags = Capability(binary.LittleEndian.Uint32(append(upperCap, lowerCap...)))

		// 1 Byte: Plugin Data Length

		sub, pos, err = readBuffer(payload, pos, 1)
		if err != nil {
			return nil, ErrHandshakeTruncated
		}
		pluginLen := sub[0]

		// 10 Bytes: Reserved

		pos += 10

		// Variable: Remaining Plugin Data Scramble

		offset := 0
		if pluginLen > 0 {
			offset = int(pluginLen) - 8
			if offset > 13 {
				offset = 13
			}
		}

		sub, pos, err = readBuffer(payload, pos, offset)
		if err != nil {
			return nil, ErrHandshakeTruncated
		}
		hs.AuthPluginData = append(hs.AuthPluginData, sub...)

		// Variable: Auth Plugin Name (NULL-Terminated)
		//	Only present if CLIENT_PLUGIN_AUTH capability is set.

		if hs.CapabilityFlags.Has(CapabilityPluginAuth) {
			sub, pos, err = readBuffer(payload, pos, nullTermStringLen(payload[pos:]))
			if err != nil {
				return nil, ErrHandshakeTruncated
			}
			hs.AuthPluginName = string(sub)
		}

		fmt.Println(pos)

		return hs, nil

	default:
		return nil, fmt.Errorf("%w: unsupported protocol version or not a mysql handshake", ErrHandshakeDecode)
	}
}

// Parses a series of bytes for a null-terminated string (C-style string) and returns
// its ending index in the bye slice.
// If no null byte is found, 0 is returned.
func nullTermStringLen(n []byte) int {
	for i := 0; i < len(n); i++ {
		if n[i] == 0 {
			return i
		}
	}
	return 0
}

// Reads the buffer at the given position up to the offset.
// The resulting sub-slice of bytes is returned, along with the new cursor position.
// If the given position + offset extend past the slice, out of bounds, an errors is returned.
func readBuffer(b []byte, pos, offset int) ([]byte, int, error) {
	if pos+offset > len(b) {
		return nil, 0, fmt.Errorf("out of bounds")
	}
	return b[pos : pos+offset], pos + offset, nil
}
