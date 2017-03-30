package core

import (
	"io"
	"regexp"
	"time"

	"github.com/taku-k/log2s3-go/pkg"
	"github.com/taku-k/log2s3-go/pkg/input/content"
)

var reTsv = regexp.MustCompile(`(?:^|[ \t])time\:([^\t]+)`)

type Aggregator struct {
	reader content.BufferedReader
	mngr   *EpochManager
	up     *Uploader
	config *pkg.UploadConfig
}

func NewAggregator(reader content.BufferedReader, cfg *pkg.UploadConfig) *Aggregator {
	mngr := NewEpochManager()
	up := NewUploader(cfg)
	return &Aggregator{
		reader: reader,
		mngr:   mngr,
		up:     up,
		config: cfg,
	}
}

func (a *Aggregator) Run() error {
	defer a.Close()

	for {
		l, err := a.reader.Readln()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		epochKey := a.parseEpoch(string(l))
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
	for _, e := range a.mngr.epochs {
		if err := a.up.Upload(e); err != nil {
			return err
		}
	}
	return nil
}

func (a *Aggregator) Close() {
	a.reader.Close()
	a.mngr.Close()
}

func (a *Aggregator) parseEpoch(l string) string {
	r := ""
	switch a.config.LogFormat {
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
	e := time.Unix(t.Unix()-t.Unix()%(int64(a.config.Step)*60), 0)
	return e.Format("20060102150405")
}
