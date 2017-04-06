package cli

import (
	"github.com/taku-k/u2s3/pkg/config"
	"github.com/taku-k/u2s3/pkg/core"
	"github.com/urfave/cli"
)

func syncFileCmd(c *cli.Context) error {
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

	agg, err := core.NewFileAggregator(cfg)
	if err != nil {
		return err
	}
	jobs := agg.GenFetchJobs()
	up := core.NewUploader(cfg)
	for f := range core.SelectUploadFiles(5, jobs) {
		if err := up.Upload(f); err != nil {
			return err
		}
	}
	return nil
}
