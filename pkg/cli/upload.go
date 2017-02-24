package cli

import (
	"github.com/taku-k/log2s3-go/pkg/aggregator"
	lio "github.com/taku-k/log2s3-go/pkg/io"
	"github.com/urfave/cli"
	"errors"
)

func uploadCmd(c *cli.Context) error {
	var reader lio.BufferedReader
	var err error
	logFmt := c.String("log-format")
	keyFmt := c.String("key")
	output := c.String("output")
	step := c.Int("step")
	bucket := c.String("bucket")

	if bucket == "" {
		return errors.New("Bucket name must be specified")
	}
	if f := c.String("file"); f != "" {
		reader, err = lio.NewFileReader(c.String("file"))
		if err != nil {
			return err
		}
	} else {
		reader = lio.NewStdinReader()
	}

	agg := aggregator.NewAggregator(reader, logFmt, keyFmt, output, bucket, step)
	return agg.Run()
}
