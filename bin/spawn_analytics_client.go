package main

import (
	"os/exec"
)

func main() {
	cmd := exec.Command("./metering")
	cmd.Start()
}
