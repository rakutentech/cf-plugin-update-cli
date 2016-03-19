package main

import (
	"os"

	"github.com/cloudfoundry/cli/plugin"
)

func main() {
	updateCLI := UpdateCLI{
		OutStream: os.Stdout,
		InStream:  os.Stdin,
	}

	// Start plugin
	plugin.Start(&updateCLI)
}
