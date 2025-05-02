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

func (m *Manager) Apply(cpuLimit float64, memoryLimit string) error {
	if cgroupsV2 {
		return m.applyV2(cpuLimit, memoryLimit)
	}
	return m.applyV1(cpuLimit, memoryLimit)
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

// cgroups v1
func (m *Manager) applyV1(cpuLimit float64, memoryLimit string) error {
	cpuPath := filepath.Join("/sys/fs/cgroup/cpu", m.Path)
	if err := os.MkdirAll(cpuPath, 0755); err != nil {
		return fmt.Errorf("failed to create CPU cgroup directory: %v", err)
	}

	period := 100000
	quota := int(float64(period) * cpuLimit)
	if quota < 1000 {
		quota = 1000 // minimum 1ms
	}

	// cpu limit
	if err := writeFile(filepath.Join(cpuPath, "cpu.cfs_period_us"), strconv.Itoa(period)); err != nil {
		return err
	}
	if err := writeFile(filepath.Join(cpuPath, "cpu.cfs_quota_us"), strconv.Itoa(quota)); err != nil {
		return err
	}

	// memory limit
	memoryPath := filepath.Join("/sys/fs/cgroup/memory", m.Path)
	if err := os.MkdirAll(memoryPath, 0755); err != nil {
		return fmt.Errorf("failed to create memory cgroup directory: %v", err)
	}

	memoryBytes, err := parseMemoryLimit(memoryLimit)
	if err != nil {
		return err
	}

	if err := writeFile(filepath.Join(memoryPath, "memory.limit_in_bytes"), strconv.FormatInt(memoryBytes, 10)); err != nil {
		return err
	}
	return err
}

func (m *Manager) AddProcess(pid int) error {
	if cgroupsV2 {
		cgroupPath := filepath.Join("/sys/fs/cgroup", m.Path)
		return writeFile(filepath.Join(cgroupPath, "cgroup.procs"), strconv.Itoa(pid))
	}

	// v1
	subsystems := []string{"cpu", "memory"}
	for _, sys := range subsystems {
		path := filepath.Join("/sys/fs/cgroup", sys, m.Path)
		if err := writeFile(filepath.Join(path, "cgroup.procs"), strconv.Itoa(pid)); err != nil {
			return err
		}
	}

	return nil
}

func (m *Manager) Destroy() error {
	if cgroupsV2 {
		return os.RemoveAll(filepath.Join("/sys/fs/cgroup", m.Path))
	}

	// v1
	subsystems := []string{"cpu", "memory"}
	for _, sys := range subsystems {
		path := filepath.Join("/sys/fs/cgroup", sys, m.Path)
		if err := os.RemoveAll(path); err != nil {
			return err
		}
	}

	return nil
}
