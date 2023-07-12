//go:build windows

package client

import (
	"os"
	"os/exec"
	"passwords/internal/common/logging"
)

// ClearScreen - clears screen in terminal (OS dependent)
func ClearScreen() {

	cmd := exec.Command("cmd", "/c", "cls")
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	logging.LogPrintln(err)
}
