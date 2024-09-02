package cmd

import "fmt"

type (
	errUnexpectedSubcommand struct {
		sc string
	}
)

func (e errUnexpectedSubcommand) Error() string {
	return fmt.Sprintf("unexpected subcommand %s", e.sc)
}

func newErrUnexpectedSubcommand(sc string) error {
	return errUnexpectedSubcommand{
		sc: sc,
	}
}
