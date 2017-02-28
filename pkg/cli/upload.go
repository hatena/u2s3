package cli

import (
	"errors"
	"github.com/taku-k/log2s3-go/pkg/aggregator"
	lio "github.com/taku-k/log2s3-go/pkg/io"
	"github.com/urfave/cli"
)

func uploadCmd(c *cli.Context) error {
	var reader lio.BufferedReader
	var err error
	file := c.String("file")
	logFmt := c.String("log-format")
	keyFmt := c.String("key")
	output := c.String("output")
	step := c.Int("step")
	bucket := c.String("bucket")
	gzipped := c.Bool("gzipped")

	if bucket == "" {
		return errors.New("Bucket name must be specified")
	}
	if file != "" {
		reader, err = lio.NewFileReader(file, gzipped)
		if err != nil {
			return err
		}
	} else {
		reader = lio.NewStdinReader(gzipped)
	}

	agg := aggregator.NewAggregator(reader, logFmt, keyFmt, output, bucket, step)
	return agg.Run()
}
