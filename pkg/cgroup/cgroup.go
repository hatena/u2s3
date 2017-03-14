package cgroup

import (
	"errors"
	"os"

	specs "github.com/opencontainers/runtime-spec/specs-go"
	"github.com/taku-k/cgroups"
	"github.com/taku-k/log2s3-go/pkg"
)

type CgroupMngr struct {
	ctrl cgroups.Cgroup
}

func NewCgroupMngr(c *pkg.UploadConfig) (*CgroupMngr, error) {
	if c.CPULimit <= 0 && c.MemoryLimit <= 0 {
		return nil, errors.New("No limit resources")
	}
	var cpu *specs.LinuxCPU
	var memory *specs.LinuxMemory
	if c.CPULimit != 0 {
		quota := int64(c.CPULimit * 1000)
		cpu = &specs.LinuxCPU{
			Quota: &quota,
		}
	}
	if c.MemoryLimit != 0 {
		memoryLimit := int64(c.MemoryLimit * 1000 * 1000)
		memory = &specs.LinuxMemory{
			Limit: &memoryLimit,
		}
	}
	ctrl, err := cgroups.New(cgroups.V2, cgroups.StaticPath("/log2s3"), &specs.LinuxResources{
		CPU:    cpu,
		Memory: memory,
	})
	if err != nil {
		return nil, err
	}
	pid := os.Getpid()
	if err := ctrl.Add(cgroups.Process{Pid: pid}); err != nil {
		return nil, err
	}
	return &CgroupMngr{ctrl}, nil
}

func (c *CgroupMngr) Close() {
	if c.ctrl != nil {
		c.ctrl.Delete()
	}
}
