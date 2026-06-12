package stream

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
)

// Buf encodes/decodes VC:MP client script streams (little-endian ints,
// big-endian 16-bit length prefixed strings).
type Buf struct {
	data []byte
	pos  int
}

func NewWriter() *Buf {
	return &Buf{}
}

func NewReader(data []byte) *Buf {
	return &Buf{data: data}
}

func (s *Buf) WriteInt(v int32) {
	var b [4]byte
	binary.LittleEndian.PutUint32(b[:], uint32(v))
	s.data = append(s.data, b[:]...)
}

func (s *Buf) WriteFloat(v float32) {
	var b [4]byte
	binary.LittleEndian.PutUint32(b[:], uint32(math.Float32bits(v)))
	s.data = append(s.data, b[:]...)
}

func (s *Buf) WriteString(str string) {
	b := []byte(str)
	if len(b) > 0xFFFF {
		b = b[:0xFFFF]
	}
	var lenb [2]byte
	binary.BigEndian.PutUint16(lenb[:], uint16(len(b)))
	s.data = append(s.data, lenb[:]...)
	s.data = append(s.data, b...)
}

func (s *Buf) Bytes() []byte {
	return s.data
}

func (s *Buf) ReadInt() (int32, error) {
	if s.pos+4 > len(s.data) {
		return 0, io.EOF
	}
	v := binary.LittleEndian.Uint32(s.data[s.pos:])
	s.pos += 4
	return int32(v), nil
}

func (s *Buf) ReadFloat() (float32, error) {
	if s.pos+4 > len(s.data) {
		return 0, io.EOF
	}
	bits := binary.LittleEndian.Uint32(s.data[s.pos:])
	s.pos += 4
	return math.Float32frombits(bits), nil
}

func (s *Buf) ReadString() (string, error) {
	if s.pos+2 > len(s.data) {
		return "", io.EOF
	}
	n := int(binary.BigEndian.Uint16(s.data[s.pos:]))
	s.pos += 2
	if s.pos+n > len(s.data) {
		return "", fmt.Errorf("stream string length %d exceeds remaining %d bytes", n, len(s.data)-s.pos)
	}
	str := string(s.data[s.pos : s.pos+n])
	s.pos += n
	return str, nil
}
