package cmd

import (
	"context"
	"errors"

	"github.com/fatih/color"
)

const (
	flagUsageManifest           string = "Path to a dotfiler.yml manifest file"
	flagUsageDestinationRootDir string = "The destination root directory. If not provided, the user's home directory will be used."
)

const (
	manifestFileName string = "dotfiler.yml"
)

var (
	colorConfirmation = color.New(color.FgCyan)
	colorWarning      = color.New(color.FgYellow)
	colorError        = color.New(color.FgRed)
)

var (
	ErrInternal = errors.New("Error handled internally")
)

type Dependencies interface {
	GetArch() string
	GetOS() string
	GetHomeDirectory() (string, error)
	GetSingleKey() (ch rune, err error)
}

// Execute routes the command to the appropriate handler.
func Execute(ctx context.Context, args []string, deps Dependencies) error {
	return rootHandler(ctx, args, deps)
}

func rootHandler(ctx context.Context, args []string, deps Dependencies) error {
	subcommands := []subcommand{
		&cmdFiles{},
	}
	if len(args) < 1 {
		return newErrSubcommandExpected(subcommands)
	}
	subcommand := args[0]
	for _, r := range subcommands {
		if r.name() == subcommand {
			if runner, err := r.init(args[1:]); err != nil {
				return err
			} else if err := runner.run(ctx, deps); err != nil {
				return err
			}
			return nil
		}
	}
	return newErrUnexpectedSubcommand(subcommand)
}

func processSubcommand(args []string, childSubcommands []subcommand) (subcommandRunner, error) {
	subcommand := args[0]
	for _, r := range childSubcommands {
		if r.name() == subcommand {
			if runner, err := r.init(args[1:]); err != nil {
				return nil, err
			} else {
				return runner, nil
			}
		}
	}
	return nil, newErrUnexpectedSubcommand(subcommand)
}
