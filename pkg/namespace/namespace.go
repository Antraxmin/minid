package namespace

import "fmt"

func CreateNetworkNamespace(name string) error {
	// execute ip netns add
	if err := executeCommand("ip", "netns", "add", name); err != nil {
		return fmt.Errorf("failed to create network namespace: %v", err)
	}
	return nil
}

func DeleteNetworkNamespace(name string) error {
	// execute ip netns delete
	if err := executeCommand("ip", "netns", "delete", name); err != nil {
		return fmt.Errorf("failed to delete network namespace")
	}
	return nil
}
