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
		Description: "Delete a cluster or a workgroup in the specified cluster",
		Usage:       "Delete a cluster or a workgroup in the specified cluster",
		Before:      configcmd.TrySetConfigFromFile,
		Flags: append(
			utils.CommonFlags(),
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
		),
		Action: func(ctx *cli.Context) error {
			err := utils.ValidateNotEmptyString(ctx,
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
			if ctx.Bool("workergroup") {
				_, err := client.DeleteWorkerGroup(ctx.String("cluster-name"), ctx.String("region"), ctx.String("project-id"), ctx.String("wg-name"))
				if err != nil {
					re, ok := err.(*cleura.RequestAPIError)
					if ok {
						if re.StatusCode == 403 {
							return fmt.Errorf("error: invalid token")
						}
					}
					return err
				}
				fmt.Printf("Workergroup: `%s` in cluster: `%s` is being deleted.\nPlease check operation status with `cleura shoot list` command\n", ctx.String("wg-name"), ctx.String("cluster-name"))
			}
			return nil
		},
	}
}
