package shootcmd

import "github.com/urfave/cli/v2"

func Command() *cli.Command {
	return &cli.Command{
		Name:        "shoot",
		Description: "Command used to perform operations with shoot clusters",
		Usage:       "Command used to perform operations with shoot clusters",
		Subcommands: []*cli.Command{
			genKubeConfigCommand(),
			getKubeConfigCommand(),
			getMonitoringCredentialsCommand(),
			listCommand(),
			createCommand(),
			deleteCommand(),
			hibernateCommand(),
			wakeupCommand(),
		},
	}
}
