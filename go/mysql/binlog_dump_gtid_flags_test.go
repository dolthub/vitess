/*
Copyright 2026 The Vitess Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package mysql

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"testing"
)

func gtidFlagsTestPacket(t *testing.T, flags uint16, payload []byte) []byte {
	t.Helper()
	var packet bytes.Buffer
	filename := "binlog-main.000002"
	fields := []any{byte(ComBinlogDumpGTID), flags, uint32(42), uint32(len(filename)),
		[]byte(filename), uint64(4), uint32(len(payload)), payload}
	for _, field := range fields {
		if err := binary.Write(&packet, binary.LittleEndian, field); err != nil {
			t.Fatal(err)
		}
	}
	return packet.Bytes()
}

func TestComBinlogDumpGTIDFlags(t *testing.T) {
	const expected = "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee:1-10"
	set, err := ParseMysql56GTIDSet(expected)
	if err != nil {
		t.Fatal(err)
	}
	payload := set.(Mysql56GTIDSet).SIDBlock()
	for _, flags := range []uint16{0, BinlogThroughGTID, BinlogThroughPosition} {
		t.Run(fmt.Sprintf("flags_%d", flags), func(t *testing.T) {
			filename, pos, result, err := new(Conn).parseComBinlogDumpGTID(gtidFlagsTestPacket(t, flags, payload))
			if err != nil {
				t.Fatal(err)
			}
			if filename != "binlog-main.000002" || pos != 4 {
				t.Fatalf("wrong file/position: %q/%d", filename, pos)
			}
			if result.GTIDSet == nil || result.GTIDSet.String() != expected {
				t.Fatalf("GTID payload lost with flags=%d: %v", flags, result.GTIDSet)
			}
		})
	}
	t.Run("no_payload", func(t *testing.T) {
		_, _, result, err := new(Conn).parseComBinlogDumpGTID(gtidFlagsTestPacket(t, 0, nil))
		if err != nil || result.GTIDSet != nil {
			t.Fatalf("no payload: result=%v err=%v", result.GTIDSet, err)
		}
	})
	t.Run("empty_set", func(t *testing.T) {
		_, _, result, err := new(Conn).parseComBinlogDumpGTID(gtidFlagsTestPacket(t, 0, make([]byte, 8)))
		if err != nil || result.GTIDSet == nil || result.GTIDSet.String() != "" {
			t.Fatalf("empty SID block: result=%v err=%v", result.GTIDSet, err)
		}
	})
	t.Run("nonblocking", func(t *testing.T) {
		_, _, result, err := new(Conn).parseComBinlogDumpGTID(gtidFlagsTestPacket(t, BinlogDumpNonBlock, payload))
		if !errors.Is(err, io.EOF) {
			t.Fatalf("expected EOF, got %v", err)
		}
		if result.GTIDSet == nil || result.GTIDSet.String() != expected {
			t.Fatalf("nonblocking GTID payload lost: %v", result.GTIDSet)
		}
	})
	t.Run("malformed_payload", func(t *testing.T) {
		_, _, _, err := new(Conn).parseComBinlogDumpGTID(gtidFlagsTestPacket(t, 0, []byte{1}))
		if err == nil {
			t.Fatal("accepted malformed SID block")
		}
	})
	t.Run("truncated_payload", func(t *testing.T) {
		packet := gtidFlagsTestPacket(t, 0, payload)
		_, _, _, err := new(Conn).parseComBinlogDumpGTID(packet[:len(packet)-1])
		if err == nil {
			t.Fatal("accepted truncated SID block")
		}
	})
}
