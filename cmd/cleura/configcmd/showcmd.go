package configcmd

import "github.com/urfave/cli/v2"

func showCommand() *cli.Command {
	return &cli.Command{
		Name:        "show",
		Description: "Show configured settings for given configuration",
		Usage:       "Show configured settings for given configuration",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config-path",
				Aliases: []string{"p"},
				Usage:   "Path to configuration file. $HOME/.config/cleura/config if not set",
			},
			&cli.StringFlag{
				Name:    "name",
				Aliases: []string{"n"},
				Usage:   "Configuration name withing config file. Choose currently active by default",
			},
		},
		Action: func(ctx *cli.Context) error {
			var config *Config
			config, err := LoadConfiguration(ctx.String("config-path"))
			if err != nil {
				return err
			}
			if ctx.String("name") == "" {
				err := config.PrintConfigurationContent(config.ActiveConfig)
				if err != nil {
					return err
				}
			} else {
				err := config.PrintConfigurationContent(ctx.String("name"))
				if err != nil {
					return err
				}
			}
			return nil
		},
	}
}
