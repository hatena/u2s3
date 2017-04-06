package core

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"time"

	gzip "github.com/klauspost/pgzip"
	"github.com/taku-k/u2s3/pkg/config"
	"github.com/taku-k/u2s3/pkg/input/content"
	"github.com/taku-k/u2s3/pkg/util"
)

var reTsv = regexp.MustCompile(`(?:^|[ \t])time\:([^\t]+)`)

type EpochAggregator struct {
	reader content.BufferedReader
	mngr   *EpochManager
	config *config.UploadConfig
}

type Epoch struct {
	fp       *os.File
	buf      *bufio.Writer
	writer   *gzip.Writer
	epochKey string
	keyFmt   string
	output   string
	seq      int
}

type EpochManager struct {
	epochs map[string]*Epoch
}

func NewEpochAggregator(cfg *config.UploadConfig) (Aggregator, error) {
	mngr := NewEpochManager()
	var reader content.BufferedReader
	var err error
	if cfg.FileName != "" {
		reader, err = content.NewFileReader(cfg.FileName)
		if err != nil {
			return nil, err
		}
	} else {
		reader = content.NewStdinReader()
	}
	return &EpochAggregator{
		reader: reader,
		mngr:   mngr,
		config: cfg,
	}, nil
}

func (a *EpochAggregator) Run() error {
	for {
		l, err := a.reader.Readln()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		epochKey := parseEpoch(string(l), a.config.LogFormat, a.config.Step)
		if epochKey == "" {
			continue
		}
		var epoch *Epoch
		if !a.mngr.HasEpoch(epochKey) {
			epoch, err = NewEpoch(epochKey, a.config.KeyFormat, a.config.OutputPrefixKey)
			if err != nil {
				return err
			}
			a.mngr.PutEpoch(epoch)
		} else {
			epoch = a.mngr.GetEpoch(epochKey)
		}
		epoch.Write(l)
	}
	return nil
}

func (a *EpochAggregator) GetUploadableFiles() []UploadableFile {
	v := make([]UploadableFile, 0, len(a.mngr.epochs))
	for _, e := range a.mngr.epochs {
		v = append(v, e)
	}
	return v
}

func (a *EpochAggregator) GenFetchJobs() chan *fetchJob {
	out := make(chan *fetchJob)
	go func() {
		for _, e := range a.mngr.epochs {
			key, err := e.GetObjectKey()
			if err != nil {
				continue
			}
			out <- &fetchJob{key, e}
		}
		close(out)
	}()
	return out
}

func (a *EpochAggregator) Close() {
	a.reader.Close()
	a.mngr.Close()
}

func parseEpoch(l, logFormat string, step int) string {
	r := ""
	switch logFormat {
	case "ssv":
		break
	case "tsv":
		m := reTsv.FindStringSubmatch(l)
		if len(m) == 2 {
			r = m[1]
		}
		if len(r) >= 2 {
			if r[0] == '[' && r[len(r)-1] == ']' {
				r = r[1 : len(r)-1]
			}
		}
		break
	}
	if r == "" {
		return ""
	}
	t, err := time.Parse("02/Jan/2006:15:04:05 -0700", r)
	if err != nil {
		return ""
	}
	e := time.Unix(t.Unix()-t.Unix()%(int64(step)*60), 0)
	return e.In(t.Location()).Format("20060102150405")
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
		seq:      0,
	}, nil
}

func (e *Epoch) GetObjectKey() (string, error) {
	t, err := time.Parse("20060102150405", e.epochKey)
	if err != nil {
		return "", err
	}
	e.seq += 1
	keyTemp := &config.UploadKeyTemplate{
		Output: e.output,
		Year:   fmt.Sprintf("%04d", t.Year()),
		Month:  fmt.Sprintf("%02d", t.Month()),
		Day:    fmt.Sprintf("%02d", t.Day()),
		Hour:   fmt.Sprintf("%02d", t.Hour()),
		Minute: fmt.Sprintf("%02d", t.Minute()),
		Second: fmt.Sprintf("%02d", t.Second()),
		Seq:    e.seq,
	}
	return util.GenerateUploadKey(keyTemp, e.keyFmt)
}

func (e *Epoch) GetFile() *os.File {
	return e.fp
}

func (e *Epoch) ResetSeq() {
	e.seq = 0
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
