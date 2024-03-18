package domaincmd

import "github.com/urfave/cli/v2"

func Command() *cli.Command {
	return &cli.Command{
		Name:        "domain",
		Description: "Command used to perform actions with available domains",
		Usage:       "Command used to perform actions with available domains",
		Subcommands: []*cli.Command{
			listCommand(),
		},
	}
}
