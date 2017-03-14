package cli

import (
	"errors"
	"fmt"

	"github.com/k0kubun/pp"
	"github.com/taku-k/log2s3-go/pkg"
	"github.com/taku-k/log2s3-go/pkg/cgroup"
	"github.com/taku-k/log2s3-go/pkg/core"
	lio "github.com/taku-k/log2s3-go/pkg/io"
	"github.com/urfave/cli"
)

func uploadCmd(c *cli.Context) error {
	var reader lio.BufferedReader
	var err error
	cfg := &pkg.UploadConfig{
		FileName:        c.String("file"),
		LogFormat:       c.String("log-format"),
		KeyFormat:       c.String("key"),
		OutputPrefixKey: c.String("output"),
		Step:            c.Int("step"),
		Bucket:          c.String("bucket"),
		Gzipped:         c.Bool("gzipeed"),
		MaxRetry:        c.Int("max-retry"),
		CPULimit:        c.Int("cpu"),
		MemoryLimit:     c.Uint("memory"),
	}

	pp.Println(cfg)

	cmngr, err := cgroup.NewCgroupMngr(cfg)
	if err == nil {
		defer cmngr.Close()
	} else {
		fmt.Println(err)
	}

	if cfg.Bucket == "" {
		return errors.New("Bucket name must be specified")
	}
	if cfg.FileName != "" {
		reader, err = lio.NewFileReader(cfg.FileName, cfg.Gzipped)
		if err != nil {
			return err
		}
	} else {
		reader = lio.NewStdinReader(cfg.Gzipped)
	}

	agg := aggregator.NewAggregator(reader, cfg)
	return agg.Run()
}
