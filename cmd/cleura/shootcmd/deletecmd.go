package shootcmd

import (
	"fmt"

	"github.com/aztekas/cleura-client-go/cmd/cleura/configcmd"
	"github.com/aztekas/cleura-client-go/cmd/cleura/utils"
	"github.com/aztekas/cleura-client-go/pkg/api/cleura"
	"github.com/urfave/cli/v2"
)

func deleteCommand() *cli.Command {
	return &cli.Command{
		Name:        "delete",
		Description: "Delete a whole cluster or a workgroup in specified cluster",
		Usage:       "Delete a whole cluster or a workgroup in specified cluster",
		Before:      configcmd.TrySetConfigFromFile,
		Flags: []cli.Flag{
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
				Usage:    "Specify Cleura project to list shoot clusters in",
				Aliases:  []string{"project"},
				EnvVars:  []string{"CLEURA_API_DEFAULT_PROJECT_ID"},
			},
			&cli.StringFlag{
				Name:     "username",
				Category: "Connection settings",
				Aliases:  []string{"u"},
				Usage:    "Username token belongs to",
				EnvVars:  []string{"CLEURA_API_USERNAME"},
			},
			&cli.StringFlag{
				Name:     "token",
				Category: "Connection settings",
				Aliases:  []string{"t"},
				Usage:    "Token to validate",
				EnvVars:  []string{"CLEURA_API_TOKEN"},
			},
			&cli.StringFlag{
				Name:     "api-host",
				Category: "Connection settings",
				Aliases:  []string{"host"},
				Usage:    "Cleura API host",
				Value:    "https://rest.cleura.cloud",
				EnvVars:  []string{"CLEURA_API_HOST"},
			},
			&cli.StringFlag{
				Name:    "config-path",
				Aliases: []string{"p"},
				Usage:   "Path to a configuration file. $HOME/.config/cleura/config if not set",
			},
			&cli.BoolFlag{
				Name:  "cluster",
				Usage: "One of --cluster or --workergroup flag is Required",
				Action: func(ctx *cli.Context, b bool) error {
					if ctx.Bool("cluster") && ctx.Bool("workergroup") {
						return fmt.Errorf("error: choose one of `--cluster` or `--workergroup`")
					}
					return nil
				},
			},
			&cli.BoolFlag{
				Name:  "workergroup",
				Usage: "One of --cluster or --workergroup flag is Required",
				Action: func(ctx *cli.Context, b bool) error {
					if ctx.Bool("cluster") && ctx.Bool("workergroup") {
						return fmt.Errorf("error: choose one of `--cluster` or `--workergroup`")
					}
					return nil
				},
			},
			&cli.StringFlag{
				Name:     "cluster-name",
				Category: "Basic cluster settings",
				Usage:    "Name of a cluster (Required)",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "wg-name",
				Category: "Workergroup settings",
				Usage:    "Workergroup name",
				Action: func(ctx *cli.Context, s string) error {
					if ctx.Bool("workergroup") && ctx.String("wg-name") == "" {
						return fmt.Errorf("error: Workergroup name must be supplied via `--wg-name` flag")
					}
					if len(ctx.String("wg-name")) > 6 {
						return fmt.Errorf("error: workergroup name must be no longer than 6 characters")
					}
					return nil
				},
			},
		},
		Action: func(ctx *cli.Context) error {
			err := utils.ValidateNotEmpty(ctx,
				"token",
				"username",
				"api-host",
				"region",
				"project-id",
			)
			if err != nil {
				return err
			}
			if !ctx.Bool("cluster") && !ctx.Bool("workergroup") {
				return fmt.Errorf("error: one of `--cluster` or `--workergroup` must be set")
			}
			token := ctx.String("token")
			username := ctx.String("username")
			host := ctx.String("api-host")
			client, err := cleura.NewClientNoPassword(&host, &username, &token)
			if err != nil {
				return err
			}
			if ctx.Bool("cluster") {
				_, err := client.DeleteShootCluster(ctx.String("cluster-name"), ctx.String("region"), ctx.String("project-id"))
				if err != nil {
					re, ok := err.(*cleura.RequestAPIError)
					if ok {
						if re.StatusCode == 403 {
							return fmt.Errorf("error: invalid token")
						}
					}
					return err
				}
				fmt.Printf("Cluster: `%s` is being deleted.\nPlease check operation status with `cleura shoot list` command\n", ctx.String("cluster-name"))
			}
			return nil
		},
	}
}
