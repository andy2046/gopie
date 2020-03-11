// Package tlv implements Type-Length-Value encoding.
package tlv

import (
	"bytes"
	"encoding/binary"
	"io"
)

// ByteSize const.
const (
	Bytes1 ByteSize = 1 << iota
	Bytes2
	Bytes4
	Bytes8
)

type (
	// ByteSize is the size of a field in bytes.
	// Used to define the size of the type and length field in a message.
	ByteSize int

	// Record represents a record of data encoded in the TLV message.
	Record struct {
		Payload []byte
		Type    uint
	}

	// Codec is the configuration for a TLV encoding/decoding task.
	Codec struct {
		// TypeBytes defines the size in bytes of the message type field.
		TypeBytes ByteSize

		// LenBytes defines the size in bytes of the message length field.
		LenBytes ByteSize
	}

	// Writer encodes records into TLV format using a Codec and writes into the provided io.Writer.
	Writer struct {
		writer io.Writer
		codec  *Codec
	}

	// Reader decodes records from TLV format using a Codec from the provided io.Reader.
	Reader struct {
		codec  *Codec
		reader io.Reader
	}
)

// NewWriter creates a new Writer.
func NewWriter(w io.Writer, codec *Codec) *Writer {
	return &Writer{
		codec:  codec,
		writer: w,
	}
}

// Write encodes records into TLV format using a Codec and writes into the provided io.Writer.
func (w *Writer) Write(rec *Record) error {
	err := writeUint(w.writer, w.codec.TypeBytes, rec.Type)
	if err != nil {
		return err
	}

	ulen := uint(len(rec.Payload))
	err = writeUint(w.writer, w.codec.LenBytes, ulen)
	if err != nil {
		return err
	}

	_, err = w.writer.Write(rec.Payload)
	return err
}

func writeUint(w io.Writer, b ByteSize, i uint) error {
	var num interface{}
	switch b {
	case Bytes1:
		num = uint8(i)
	case Bytes2:
		num = uint16(i)
	case Bytes4:
		num = uint32(i)
	case Bytes8:
		num = uint64(i)
	}
	return binary.Write(w, binary.LittleEndian, num)
}

// NewReader creates a new Reader.
func NewReader(reader io.Reader, codec *Codec) *Reader {
	return &Reader{codec: codec, reader: reader}
}

// Next reads a single Record from the io.Reader.
func (r *Reader) Next() (*Record, error) {
	// get type
	typeBytes := make([]byte, r.codec.TypeBytes)
	_, err := r.reader.Read(typeBytes)
	if err != nil {
		return nil, err
	}
	typ := readUint(typeBytes, r.codec.TypeBytes)

	// get len
	payloadLenBytes := make([]byte, r.codec.LenBytes)
	_, err = r.reader.Read(payloadLenBytes)
	if err != nil && err != io.EOF {
		return nil, err
	}
	payloadLen := readUint(payloadLenBytes, r.codec.LenBytes)

	if err == io.EOF && payloadLen != 0 {
		return nil, err
	}

	// get value
	v := make([]byte, payloadLen)
	_, err = r.reader.Read(v)
	if err != nil && err != io.EOF {
		return nil, err
	}

	return &Record{
		Type:    typ,
		Payload: v,
	}, nil

}

func readUint(b []byte, sz ByteSize) uint {
	reader := bytes.NewReader(b)
	switch sz {
	case Bytes1:
		var i uint8
		binary.Read(reader, binary.LittleEndian, &i)
		return uint(i)
	case Bytes2:
		var i uint16
		binary.Read(reader, binary.LittleEndian, &i)
		return uint(i)
	case Bytes4:
		var i uint32
		binary.Read(reader, binary.LittleEndian, &i)
		return uint(i)
	case Bytes8:
		var i uint64
		binary.Read(reader, binary.LittleEndian, &i)
		return uint(i)
	default:
		return 0
	}
}
