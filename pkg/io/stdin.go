package io

import (
	"bufio"
	"os"
)

type StdinReader struct {
	reader *bufio.Reader
}

func NewStdinReader() *StdinReader {
	return &StdinReader{
		bufio.NewReader(os.Stdin),
	}
}

func (r *StdinReader) Readln() ([]byte, error) {
	l, _, err := r.reader.ReadLine()
	return l, err
}

func (r *StdinReader) Close() {}
