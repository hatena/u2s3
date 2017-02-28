package io

import (
	"bufio"
	"os"
	gzip "github.com/klauspost/pgzip"
)

type StdinReader struct {
	reader *bufio.Reader
}

func NewStdinReader(gzipped bool) *StdinReader {
	var buf *bufio.Reader
	if gzipped {
		r, _ := gzip.NewReader(os.Stdin)
		buf = bufio.NewReader(r)
	} else {
		buf = bufio.NewReader(os.Stdin)
	}
	return &StdinReader{buf}
}

func (r *StdinReader) Readln() ([]byte, error) {
	l, _, err := r.reader.ReadLine()
	return l, err
}

func (r *StdinReader) Close() {}
