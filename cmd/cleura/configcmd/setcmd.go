package configcmd

import "github.com/urfave/cli/v2"

func setCommand() *cli.Command {
	return &cli.Command{
		Name:        "set",
		Description: "Set active configuration from the list of configurations defined in the configuration file",
		Usage: "Set active configuration from the list of configurations defined in the configuration file",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "name",
				Required: true,
				Aliases:  []string{"n"},
				Usage:    "Configuration name",
			},
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
			err = config.SetActive(ctx.String("name"))
			if err != nil {
				return err
			}

			return nil

		},
	}
}
