package aggregator

import (
	"io"
	"regexp"

	lio "github.com/taku-k/log2s3-go/pkg/io"
)

var reTsv = regexp.MustCompile(`(?:^|[ \t])time\:([^\t]+)`)

type Aggregator struct {
	reader  lio.BufferedReader
	mngr    *EpochManager
	cmpr    *Compressor
	up      *Uploader
	logType string
	keyFmt  string
}

func NewAggregator(reader lio.BufferedReader, logType, keyFmt string) *Aggregator {
	mngr := NewEpochManager()
	cmpr := NewCompressor()
	up := NewUploader()
	return &Aggregator{
		reader:  reader,
		mngr:    mngr,
		cmpr:    cmpr,
		up:      up,
		logType: logType,
		keyFmt:  keyFmt}
}

func (a *Aggregator) Run() error {
	for {
		l, err := a.reader.Readln()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		epochKey := a.parseEpoch(l)
		if epochKey == "" {
			continue
		}
		var epoch *Epoch
		if !a.mngr.HasEpoch(epochKey) {
			epoch = NewEpoch(epochKey, a.keyFmt)
		} else {
			epoch = a.mngr.GetEpoch(epochKey)
		}
		a.cmpr.Compress(epoch, l)
	}
	a.up.Upload()
	return nil
}

func (a *Aggregator) parseEpoch(l string) string {
	s := ""
	switch a.logType {
	case "ssv":
		break
	case "tsv":
		s = reTsv.FindString(l)
		break
	}
	return s
}
