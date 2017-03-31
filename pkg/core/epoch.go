package core

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	gzip "github.com/klauspost/pgzip"
	"github.com/taku-k/u2s3/pkg"
	"github.com/taku-k/u2s3/pkg/util"
)

type Epoch struct {
	fp       *os.File
	buf      *bufio.Writer
	writer   *gzip.Writer
	epochKey string
	keyFmt   string
	output   string
}

type EpochManager struct {
	epochs map[string]*Epoch
}

func NewEpoch(epochKey, keyFmt, output string) (*Epoch, error) {
	fp, err := ioutil.TempFile("", "u2s3")
	if err != nil {
		return nil, err
	}
	buf := bufio.NewWriter(fp)
	w, _ := gzip.NewWriterLevel(buf, gzip.BestCompression)
	return &Epoch{
		fp:       fp,
		buf:      buf,
		writer:   w,
		epochKey: epochKey,
		keyFmt:   keyFmt,
		output:   output,
	}, nil
}

func (e *Epoch) GetObjectKey(seq int) (string, error) {
	t, err := time.Parse("20060102150405", e.epochKey)
	if err != nil {
		return "", err
	}
	keyTemp := &pkg.UploadKeyTemplate{
		Output: e.output,
		Year:   fmt.Sprintf("%04d", t.Year()),
		Month:  fmt.Sprintf("%02d", t.Month()),
		Day:    fmt.Sprintf("%02d", t.Day()),
		Hour:   fmt.Sprintf("%02d", t.Hour()),
		Minute: fmt.Sprintf("%02d", t.Minute()),
		Second: fmt.Sprintf("%02d", t.Second()),
		Seq:    seq,
	}
	return util.GenerateUploadKey(keyTemp, e.keyFmt)
}

func (e *Epoch) GetFile() *os.File {
	return e.fp
}

func (e *Epoch) Write(l []byte) {
	_, err := e.writer.Write(l)
	if err != nil {
		return
	}
}

func (e *Epoch) Flush() {
	e.writer.Close()
	e.buf.Flush()
}

func (e *Epoch) Remove() {
	e.writer.Close()
	e.buf.Flush()
	e.fp.Close()
	os.Remove(e.fp.Name())
}

func NewEpochManager() *EpochManager {
	return &EpochManager{
		epochs: make(map[string]*Epoch, 100),
	}
}

func (m *EpochManager) HasEpoch(key string) bool {
	_, ok := m.epochs[key]
	return ok
}

func (m *EpochManager) GetEpoch(key string) *Epoch {
	return m.epochs[key]
}

func (m *EpochManager) PutEpoch(e *Epoch) {
	m.epochs[e.epochKey] = e
}

func (m *EpochManager) Close() {
	for _, e := range m.epochs {
		e.Remove()
	}
}
