package shootcmd

import (
	"fmt"

	"github.com/aztekas/cleura-client-go/cmd/cleura/common"
	"github.com/aztekas/cleura-client-go/cmd/cleura/configcmd"
	"github.com/aztekas/cleura-client-go/pkg/api/cleura"
	"github.com/urfave/cli/v2"
)

func wakeupCommand() *cli.Command {
	commonFlags := append(common.CleuraAuthFlags(), common.LocationFlags()...)
	return &cli.Command{
		Name:        "wakeup",
		Description: "Wakeup specified shoot cluster",
		Usage:       "Wakeup specified shoot cluster",
		Before:      configcmd.TrySetConfigFromFile,
		Flags: append(
			commonFlags,
			&cli.StringFlag{
				Name:     "cluster-name",
				Category: "Basic cluster settings",
				Usage:    "Name of a cluster (Required)",
				Required: true,
			},
		),
		Action: func(ctx *cli.Context) error {
			err := common.ValidateNotEmptyString(ctx,
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
			err = client.WakeUpCluster(ctx.String("region"), ctx.String("project-id"), ctx.String("cluster-name"))
			if err != nil {
				return err
			}
			fmt.Printf("Cluster: `%s` will wake up soon.\nPlease check status with `cleura shoot list` command\n", ctx.String("cluster-name"))
			return nil
		},
	}
}
