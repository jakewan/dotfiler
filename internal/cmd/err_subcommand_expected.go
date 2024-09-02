package cmd

import (
	"fmt"
	"slices"
	"strings"
)

type (
	errSubcommandExpected struct {
		oneOf []string
	}
)

func (e errSubcommandExpected) Error() string {
	return fmt.Sprintf("subcommand expected (%s)", strings.Join(e.oneOf, ","))
}

func newErrSubcommandExpected(runners []subcommand) error {
	names := []string{}
	for _, r := range runners {
		names = append(names, r.name())
	}
	slices.Sort(names)
	return errSubcommandExpected{oneOf: names}
}
