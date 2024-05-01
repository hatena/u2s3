package core

import (
	"errors"
	"io"
	"os"
	"path/filepath"

	gzip "github.com/klauspost/pgzip"

	"github.com/hatena/u2s3/pkg/config"
	"github.com/hatena/u2s3/pkg/util"
)

type FileAggregator struct {
	config *config.UploadConfig
	files  map[string]*File
}

type File struct {
	fn      string
	inFp    *os.File
	outFp   *os.File
	keyFmt  string
	keyTemp *config.UploadKeyTemplate
	seq     int
}

func NewFileAggregator(cfg *config.UploadConfig) (Aggregator, error) {
	if cfg.FileName == "" {
		return nil, errors.New("Filename must be specified when using filename-aware uploading.")
	}
	matches, _ := filepath.Glob(cfg.FileName)
	files := make(map[string]*File, len(matches))
	for _, f := range matches {
		files[f] = NewFile(f, cfg.FilenameFormat, cfg.KeyFormat, cfg.OutputPrefixKey)
	}
	return &FileAggregator{
		config: cfg,
		files:  files,
	}, nil
}

func (a *FileAggregator) Run() error {
	for _, f := range a.files {
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

func (a *FileAggregator) GenFetchJobs() chan *fetchJob {
	out := make(chan *fetchJob)
	go func() {
		for _, f := range a.files {
			key, err := f.GetObjectKey()
			if err != nil {
				continue
			}
			out <- &fetchJob{key, f}
		}
		close(out)
	}()
	return out
}

func (a *FileAggregator) Close() {
	for _, f := range a.files {
		f.Remove()
	}
}

func NewFile(fn string, nameFmt, keyFmt, output string) *File {
	params := util.GetParams(nameFmt, fn)
	keyTemp := &config.UploadKeyTemplate{
		Output: output,
		Year:   params["Year"],
		Month:  params["Month"],
		Day:    params["Day"],
		Hour:   params["Hour"],
		Minute: params["Minute"],
		Second: params["Second"],
	}
	return &File{
		fn:      fn,
		inFp:    nil,
		outFp:   nil,
		keyFmt:  keyFmt,
		keyTemp: keyTemp,
		seq:     0,
	}
}

func (f *File) GetObjectKey() (string, error) {
	f.seq += 1
	f.keyTemp.Seq = f.seq
	return util.GenerateUploadKey(f.keyTemp, f.keyFmt)
}

func (f *File) GetFile() *os.File {
	op, err := os.Open(f.outFp.Name())
	if err != nil {
		panic(err)
	}
	return op
}

func (f *File) ResetSeq() {
	f.seq = 0
}

func (f *File) Flush() {}

func (f *File) Remove() {
	f.inFp.Close()
	f.outFp.Close()
	os.Remove(f.outFp.Name())
}

func (f *File) compress() error {
	inFp, err := os.Open(f.fn)
	if err != nil {
		return err
	}
	f.inFp = inFp
	outFp, err := os.CreateTemp("", "u2s3")
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
	gw, err := gzip.NewWriterLevel(outFp, gzip.BestCompression)
	if err != nil {
		return err
	}
	_, err = io.Copy(gw, in)
	if err != nil {
		return err
	}
	gw.Close()
	outFp.Close()
	return nil
}
