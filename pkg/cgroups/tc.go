package cgroups

import (
	"fmt"
	"os/exec"

	"github.com/taku-k/log2s3-go/pkg"
	"strings"
)

const (
	MAJOR_ID = 10
)

func createLimitBW(c *pkg.UploadConfig, minor int) error {
	args := strings.Split(fmt.Sprintf("qdisc del dev %s root", c.Device), " ")
	_ = exec.Command("tc", args...).Run()
	args = strings.Split(fmt.Sprintf("qdisc add dev %s root handle %d: htb", c.Device, MAJOR_ID), " ")
	if err := exec.Command("tc", args...).Run(); err != nil {
		return err
	}
	args = strings.Split(fmt.Sprintf("class add dev %s parent %d: classid %d:%d htb rate %dmbit", c.Device, MAJOR_ID, MAJOR_ID, minor, c.RateLimit), " ")
	if err := exec.Command("tc", args...).Run(); err != nil {
		return err
	}
	args = strings.Split(fmt.Sprintf("filter add dev %s parent %d: protocol ip prio %d handle 1: cgroup", c.Device, MAJOR_ID, MAJOR_ID), " ")
	if err := exec.Command("tc", args...).Run(); err != nil {
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
