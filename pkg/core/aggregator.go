package core

import (
	"io"
	"regexp"
	"time"

	"github.com/taku-k/u2s3/pkg"
	"github.com/taku-k/u2s3/pkg/input/content"
)

type Aggregator interface {
	Run() error
	GetUploadableFiles() []UploadableFile
	Close()
}

var reTsv = regexp.MustCompile(`(?:^|[ \t])time\:([^\t]+)`)

type EpochAggregator struct {
	reader content.BufferedReader
	mngr   *EpochManager
	config *pkg.UploadConfig
}

func NewEpochAggregator(cfg *pkg.UploadConfig) (Aggregator, error) {
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
	return e.Format("20060102150405")
}
