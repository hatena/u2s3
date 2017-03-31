package cli

import (
	"errors"
	"fmt"

	"github.com/k0kubun/pp"
	"github.com/taku-k/u2s3/pkg"
	"github.com/taku-k/u2s3/pkg/core"
	"github.com/taku-k/u2s3/pkg/input/content"
	"github.com/taku-k/u2s3/pkg/resourcelimit"
	"github.com/urfave/cli"
)

func uploadCmd(c *cli.Context) error {
	var reader content.BufferedReader
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
		MemoryLimit:     c.Int("memory"),
		RateLimit:       c.Int("rate"),
		Device:          c.String("dev"),
	}

	pp.Println(cfg)

	cmngr, err := resourcelimit.NewCgroupMngr(cfg)
	if err == nil {
		defer cmngr.Close()
	} else {
		fmt.Println(err)
	}

	if cfg.Bucket == "" {
		return errors.New("Bucket name must be specified")
	}
	if cfg.FileName != "" {
		reader, err = content.NewFileReader(cfg.FileName, cfg.Gzipped)
		if err != nil {
			return err
		}
	} else {
		reader = content.NewStdinReader(cfg.Gzipped)
	}

	agg := core.NewAggregator(reader, cfg)
	return agg.Run()
}
