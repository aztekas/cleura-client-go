package projectcmd

import "github.com/urfave/cli/v2"

func Command() *cli.Command {
	return &cli.Command{
		Name:        "project",
		Description: "Command used to perform operations with projects",
		Subcommands: []*cli.Command{
			listCommand(),
		},
	}
}
