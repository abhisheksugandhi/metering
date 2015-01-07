package main

import (
     "os/exec"
)

func main() {
     cmd := exec.Command("./bin/metering")
     cmd.Start()
}