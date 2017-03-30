package content

import (
	"bufio"
	"errors"
	"io"
	"os"
	"path/filepath"
	gzip "github.com/klauspost/pgzip"
)

type FileReader struct {
	files  []string
	pos    int
	fp     *os.File
	reader *bufio.Reader
	gzipped bool
}

func NewFileReader(f string, gzipped bool) (*FileReader, error) {
	matches, _ := filepath.Glob(f)
	r := &FileReader{matches, 0, nil, nil, gzipped}
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
		if r.gzipped {
			g, err := gzip.NewReader(r.fp)
			if err != nil {
				return err
			}
			r.reader = bufio.NewReader(g)
		} else {
			r.reader = bufio.NewReader(r.fp)
		}
	}
	return nil
}

func (r *FileReader) Readln() ([]byte, error) {
	l, err := r.reader.ReadBytes('\n')
	if err == io.EOF {
		err = r.ready()
		if err != nil {
			return nil, err
		}
		return r.Readln()
	}
	return l, nil
}

func (r *FileReader) Close() {
	r.fp.Close()
}
