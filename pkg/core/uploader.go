package core

import (
	"os"

	"github.com/taku-k/u2s3/pkg"
	"github.com/taku-k/u2s3/pkg/aws"
)

type UploadableFile interface {
	GetObjectKey(seq int) (string, error)
	GetFile() *os.File
	Close()
}

type Uploader struct {
	cli *aws.S3Cli
}

func NewUploader(config *pkg.UploadConfig) *Uploader {
	cli := aws.NewS3Cli(config)
	return &Uploader{cli}
}

func (u *Uploader) Upload(uf UploadableFile) error {
	uf.Close()

	seq := 0
	key := ""
	var err error
	for {
		seq += 1
		key, err = uf.GetObjectKey(seq)
		if err != nil {
			return err
		}
		if !u.cli.HasKey(key) {
			break
		}
	}
	if err := u.cli.Upload(key, uf.GetFile()); err != nil {
		return err
	}
	return nil
}
