package cgroups

import "os"

var cgroupsV2 = false

func init() {
	if _, err := os.Stat("/sys/fs/cgroup/cgroup.controllers"); err == nil {
		cgroupsV2 = true
	}
}

type Manager struct {
	Path string // cgroup path
}

// create new cgroups manager
func NewManager(path string) *Manager {
	return &Manager{
		Path: path,
	}
}
