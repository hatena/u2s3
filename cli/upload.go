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
		KeyFormat:       c.String("output-key"),
		OutputPrefixKey: c.String("output"),
		Step:            c.Int("step"),
		Bucket:          c.String("bucket"),
		MaxRetry:        c.Int("max-retry"),
		CPULimit:        c.Int("cpu"),
		MemoryLimit:     c.Int("memory"),
		RateLimit:       c.Int("rate"),
		Device:          c.String("dev"),
		ContentAware:    c.Bool("content-aware"),
		FilenameFormat:  c.String("filename-format"),
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

	var agg core.Aggregator
	if cfg.ContentAware {
		agg, err = core.NewEpochAggregator(cfg)
	} else {
		agg, err = core.NewFileAggregator(cfg)
	}
	if err != nil {
		return err
	}
	defer agg.Close()
	if err := agg.Run(); err != nil {
		return err
	}
	up := core.NewUploader(cfg)
	for _, f := range agg.GetUploadableFiles() {
		if err := up.Upload(f); err != nil {
			return err
		}
	}
	return nil
}
