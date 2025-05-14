package common

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/urfave/cli/v2"
)

func CleuraAuthFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:     "token",
			Category: "Cleura auth settings",
			Aliases:  []string{"t"},
			Usage:    "Cleura API TOKEN",
			EnvVars:  []string{"CLEURA_API_TOKEN"},
		},
		&cli.StringFlag{
			Name:     "username",
			Category: "Cleura auth settings",
			Aliases:  []string{"u"},
			Usage:    "Username token belongs to",
			EnvVars:  []string{"CLEURA_API_USERNAME"},
		},
		&cli.StringFlag{
			Name:     "api-host",
			Category: "Cleura auth settings",
			Aliases:  []string{"host"},
			Usage:    "Cleura API host",
			Value:    "https://rest.cleura.cloud",
			EnvVars:  []string{"CLEURA_API_HOST"},
		},
		&cli.StringFlag{
			Name:     "config-path",
			Category: "Cleura auth settings",
			Aliases:  []string{"p"},
			Usage:    "Path to a configuration file. $HOME/.config/cleura/config if not set",
		},
	}
}

func LocationFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "gardener-domain",
			Category:    "Location settings",
			Usage:       "Specify gardener domain, defaults to 'public'",
			EnvVars:     []string{"CLEURA_API_GARDENER_DOMAIN"},
			DefaultText: "public",
		},
		&cli.StringFlag{
			Name:     "region",
			Category: "Location settings",
			Aliases:  []string{"r"},
			Usage:    "Specify region",
			EnvVars:  []string{"CLEURA_API_DEFAULT_REGION"},
		},
		&cli.StringFlag{
			Name:     "project-id",
			Category: "Location settings",
			Usage:    "Specify Cleura project",
			Aliases:  []string{"project"},
			EnvVars:  []string{"CLEURA_API_DEFAULT_PROJECT_ID"},
		},
	}
}

func ValidateNotEmptyString(ctx *cli.Context, flags ...string) error {
	for _, flag := range flags {
		if ctx.String(flag) == "" {
			return fmt.Errorf("required flag: `%s` is not set or empty", flag)
		}
	}
	return nil
}

func CliLogger(level string) *slog.Logger {
	logLevel := &slog.LevelVar{}
	switch level {
	case "warn":
		logLevel.Set(slog.LevelWarn)
	case "error":
		logLevel.Set(slog.LevelError)
	case "info":
		logLevel.Set(slog.LevelInfo)
	case "debug":
		logLevel.Set(slog.LevelDebug)
	}
	opts := &slog.HandlerOptions{
		Level: logLevel,
		//AddSource: true,
	}
	handler := slog.NewTextHandler(os.Stdout, opts)
	return slog.New(handler)
}
