package cmd

import "context"

type (
	subcommand interface {
		name() string
		init([]string) (subcommandRunner, error)
	}

	subcommandRunner interface {
		subcommand
		run(context.Context, Dependencies) error
	}
)
