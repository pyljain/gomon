package main

import (
	"os"
	"os/exec"
	"syscall"
)

func startProcess(pathToMainFile string) (*exec.Cmd, error) {
	cmd := exec.Command("go", "run", pathToMainFile)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	err := cmd.Start()
	if err != nil {
		return nil, err
	}
	// cmd.Wait()

	return cmd, nil
}
