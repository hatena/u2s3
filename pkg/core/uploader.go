package core

import (
	"log"
	"os"

	"github.com/taku-k/u2s3/pkg/aws"
	"github.com/taku-k/u2s3/pkg/config"
)

type UploadableFile interface {
	GetObjectKey(seq int) (string, error)
	GetFile() *os.File
	Flush()
}

type Uploader struct {
	cli *aws.S3Cli
}

func NewUploader(config *config.UploadConfig) *Uploader {
	cli := aws.NewS3Cli(config)
	return &Uploader{cli}
}

func (u *Uploader) Upload(uf UploadableFile) error {
	uf.Flush()

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
	log.Printf("[info] Uploaded %s\n", key)
	return nil
}
