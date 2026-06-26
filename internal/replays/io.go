package replays

import (
	"bytes"
	"encoding/binary"
)

func writeU8(buf *bytes.Buffer, value uint8) {
	buf.WriteByte(value)
}

func writeU16(buf *bytes.Buffer, value uint16) {
	var b [2]byte
	binary.LittleEndian.PutUint16(b[:], value)
	buf.Write(b[:])
}

func writeU32(buf *bytes.Buffer, value uint32) {
	var b [4]byte
	binary.LittleEndian.PutUint32(b[:], value)
	buf.Write(b[:])
}

func writeU64(buf *bytes.Buffer, value uint64) {
	var b [8]byte
	binary.LittleEndian.PutUint64(b[:], value)
	buf.Write(b[:])
}

func writeS32(buf *bytes.Buffer, value int32) {
	writeU32(buf, uint32(value))
}

func writeS64(buf *bytes.Buffer, value int64) {
	writeU64(buf, uint64(value))
}

func writeBool(buf *bytes.Buffer, value bool) {
	if value {
		writeU8(buf, 1)
		return
	}
	writeU8(buf, 0)
}

func writeULEB128(buf *bytes.Buffer, value uint64) {
	if value == 0 {
		buf.WriteByte(0)
		return
	}
	for value != 0 {
		current := byte(value & 0x7f)
		value >>= 7
		if value != 0 {
			current |= 0x80
		}
		buf.WriteByte(current)
	}
}

func writeString(buf *bytes.Buffer, value string) {
	if value == "" {
		writeU8(buf, 0x00)
		return
	}
	writeU8(buf, 0x0b)
	writeULEB128(buf, uint64(len(value)))
	buf.WriteString(value)
}
