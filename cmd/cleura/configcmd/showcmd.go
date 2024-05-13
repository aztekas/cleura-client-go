package configcmd

import (
	"fmt"

	"github.com/aztekas/cleura-client-go/pkg/configfile"
	"github.com/urfave/cli/v2"
	"golang.org/x/exp/maps"
)

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
				Usage:   "Configuration name within config file. Choose currently active by default",
			},
		},
		Action: func(ctx *cli.Context) error {
			var config *configfile.Configuration
			config, err := configfile.InitConfiguration(ctx.String("config-path"))
			if err != nil {
				return err
			}
			profileMap, err := config.GetProfileMap(ctx.String("name"))
			if err != nil {
				return err
			}
			fmt.Printf("Details for profile: `%s`\n", config.GetActiveProfile())
			for _, key := range maps.Keys(profileMap) {
				if key == "token" && profileMap[key] != "" {
					fmt.Printf("%-10s: ****hidden****\n", "token")
				} else {
					fmt.Printf("%-10s: %s\n", key, profileMap[key])
				}

			}
			return nil
		},
	}
}
