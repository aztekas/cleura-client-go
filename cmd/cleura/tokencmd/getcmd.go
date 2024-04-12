package tokencmd

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/aztekas/cleura-client-go/cmd/cleura/common"
	"github.com/aztekas/cleura-client-go/pkg/api/cleura"
	"github.com/aztekas/cleura-client-go/pkg/configfile"
	"github.com/aztekas/cleura-client-go/pkg/util"
	"github.com/urfave/cli/v2"
)

func getCommand() *cli.Command {
	return &cli.Command{
		Name:        "get",
		Description: "Receive token from Cleura API using username and password",
		Usage:       "Receive token from Cleura API using username and password",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "username",
				Aliases: []string{"u"},
				Usage:   "Username for token request",
				EnvVars: []string{"CLEURA_API_USERNAME"},
			},
			&cli.StringFlag{
				Name:    "password",
				Aliases: []string{"p"},
				Usage:   "Password for token request.",
				EnvVars: []string{"CLEURA_API_PASSWORD"},
			},
			&cli.StringFlag{
				Name:    "api-host",
				Aliases: []string{"host"},
				Usage:   "Cleura API host",
				Value:   "https://rest.cleura.cloud",
				EnvVars: []string{"CLEURA_API_HOST"},
			},
			&cli.BoolFlag{
				Name:       "update-config",
				Usage:      "Save token to active configuration. NB: token saved in open text",
				Value:      true,
				HasBeenSet: true,
			},
			&cli.StringFlag{
				Name:  "config-path",
				Usage: "Path to configuration file. $HOME/.config/cleura/config if not set",
			},
			&cli.BoolFlag{
				Name:    "interactive",
				Usage:   "Interactive mode. Input username and password in interactive mode",
				Aliases: []string{"i"},
				Action: func(ctx *cli.Context, b bool) error {
					if ctx.String("username") != "" || ctx.String("password") != "" {
						return fmt.Errorf("error: --username (-u)/--password (-p) flags and CLEURA_API_PASSWORD/CLEURA_API_USERNAME environmental variables not supported in interactive mode")
					}
					return nil
				},
			},
			&cli.BoolFlag{
				Name:    "two-factor",
				Usage:   "Set this flag if two-factor authentication (sms) is enabled in your cleura profile ",
				Aliases: []string{"2fa"},
				Value:   false,
			},
		},
		Action: func(ctx *cli.Context) error {
			logger := common.CliLogger(ctx.String("loglevel"))
			var host, username, password string
			var client *cleura.Client
			var err error
			if ctx.Bool("interactive") {
				username, err = util.GetUserInput("username", false)
				if err != nil {
					return err
				}
				password, err = util.GetUserInput("password", true)
				if err != nil {
					return err
				}
			} else {
				username = ctx.String("username")
				password = ctx.String("password")
			}

			// Validate not empty
			if username == "" || password == "" {
				return errors.New("error: password and username must be supplied")
			}
			host = ctx.String("api-host")

			// Handle two-factor authentication
			if ctx.Bool("two-factor") {
				client, err = cleura.NewClient(&host, &username, &password, true)
				if err != nil {
					return err
				}
				input, err := util.GetUserInput("SMS Code", false)
				if err != nil {
					return err
				}
				twoFaCode, err := strconv.Atoi(input)
				if err != nil {
					return err
				}
				err = client.GetTokenWith2FA(twoFaCode)
				if err != nil {
					return err
				}
			} else {
				client, err = cleura.NewClient(&host, &username, &password, false)
				if err != nil {
					return err
				}
			}
			fmt.Printf("\nexport CLEURA_API_TOKEN=%v\nexport CLEURA_API_USERNAME=%v\nexport CLEURA_API_HOST=%v\n", client.Token, client.Auth.Username, ctx.String("api-host"))
			if ctx.Bool("update-config") {
				config, err := configfile.InitConfiguration(ctx.String("config-path"))
				if err != nil {
					return fmt.Errorf("error updating configuration file: `%s`, %w", ctx.String("config-path"), err)
				}
				err = config.SetProfileField("token", client.Token)
				if err != nil {
					return err
				}
				logger.Info("Token is updated")
			}
			return nil
		},
	}

}
