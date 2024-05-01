package cli

import (
	"errors"
	"log"

	"github.com/k0kubun/pp"
	"github.com/hatena/u2s3/pkg/config"
	"github.com/hatena/u2s3/pkg/core"
	"github.com/hatena/u2s3/pkg/resourcelimit"
	"github.com/urfave/cli"
)

func uploadLogCmd(c *cli.Context) error {
	cfg := &config.UploadConfig{
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

	return uploadCmdBase(cfg, core.NewEpochAggregator)
}

func uploadFileCmd(c *cli.Context) error {
	cfg := &config.UploadConfig{
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

	return uploadCmdBase(cfg, core.NewFileAggregator)
}

func uploadCmdBase(cfg *config.UploadConfig, newAggFunc func(cfg *config.UploadConfig) (core.Aggregator, error)) error {
	pp.Println(cfg)

	initResourceLimit(cfg)

	if cfg.Bucket == "" {
		return errors.New("Bucket name must be specified")
	}

	agg, err := newAggFunc(cfg)
	if err != nil {
		return err
	}
	defer agg.Close()
	if err := agg.Run(); err != nil {
		return err
	}
	return upload(cfg, agg.GetUploadableFiles())
}

func initResourceLimit(cfg *config.UploadConfig) {
	cmngr, err := resourcelimit.NewCgroupMngr(cfg)
	if err == nil {
		defer cmngr.Close()
	} else {
		log.Printf("[info] %s\n", err)
	}
}

func upload(cfg *config.UploadConfig, files []core.UploadableFile) error {
	up := core.NewUploader(cfg)
	for _, f := range files {
		if err := up.Upload(f); err != nil {
			return err
		}
	}
	return nil
}
