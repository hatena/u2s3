package io

import (
	"bufio"
	"os"
)

type BufferedReader interface {
	Readln() (string, error)
}

type StdinReader struct {
	reader *bufio.Reader
}

func NewStdinReader() *StdinReader {
	return &StdinReader{
		bufio.NewReader(os.Stdin),
	}
}

func (r *StdinReader) Readln() (string, error) {
	return r.reader.ReadString('\n')
}