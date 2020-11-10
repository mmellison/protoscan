/*
   --------------------------------------------------------------------------
   Copyright (c) Arroyo Networks - 2020 - All Rights Reserved
   Proprietary and Confidential

   Unauthorized copying of this file, via any medium, is strictly prohibited.
   --------------------------------------------------------------------------
*/

package mysql

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
)

var ErrPacketDecode = fmt.Errorf("packet decode")

// Packet represents the basic MySQL packet.
// See https://dev.mysql.com/doc/internals/en/mysql-packet.html
type Packet struct {
	// SequenceID is the packet's sequence ID.
	SequenceID uint8

	// Payload of the packet.
	Payload []byte
}

// ReadPacket attempts to read a MySQL packet from the given reader and produce
// the packet's payload.
func ReadPacket(r io.Reader) (*Packet, error) {

	// Packet:
	//	3 Bytes: Payload Length
	//	1 Byte: Sequence ID
	//	PAYLOAD

	p := &Packet{}

	// Length

	buf := make([]byte, 3)
	n, err := r.Read(buf)
	if err != nil {
		if errors.Is(err, os.ErrDeadlineExceeded) {
			return nil, err
		}
		return nil, fmt.Errorf("%w: %v", ErrPacketDecode, err)
	}
	if n != 3 {
		return nil, fmt.Errorf("%w: truncated header", ErrPacketDecode)
	}

	length := binary.LittleEndian.Uint32(append(buf, 0))

	// Sequence ID

	buf = make([]byte, 1)
	n, err = r.Read(buf)
	if err != nil {
		if errors.Is(err, os.ErrDeadlineExceeded) {
			return nil, err
		}
		return nil, fmt.Errorf("%w: %v", ErrPacketDecode, err)
	}
	if n != 1 {
		return nil, fmt.Errorf("%w: truncated header", ErrPacketDecode)
	}

	p.SequenceID = buf[0]

	// Payload

	buf = make([]byte, length)
	n, err = r.Read(buf)
	if err != nil {
		if errors.Is(err, os.ErrDeadlineExceeded) {
			return nil, err
		}
		return nil, fmt.Errorf("%w: %v", ErrPacketDecode, err)
	}
	if n != int(length) {
		return nil, fmt.Errorf("%w: truncated payload", ErrPacketDecode)
	}

	p.Payload = buf

	return p, nil
}
