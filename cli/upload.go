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

func uploadLogCmd(c *cli.Context) error {
	cfg := &pkg.UploadConfig{
		Step:            c.Int("step"),
		LogFormat:       c.String("log-format"),
		FileName:        c.String("file"),
		KeyFormat:       c.String("key-format"),
		OutputPrefixKey: c.String("output"),
		Bucket:          c.String("bucket"),
		MaxRetry:        c.Int("max-retry"),
		CPULimit:        c.Int("cpu"),
		MemoryLimit:     c.Int("memory"),
		RateLimit:       c.Int("rate"),
		Device:          c.String("dev"),
	}

	pp.Println(cfg)

	initResourceLimit(cfg)

	if cfg.Bucket == "" {
		return errors.New("Bucket name must be specified")
	}

	agg, err := core.NewEpochAggregator(cfg)
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

func uploadCmd(c *cli.Context) error {
	cfg := &pkg.UploadConfig{
		FilenameFormat:  c.String("filename-format"),
		FileName:        c.String("file"),
		KeyFormat:       c.String("key-format"),
		OutputPrefixKey: c.String("output"),
		Bucket:          c.String("bucket"),
		MaxRetry:        c.Int("max-retry"),
		CPULimit:        c.Int("cpu"),
		MemoryLimit:     c.Int("memory"),
		RateLimit:       c.Int("rate"),
		Device:          c.String("dev"),
	}

	pp.Println(cfg)

	initResourceLimit(cfg)

	if cfg.Bucket == "" {
		return errors.New("Bucket name must be specified")
	}

	agg, err := core.NewFileAggregator(cfg)
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

func initResourceLimit(cfg *pkg.UploadConfig) {
	cmngr, err := resourcelimit.NewCgroupMngr(cfg)
	if err == nil {
		defer cmngr.Close()
	} else {
		fmt.Println(err)
	}
}
