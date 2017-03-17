package cgroups

import (
	"fmt"
	"os/exec"

	"github.com/taku-k/log2s3-go/pkg"
)

const (
	MAJOR_ID = 10
)

func createLimitBW(c *pkg.UploadConfig, minor int) error {
	cmd := fmt.Sprintf("tc qdisc del dev %s root", c.Device)
	_ = exec.Command(cmd).Run()
	cmd = fmt.Sprintf("tc qdisc add dev %s root handle %d: htb", c.Device, MAJOR_ID)
	if err := exec.Command(cmd).Run(); err != nil {
		return err
	}
	cmd = fmt.Sprintf("tc class add dev %s parent %d: classid %d:%d htb rate %dmbit", c.Device, MAJOR_ID, MAJOR_ID, minor, c.RateLimit)
	if err := exec.Command(cmd).Run(); err != nil {
		return err
	}
	cmd = fmt.Sprintf("tc filter add dev %s parent %d: protocol ip prio %d handle 1: cgroup", c.Device, MAJOR_ID, MAJOR_ID)
	if err := exec.Command(cmd).Run(); err != nil {
		return err
	}
	return nil
}

func deleteLimitBW(c *pkg.UploadConfig) {
	cmd := fmt.Sprintf("tc qdisc del dev %s root", c.Device)
	_ = exec.Command(cmd).Run()
}

func isEnableLimitBW(c *pkg.UploadConfig) bool {
	return c.RateLimit > 0
}
