//go:build !windows
// +build !windows

package app

import (
	"log"
	"os"
	"os/exec"
)

const (
	clearCommand          = "clear"
)

func cleanConsole() {
	cmd := exec.Command(clearCommand)
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		log.Fatalln(err)
	}
}
