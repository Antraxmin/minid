package cgroups

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func readFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func writeFile(path, data string) error {
	return os.WriteFile(path, []byte(data), 0644)
}

func parseMemoryLimit(limit string) (int64, error) {
	limit = strings.ToLower(limit)

	var multiplier int64 = 1
	if strings.HasSuffix(limit, "k") {
		multiplier = 1024
		limit = limit[:len(limit)-1]
	} else if strings.HasSuffix(limit, "m") {
		multiplier = 1024 * 1024
		limit = limit[:len(limit)-1]
	} else if strings.HasSuffix(limit, "g") {
		multiplier = 1024 * 1024 * 1024
		limit = limit[:len(limit)-1]
	}

	value, err := strconv.ParseInt(limit, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("Invalid memory limit format: %s", limit)
	}

	return value * multiplier, nil
}
