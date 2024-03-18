package configcmd

import (
	"os"

	"github.com/aztekas/cleura-client-go/cmd/cleura/utils"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v2"
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
			config := &Config{
				ActiveConfig: "default",
				Configurations: map[string]Configuration{
					"default": {},
				},
			}
			templateByte, err := yaml.Marshal(config)
			if err != nil {
				return err
			}
			path, err := utils.ChoosePath(ctx.String("output-file"))
			if err != nil {
				return err
			}
			f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0600)
			if err != nil {
				return err
			}
			if _, err := f.Write(templateByte); err != nil {
				f.Close() // ignore error; Write error takes precedence
				return err
			}
			if err := f.Close(); err != nil {
				return err
			}
			return nil
		},
	}
}
