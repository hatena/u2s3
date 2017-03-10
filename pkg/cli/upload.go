package cli

import (
	"errors"
	"fmt"
	"os"

	"github.com/crosbymichael/cgroups"
	specs "github.com/opencontainers/runtime-spec/specs-go"
	"github.com/taku-k/log2s3-go/pkg"
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
	}

	shares := uint64(10)
	ctrl, err := cgroups.New(cgroups.V1, cgroups.StaticPath("/test"), &specs.LinuxResources{
		CPU: &specs.LinuxCPU{
			Shares: &shares,
			Cpus:   "0",
			Mems:   "0",
		},
	})
	defer ctrl.Delete()
	pid := os.Getpid()
	if err := ctrl.Add(cgroups.Process{Pid: pid}); err != nil {
		panic(err)
	}

	stats, err := ctrl.Stat()
	fmt.Println(stats)

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
