package cli

import (
	"errors"
	"fmt"

	"github.com/k0kubun/pp"
	"github.com/taku-k/u2s3/pkg"
	"github.com/taku-k/u2s3/pkg/core"
	"github.com/taku-k/u2s3/pkg/resourcelimit"
	"github.com/urfave/cli"
)

func uploadCmd(c *cli.Context) error {
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
		ContentAware:    c.Bool("content-aware"),
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

	if cfg.ContentAware {
		agg, err := core.NewAggregator(cfg)
		if err != nil {
			return err
		}
		return agg.Run()
	}
	return nil
}
