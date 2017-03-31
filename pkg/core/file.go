package core

import (
	"bufio"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	gzip "github.com/klauspost/pgzip"
	"github.com/taku-k/u2s3/pkg"
	"github.com/taku-k/u2s3/pkg/util"
)

type File struct {
	fn    string
	inFp  *os.File
	outFp *os.File
}

type FileAggregator struct {
	config *pkg.UploadConfig
	files  map[string]*File
}

func NewFileAggregator(cfg *pkg.UploadConfig) (Aggregator, error) {
	if cfg.FileName == "" {
		return nil, errors.New("Filename must be specified when using filename-aware uploading.")
	}
	matches, _ := filepath.Glob(cfg.FileName)
	files := make(map[string]*File, len(matches))
	for _, f := range matches {
		files[f] = nil
	}
	return &FileAggregator{
		config: cfg,
		files:  files,
	}, nil
}

func (a *FileAggregator) Run() error {
	for fn := range a.files {
		f, err := NewFile(fn)
		if err != nil {
			return err
		}
		a.files[fn] = f
		f.compress()
	}
	return nil
}

func (a *FileAggregator) GetUploadableFiles() []UploadableFile {
	v := make([]UploadableFile, 0, len(a.files))
	for _, f := range a.files {
		v = append(v, f)
	}
	return v
}

func (a *FileAggregator) Close() {
	for _, f := range a.files {
		f.Remove()
	}
}

func NewFile(fn string) (*File, error) {
	return &File{
		fn:    fn,
		inFp:  nil,
		outFp: nil,
	}, nil
}

func (f *File) GetObjectKey(seq int) (string, error) {
	return "", nil
}

func (f *File) GetFile() *os.File {
	return f.outFp
}

func (f *File) Close() {
	f.inFp.Close()
	f.outFp.Close()
}

func (f *File) Remove() {
	os.Remove(f.outFp.Name())
}

func (f *File) compress() error {
	inFp, err := os.Open(f.fn)
	if err != nil {
		return err
	}
	f.inFp = inFp
	outFp, err := ioutil.TempFile("", "u2s3")
	if err != nil {
		return err
	}
	f.outFp = outFp
	var in io.Reader
	if util.IsGzipped(inFp) {
		in, err = gzip.NewReader(inFp)
		if err != nil {
			return err
		}
	} else {
		in = inFp
	}
	scanner := bufio.NewScanner(in)
	outBuf := bufio.NewWriter(outFp)
	w, _ := gzip.NewWriterLevel(outBuf, gzip.BestCompression)
	for scanner.Scan() {
		w.Write(scanner.Bytes())
	}
	inFp.Close()
	outFp.Close()
	return nil
}
