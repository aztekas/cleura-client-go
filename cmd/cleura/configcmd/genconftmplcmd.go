package configcmd

import (
	"github.com/aztekas/cleura-client-go/pkg/configfile"
	"github.com/urfave/cli/v2"
)

func generateConfigTemplateCommand() *cli.Command {
	return &cli.Command{
		Name:        "generate-template",
		Description: "Generate configuration file template on the given path",
		Usage:       "Generate configuration file template on the given path",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "output-file",
				Aliases: []string{"o"},
				Usage:   "Path to configuration file. $HOME/.config/cleura/config if not set. NB: Overwrites existing if found",
			},
		},
		Action: func(ctx *cli.Context) error {
			err := configfile.CreateConfigTemplateFile(ctx.String("output-file"))
			if err != nil {
				return err
			}
			return nil
		},
	}
}
