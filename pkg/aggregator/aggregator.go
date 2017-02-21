package aggregator

import (
	"io"

	lio "github.com/taku-k/log2s3-go/pkg/io"
)

type Aggregator struct {
	reader lio.BufferedReader
	cmpr   *Compressor
	up     *Uploader
}

func NewAggregator(reader lio.BufferedReader) *Aggregator {
	cmpr := NewCompressor()
	up := NewUploader()
	return &Aggregator{reader, cmpr, up}
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
		if !a.cmpr.HasEpochFile(epochKey) {
			epoch := NewEpoch()
			a.cmpr.AddEpoch(epoch)
			a.up.AddEpoch(epoch)
		}
		a.cmpr.Compress(epochKey, l)
	}
	a.up.Upload()
	return nil
}

func (a *Aggregator) parseEpoch(l string) string {
	return ""
}
