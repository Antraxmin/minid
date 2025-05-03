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

func SetupVethPair(containerID, bridgeName string) error {
	// create veth name
	vethName := fmt.Sprintf("veth-%s", containerID[:8])
	peerName := fmt.Sprintf("eth0-%s", containerID[:8])

	// create veth pair
	if err := executeCommand("ip", "link", "add", vethName, "type", "veth", "peer", "name", peerName); err != nil {
		return fmt.Errorf("failed to create veth pair: %v", err)
	}

	// connect to veth bridge
	if err := executeCommand("ip", "link", "set", vethName, "master", bridgeName); err != nil {
		return fmt.Errorf("failed to connect veth to bridge: %v", err)
	}

	// activate veth
	if err := executeCommand("ip", "link", "set", vethName, "up"); err != nil {
		return fmt.Errorf("failed to activate veth: %v", err)
	}

	// Move peer to container namespace
	if err := executeCommand("ip", "link", "set", peerName, "netns", containerID); err != nil {
		return fmt.Errorf("failed to move peer to container namespace: %v", err)
	}

	return nil
}

func ConfigureContainerNetwork(namespace, ifName string) error {
	if err := executeCommandInNetNS(namespace, "ip", "link", "set", "lo", "up"); err != nil {
		return fmt.Errorf("failed to activate lo interface: %v", err)
	}

	if err := executeCommandInNetNS(namespace, "ip", "link", "set", ifName, "name", "eth0"); err != nil {
		return fmt.Errorf("failed to change interface name: %v", err)
	}

	if err := executeCommandInNetNS(namespace, "ip", "addr", "add", "10.0.0.2/24", "dev", "eth0"); err != nil {
		return fmt.Errorf("failed to set IP address: %v", err)
	}

	if err := executeCommandInNetNS(namespace, "ip", "link", "set", "eth0", "up"); err != nil {
		return fmt.Errorf("failed to activate eth0: %v", err)
	}

	if err := executeCommandInNetNS(namespace, "ip", "route", "add", "default", "via", "10.0.0.1"); err != nil {
		return fmt.Errorf("failed to set default route: %v", err)
	}

	return nil
}
