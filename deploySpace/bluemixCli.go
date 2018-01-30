package main

import (
	"os"
	"os/exec"
	"strings"
)

type LiveProvisioner struct {
}

const BLUEMIXCLI = "bluemix"

func parseArgs(args ...string) []string {
	if len(args) > 1 {
		return append(strings.Split(args[0], " "), args[1])
	}
	return strings.Split(args[0], " ")
}

func (l *LiveProvisioner) CallBluemixWithArgs(args ...string) error {

	cmd := exec.Command(BLUEMIXCLI, parseArgs(args...)...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
