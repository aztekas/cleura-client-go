package configcmd

import (
	"fmt"

	"github.com/aztekas/cleura-client-go/cmd/cleura/common"
	"github.com/aztekas/cleura-client-go/pkg/configfile"
	"github.com/urfave/cli/v2"
)

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
			logger := common.CliLogger(ctx.String("loglevel"))
			logger.Debug("cleura config list")
			var config *configfile.Configuration
			config, err := configfile.InitConfiguration(ctx.String("config-path"))
			if err != nil {
				return err
			}

			fmt.Printf("Available profiles: (in %s)\n", config.Location)
			for i, prof := range config.ProfilesSlice() {
				if prof == config.GetActiveProfile() {
					prof += " (active)"
				}
				fmt.Printf("%d. %s\n", i, prof)
			}
			return nil
		},
	}
}
