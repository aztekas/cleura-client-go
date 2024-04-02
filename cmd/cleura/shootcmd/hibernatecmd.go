package shootcmd

import (
	"fmt"

	"github.com/aztekas/cleura-client-go/cmd/cleura/configcmd"
	"github.com/aztekas/cleura-client-go/cmd/cleura/utils"
	"github.com/aztekas/cleura-client-go/pkg/api/cleura"
	"github.com/urfave/cli/v2"
)

func hibernateCommand() *cli.Command {
	return &cli.Command{
		Name:        "hibernate",
		Description: "Hibernate specified shoot cluster",
		Usage:       "Hibernate specified shoot cluster",
		Before:      configcmd.TrySetConfigFromFile,
		Flags: append(
			utils.CommonFlags(),
			&cli.StringFlag{
				Name:    "region",
				Aliases: []string{"r"},
				Usage:   "Specify region",
				EnvVars: []string{"CLEURA_API_DEFAULT_REGION"},
			},
			&cli.StringFlag{
				Name:    "project-id",
				Usage:   "Specify Cleura project to list shoot clusters in",
				Aliases: []string{"project"},
				EnvVars: []string{"CLEURA_API_DEFAULT_PROJECT_ID"},
			},
			&cli.StringFlag{
				Name:     "cluster-name",
				Category: "Basic cluster settings",
				Usage:    "Name of a cluster (Required)",
				Required: true,
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
			token := ctx.String("token")
			username := ctx.String("username")
			host := ctx.String("api-host")

			client, err := cleura.NewClientNoPassword(&host, &username, &token)
			if err != nil {
				return err
			}
			err = client.HibernateCluster(ctx.String("region"), ctx.String("project-id"), ctx.String("cluster-name"))
			if err != nil {
				return err
			}
			fmt.Printf("Cluster: `%s` is being hibernated.\nPlease check status with `cleura shoot list` command\n", ctx.String("cluster-name"))
			return nil
		},
	}
}
