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

func (r *StdinReader) Readln() (string, error) {
	l, _, err := r.reader.ReadLine()
	return string(l), err
}

func (r *StdinReader) Close() {}
