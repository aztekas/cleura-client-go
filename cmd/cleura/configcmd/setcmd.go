package configcmd

import (
	"github.com/aztekas/cleura-client-go/pkg/configfile"
	"github.com/urfave/cli/v2"
)

func setCommand() *cli.Command {
	return &cli.Command{
		Name:        "set",
		Description: "Set active profile from the list of profiles defined in the configuration file",
		Usage:       "Set active profile from the list of profiles defined in the configuration file",
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
			var config *configfile.Configuration
			config, err := configfile.InitConfiguration(ctx.String("config-path"))
			if err != nil {
				return err
			}
			err = config.SetActiveProfile(ctx.String("name"))
			if err != nil {
				return err
			}

			return nil

		},
	}
}
