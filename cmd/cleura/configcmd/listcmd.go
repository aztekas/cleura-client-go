package configcmd

import "github.com/urfave/cli/v2"

func listCommand() *cli.Command {
	return &cli.Command{
		Name:        "list",
		Description: "List available configurations",
		Usage:       "List available configurations defined in configuration file",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config-path",
				Aliases: []string{"p"},
				Usage:   "Path to configuration file. $HOME/.config/cleura/config if not set",
			},
		},
		Action: func(ctx *cli.Context) error {
			var config *Config
			config, err := LoadConfiguration(ctx.String("config-path"))
			if err != nil {
				return err
			}
			config.PrintConfigurations()
			return nil
		},
	}
}
