package core

import (
	"github.com/taku-k/u2s3/pkg"
	"github.com/taku-k/u2s3/pkg/aws"
)

type Uploader struct {
	cli *aws.S3Cli
}

func NewUploader(config *pkg.UploadConfig) *Uploader {
	cli := aws.NewS3Cli(config)
	return &Uploader{cli}
}

func (u *Uploader) Upload(e *Epoch) error {
	e.Close()

	seq := 0
	key := ""
	var err error
	for {
		seq += 1
		key, err = e.GetObjectKey(seq)
		if err != nil {
			return err
		}
		if !u.cli.HasKey(key) {
			break
		}
	}
	if err := u.cli.Upload(key, e.fp); err != nil {
		return err
	}
	return nil
}
