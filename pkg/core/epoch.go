package core

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"text/template"
	"time"

	gzip "github.com/klauspost/pgzip"
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
	fp, err := ioutil.TempFile("", "log2s3")
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
	host, err := os.Hostname()
	if err != nil {
		return "", err
	}
	temp := template.New("key")
	template.Must(temp.Parse(e.keyFmt))
	var res bytes.Buffer
	err = temp.Execute(&res, map[string]interface{}{
		"Output":   e.output,
		"Year":     fmt.Sprintf("%04d", t.Year()),
		"Month":    fmt.Sprintf("%02d", t.Month()),
		"Day":      fmt.Sprintf("%02d", t.Day()),
		"Hour":     fmt.Sprintf("%02d", t.Hour()),
		"Minute":   fmt.Sprintf("%02d", t.Minute()),
		"Second":   fmt.Sprintf("%02d", t.Second()),
		"Hostname": host,
		"Seq":      seq,
	})
	if err != nil {
		return "", err
	}
	return res.String(), nil
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

func (e *Epoch) Close() {
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
