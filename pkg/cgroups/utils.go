package cgroups

import "os"

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
