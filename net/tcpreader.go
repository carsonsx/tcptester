package net

import (
	"bufio"
	"encoding/binary"
	"io"
)

type Reader interface {
	Read(c *TCPClient) ([]byte, []byte, error)
}

type ProtobufReader struct {
}

func (r *ProtobufReader) Read(c *TCPClient) ([]byte, []byte, error) {
	b := make([]byte, 2)
	_, err := io.ReadFull(c.conn, b)
	if err != nil {
		return nil, nil, err
	}
	l := binary.BigEndian.Uint16(b)
	data := make([]byte, l)
	if l == 0 {
		return nil, nil, nil
	}
	_, err = io.ReadFull(c.conn, data)
	if err != nil {
		return nil, nil, err
	}
	return data, data, nil
}

type LineReader struct {
	_reader *bufio.Reader
}

func (r *LineReader) Read(c *TCPClient) ([]byte, []byte, error) {
	if r._reader == nil {
		r._reader = bufio.NewReader(c.conn)
	}
	line, err := r._reader.ReadSlice('\n')
	if err != nil {
		return line, nil, err
	}
	data := line
	if line[len(line)-1] == '\n' {
		drop := 1
		if len(line) > 1 && line[len(line)-2] == '\r' {
			drop = 2
		}
		data = line[:len(line)-drop]
	}
	return line, data, err
}
