package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"regexp"

	"github.com/cloudfoundry/cli/plugin"
)

const (
	CfExecutable = "cf"
)

// reVersion is regular expression for cf version output
var reVersion = regexp.MustCompile(`[0-9]+\.[0-9]+\.[0-9]+`)

// CLIContext is the context which can be retrieved
// from cf command.
type CLIContext struct {
	Version string
	CfPath  string

	// Embeded because some value is needed to
	// be retrieved dynamically.
	plugin.CliConnection
}

// NewCLIContext retrieved current cf command context
func NewCLIContext(cliConn plugin.CliConnection) (*CLIContext, error) {

	version, err := cfVersion()
	if err != nil {
		return nil, err
	}

	cfPath, err := exec.LookPath(CfExecutable)
	if err != nil {
		return nil, err
	}

	return &CLIContext{
		Version:       version,
		CfPath:        cfPath,
		CliConnection: cliConn,
	}, nil
}

// cfVersion gets the cf bin version
func cfVersion() (string, error) {
	cmd := exec.Command(CfExecutable, "-v")

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf(
			"failed to get cf version string: %s\n"+
				"Output of cf command:\n%s", err, stderr.String())
	}

	// example of version string output is `cf version 6.14.0+2654a47-2015-11-18`
	verStr := stdout.String()
	return reVersion.FindString(verStr), nil
}
