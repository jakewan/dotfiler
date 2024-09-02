package cmd

type cmdFiles struct{}

func (cmd *cmdFiles) init(args []string) (subcommandRunner, error) {
	runners := []subcommand{
		newCmdFilesUpdate(),
	}
	if len(args) < 1 {
		return nil, newErrSubcommandExpected(runners)
	}
	return processSubcommand(args, runners)
}

func (cmd *cmdFiles) name() string {
	return "files"
}
