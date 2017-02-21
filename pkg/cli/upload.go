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
	agg := aggregator.NewAggregator(reader)
	return agg.Run()
}
