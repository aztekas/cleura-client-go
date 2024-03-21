package utils

import "github.com/urfave/cli/v2"

func CommonFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:    "token",
			Aliases: []string{"t"},
			Usage:   "Token to validate",
			EnvVars: []string{"CLEURA_API_TOKEN"},
		},
		&cli.StringFlag{
			Name:    "username",
			Aliases: []string{"u"},
			Usage:   "Username token belongs to",
			EnvVars: []string{"CLEURA_API_USERNAME"},
		},
		&cli.StringFlag{
			Name:    "api-host",
			Aliases: []string{"host"},
			Usage:   "Cleura API host",
			Value:   "https://rest.cleura.cloud",
			EnvVars: []string{"CLEURA_API_HOST"},
		},
		&cli.StringFlag{
			Name:    "config-path",
			Aliases: []string{"p"},
			Usage:   "Path to configuration file. $HOME/.config/cleura/config if not set",
		},
	}
}
