package io

import (
	"bufio"
	"errors"
	"io"
	"os"
	"path/filepath"
)

type FileReader struct {
	files  []string
	pos    int
	fp     *os.File
	reader *bufio.Reader
}

func NewFileReader(f string) (*FileReader, error) {
	matches, _ := filepath.Glob(f)
	r := &FileReader{matches, 0, nil, nil}
	err := r.ready()
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (r *FileReader) ready() error {
	if r.fp != nil {
		r.fp.Close()
	}
	if len(r.files) == 0 {
		return errors.New("no files")
	} else if len(r.files) <= r.pos {
		return io.EOF
	} else {
		fp, err := os.Open(r.files[r.pos])
		if err != nil {
			return err
		}
		r.pos += 1
		r.fp = fp
		r.reader = bufio.NewReader(r.fp)
	}
	return nil
}

func (r *FileReader) Readln() (string, error) {
	l, _, err := r.reader.ReadLine()
	if err == io.EOF {
		err = r.ready()
		if err != nil {
			return "", err
		}
		return r.Readln()
	}
	return string(l), nil
}

func (r *FileReader) Close() {
	r.fp.Close()
}
