package cgroups

import "github.com/taku-k/log2s3-go/pkg"

type DummyCgroupMngr struct{}

func (c *DummyCgroupMngr) Close() {}

func NewCgroupMngr(c *pkg.UploadConfig) (*DummyCgroupMngr, error) {
	return &DummyCgroupMngr{}, nil
}
