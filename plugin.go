package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/cloudfoundry/cli/plugin"
	"github.com/tcnksm/go-input"
	"github.com/tcnksm/go-latest"
)

// Exit codes are int values that represent an exit code
// for a particular error.
const (
	ExitCodeOK    int = 0
	ExitCodeError int = 1 + iota
)

// EnvDebug is environmental variable for enabling debug
// output
const EnvDebug = "DEBUG_PLUGIN"

// Debugf prints debug output when EnvDebug is given
func Debugf(format string, args ...interface{}) {
	if env := os.Getenv(EnvDebug); len(env) != 0 {
		fmt.Fprintf(os.Stdout, "[DEBUG] "+format+"\n", args...)
	}
}

// UpdateCLI
type UpdateCLI struct {
	OutStream io.Writer
	InStream  io.Reader
}

// Run starts plugin process.
func (p *UpdateCLI) Run(cliConn plugin.CliConnection, arg []string) {
	Debugf("Run update-cli plugin")
	Debugf("Arg: %#v", arg)

	// Ensure plugin is called.
	// Plugin is also called when install/uninstall via cf command.
	// Ignore such other calls.
	if len(arg) < 1 || arg[0] != Name {
		os.Exit(ExitCodeOK)
	}

	// Read CLI context (Currently, ctx val is not used but in future it should).
	ctx, err := NewCLIContext(cliConn)
	if err != nil {
		fmt.Fprintf(p.OutStream, "Failed to read cf command context: %s\n", err)
		os.Exit(ExitCodeError)
	}

	// Call run instead of doing the work here so we can use
	// `defer` statements within the function and have them work properly.
	// (defers aren't called with os.Exit)
	os.Exit(p.run(ctx, arg[1:]))
}

func (p *UpdateCLI) run(ctx *CLIContext, args []string) int {
	var (
		check   bool
		version bool
	)

	flags := flag.NewFlagSet("plugin", flag.ContinueOnError)
	flags.SetOutput(p.OutStream)
	flags.Usage = func() {
		fmt.Fprintln(p.OutStream, p.Usage())
	}

	flags.BoolVar(&check, "check", false, "")

	flags.BoolVar(&version, "version", false, "")
	flags.BoolVar(&version, "v", false, "(Short)")

	if err := flags.Parse(args); err != nil {
		return ExitCodeError
	}

	// Show version information
	if version {
		var buf bytes.Buffer
		fmt.Fprintf(&buf, "%s v%s", Name, VersionStr())

		if len(GitCommit) != 0 {
			fmt.Fprintf(&buf, " (%s)", GitCommit)
		}

		fmt.Fprintln(p.OutStream, buf.String())
		return ExitCodeOK
	}

	Debugf("Cf version: %s", ctx.Version)
	Debugf("Cf Path to update: %s", ctx.CfPath)

	// Check version is latest or not
	githubTag := &latest.GithubTag{
		Owner:             "cloudfoundry",
		Repository:        "cli",
		FixVersionStrFunc: latest.DeleteFrontV(),
		TagFilterFunc: func(s string) bool {
			return strings.Contains(s, ".")
		},
	}

	res, err := latest.Check(githubTag, ctx.Version)
	if err != nil {
		fmt.Fprintf(p.OutStream, "Error: %s\n", err)
		return ExitCodeError
	}

	if !res.Outdated {
		fmt.Fprintf(p.OutStream, "You are using the latest version of cf cli (v%s)\n",
			ctx.Version)
		return ExitCodeOK
	}

	fmt.Fprintf(p.OutStream, "Your cf version v%s is not latest (v%s)\n",
		ctx.Version, res.Current)

	if check {
		// Just checking
		return ExitCodeOK
	}

	ui := &input.UI{
		Writer: p.OutStream,
		Reader: p.InStream,
	}
	query := "Do you want to update?: [Y/n]"
	ans, err := ui.Ask(query, &input.Options{
		Default:     "y",
		Loop:        true,
		Required:    true,
		HideDefault: true, HideOrder: true,
		ValidateFunc: func(s string) error {
			if s == "y" || s == "Y" || s == "n" {
				return nil
			}

			return fmt.Errorf("Input 'y' or 'n'")
		},
	})
	Debugf("Answer: %s", ans)

	if ans == "n" {
		fmt.Fprintf(p.OutStream, "Aboting\n")
		return ExitCodeOK
	}

	fmt.Fprintf(p.OutStream, "Start to update to v%s\n", res.Current)
	installer, err := NewInstaller(res.Current)
	if err != nil {
		fmt.Fprintf(p.OutStream, "Error: %s\n", err)
		return ExitCodeError
	}
	installer.OutStream = p.OutStream

	// On windows we can't remove an existing file or remove the running binary
	// so we download the file to cfPath.new and  move the running binary
	// to cfPath.old (deleting any existing file first) rename the downloaded
	// file to cfPath. This is the same way as heroku/heroku-cli does.
	cfPath := ctx.CfPath
	if err := installer.Install(cfPath + ".new"); err != nil {
		fmt.Fprintf(p.OutStream, "Error: %s\n", err)
		return ExitCodeError
	}

	// Rename running binary
	if err := os.Rename(cfPath, cfPath+".old"); err != nil {
		fmt.Fprintf(p.OutStream, "Error: %s\n", err)
		return ExitCodeError
	}
	defer os.Remove(cfPath + ".old")

	if err := os.Rename(cfPath+".new", cfPath); err != nil {
		fmt.Fprintf(p.OutStream, "Error: %s\n", err)
		return ExitCodeError
	}

	fmt.Fprintf(p.OutStream, "Successfully updated\n")

	return ExitCodeOK
}

func (p *UpdateCLI) GetMetadata() plugin.PluginMetadata {
	return plugin.PluginMetadata{
		Name:    Name,
		Version: Version,
		Commands: []plugin.Command{
			{
				Name:     "update-cli",
				HelpText: "Update cf cli to the latest version",
				UsageDetails: plugin.Usage{
					Usage: p.Usage(),
				},
			},
		},
	}
}

func (p *UpdateCLI) Usage() string {
	text := `cf update-cli [option]

Options:

   -check   Check current cf cli is latest or not.

`
	return text
}
