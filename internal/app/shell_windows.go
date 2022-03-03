//go:build windows
// +build windows

package app

import (
	"log"
	"os"
	"os/exec"
)

const (
	powerShell          = "powershell"
	cmdShell            = "cmd"
	clearCommand        = "cls"
	charSelectedService = " + "
)

func cleanConsole() {
	cmd := exec.Command(powerShell, "-c", clearCommand)
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		cmd = exec.Command(cmdShell, "-c", clearCommand)
		cmd.Stdout = os.Stdout
		if err = cmd.Run(); err != nil {
			log.Fatalln(err)
		}
	}
}
