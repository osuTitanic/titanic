package replays

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
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

func readU8(reader io.Reader) (uint8, error) {
	var b [1]byte
	if _, err := io.ReadFull(reader, b[:]); err != nil {
		return 0, err
	}
	return b[0], nil
}

func readU16(reader io.Reader) (uint16, error) {
	var b [2]byte
	if _, err := io.ReadFull(reader, b[:]); err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint16(b[:]), nil
}

func readU32(reader io.Reader) (uint32, error) {
	var b [4]byte
	if _, err := io.ReadFull(reader, b[:]); err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint32(b[:]), nil
}

func readU64(reader io.Reader) (uint64, error) {
	var b [8]byte
	if _, err := io.ReadFull(reader, b[:]); err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint64(b[:]), nil
}

func readS32(reader io.Reader) (int32, error) {
	value, err := readU32(reader)
	return int32(value), err
}

func readS64(reader io.Reader) (int64, error) {
	value, err := readU64(reader)
	return int64(value), err
}

func readBool(reader io.Reader) (bool, error) {
	value, err := readU8(reader)
	if err != nil {
		return false, err
	}

	switch value {
	case 0:
		return false, nil
	case 1:
		return true, nil
	default:
		return false, fmt.Errorf("invalid boolean value: %d", value)
	}
}

func readULEB128(reader io.Reader) (uint64, error) {
	var value uint64

	for shift := uint(0); shift < 64; shift += 7 {
		b, err := readU8(reader)
		if err != nil {
			return 0, err
		}
		if shift == 63 && b > 1 {
			return 0, fmt.Errorf("uleb128 overflow")
		}
		value |= uint64(b&0x7f) << shift

		if b&0x80 == 0 {
			return value, nil
		}
	}

	return 0, fmt.Errorf("uleb128 overflow")
}

func readString(reader io.Reader) (string, error) {
	marker, err := readU8(reader)
	if err != nil {
		return "", err
	}

	switch marker {
	case 0x00:
		return "", nil

	case 0x0b:
		length, err := readULEB128(reader)
		if err != nil {
			return "", fmt.Errorf("read string length: %w", err)
		}

		data := make([]byte, int(length))
		if _, err := io.ReadFull(reader, data); err != nil {
			return "", err
		}
		return string(data), nil

	default:
		return "", fmt.Errorf("invalid string marker: %#x", marker)
	}
}
