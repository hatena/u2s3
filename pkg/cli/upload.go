package cli

import (
	"github.com/taku-k/log2s3-go/pkg/aggregator"
	lio "github.com/taku-k/log2s3-go/pkg/io"
	"github.com/urfave/cli"
)

func uploadCmd(c *cli.Context) error {
	var reader lio.BufferedReader
	var err error
	if f := c.String("file"); f != "" {
		reader, err = lio.NewFileReader(c.String("file"))
		if err != nil {
			return err
		}
	} else {
		reader = lio.NewStdinReader()
	}
	logFmt := c.String("log-format")
	keyFmt := c.String("key")
	output := c.String("output")
	step := c.Int("step")
	agg := aggregator.NewAggregator(reader, logFmt, keyFmt, output, step)
	return agg.Run()
}
