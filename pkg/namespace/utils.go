package namespace

import (
	"os"
	"os/exec"
)

func executeCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func executeCommandInNetNS(namespace string, name string, args ...string) error {
	nsenterArgs := []string{"--net=/var/run/netns/" + namespace}
	nsenterArgs = append(nsenterArgs, name)
	nsenterArgs = append(nsenterArgs, args...)

	return executeCommand("nsenter", nsenterArgs...)
}
