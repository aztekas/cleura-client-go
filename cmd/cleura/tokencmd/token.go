package tokencmd

import "github.com/urfave/cli/v2"

func Command() *cli.Command {
	return &cli.Command{
		Name:        "token",
		Description: "Command used to perform actions with Cleura API tokens",
		Subcommands: []*cli.Command{
			getCommand(),
			revokeCommand(),
			validateCommand(),
		},
	}
}
