package tokencmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/aztekas/cleura-client-go/cmd/cleura/configcmd"
	"github.com/aztekas/cleura-client-go/pkg/api/cleura"
	"github.com/urfave/cli/v2"
)

func getCommand() *cli.Command {
	return &cli.Command{
		Name:        "get",
		Description: "Receive token from Cleura API using username and password",
		Usage: "Receive token from Cleura API using username and password",
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
				Usage:   "Password for token request",
				EnvVars: []string{"CLEURA_API_PASSWORD"},
			},
			&cli.StringFlag{
				Name:    "credentials-file",
				Aliases: []string{"c"},
				Usage:   "Path to credentials json file",
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
				Value:      false,
				HasBeenSet: true,
			},
			&cli.StringFlag{
				Name:  "path",
				Usage: "Path to configuration file. $HOME/.config/cleura/config if not set",
			},
			//Add interactive mode
			//Add two factor mode
		},
		Action: func(ctx *cli.Context) error {
			var host, username, password string
			if ctx.String("credentials-file") != "" {
				p := ctx.String("credentials-file")
				credentials := struct {
					Username string `json:"username"`
					Password string `json:"password"`
				}{}
				file, err := os.Open(filepath.Join(p))
				if err != nil {
					return err
				}
				defer file.Close()
				jsonByte, err := io.ReadAll(file)
				if err != nil {
					return err
				}
				err = json.Unmarshal(jsonByte, &credentials)
				if err != nil {
					return err
				}
				username = credentials.Username
				password = credentials.Password

			} else {
				username = ctx.String("username")
				password = ctx.String("password")
			}
			//Add option to supply u&p via console
			if username == "" || password == "" {
				return errors.New("error: password and username must be supplied")
			}
			host = ctx.String("api-host")

			client, err := cleura.NewClient(&host, &username, &password)
			if err != nil {
				return err
			}
			if ctx.Bool("update-config") {
				config, err := configcmd.LoadConfiguration(ctx.String("path"))
				if err != nil {
					return err
				}
				err = config.UpdateConfiguration("token", client.Token)
				if err != nil {
					return err
				}
				fmt.Println("Token is updated")
			}
			fmt.Printf("export CLEURA_API_TOKEN=%v\nexport CLEURA_API_USERNAME=%v\nexport CLEURA_API_HOST=%v\n", client.Token, client.Auth.Username, ctx.String("api-host"))
			return nil

		},
	}

}
