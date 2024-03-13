package configcmd

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

func showActiveCommand() *cli.Command {
	return &cli.Command{
		Name:        "show-active",
		Description: "Show currently active configuration",
		Usage:       "Show currently active configuration",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "path",
				Aliases: []string{"p"},
				Usage:   "Path to configuration file. $HOME/.config/cleura/config if not set",
			},
		},
		Action: func(ctx *cli.Context) error {
			var config *Config
			config, err := LoadConfiguration(ctx.String("path"))
			if err != nil {
				return err
			}
			activeConf := config.ActiveConfig
			if activeConf == "" {
				fmt.Println("Active configuration is not set, try setting with `cleura config set --name <configuration_name>`")
			} else {
				fmt.Printf("Active configuration: `%s` (in %s)\n", activeConf, config.location)
			}
			return nil
		},
	}
}
