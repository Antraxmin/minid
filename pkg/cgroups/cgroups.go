package cgroups

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

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

func (m *Manager) applyV2(cpuLimit float64, memoryLimit string) error {
	// Integration cgroup v2 path
	cgroupPath := filepath.Join("/sys/fs/cgroup", m.Path)
	if err := os.MkdirAll(cgroupPath, 0755); err != nil {
		return fmt.Errorf("Failed to create cgroup directory: %v", err)
	}

	// activate controller
	if err := writeFile(filepath.Join(cgroupPath, "cgroup.subtree_control"), "+cpu +memory"); err != nil {
		return err
	}

	// set CPU limit
	period := 100000
	quota := int(float64(period) * cpuLimit)
	if quota < 1000 {
		quota = 1000 // minimum quota
	}

	// write CPU limit file
	if err := writeFile(filepath.Join(cgroupPath, "cpu.max"),
		fmt.Sprintf("%d %d", quota, period)); err != nil {
		return err
	}

	memoryBytes, err := parseMemoryLimit(memoryLimit)
	if err != nil {
		return err
	}

	// write memory limit file
	if err := writeFile(filepath.Join(cgroupPath, "memory.max"),
		strconv.FormatInt(memoryBytes, 10)); err != nil {
		return err
	}

	return nil
}
