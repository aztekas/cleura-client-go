package shootcmd

import (
	"encoding/json"
	"fmt"

	"github.com/aztekas/cleura-client-go/cmd/cleura/common"
	"github.com/aztekas/cleura-client-go/cmd/cleura/configcmd"
	"github.com/aztekas/cleura-client-go/pkg/api/cleura"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/urfave/cli/v2"
)

func listCommand() *cli.Command {
	commonFlags := append(common.CleuraAuthFlags(), common.LocationFlags()...)
	return &cli.Command{
		Name:        "list",
		Description: "List shoot clusters in a given project and region",
		Usage:       "List shoot clusters in a given project and region",
		Before:      configcmd.TrySetConfigFromFile,
		Flags: append(
			commonFlags,
			&cli.BoolFlag{
				Name:  "raw",
				Usage: "Output in raw json",
			},
		),
		Action: func(ctx *cli.Context) error {
			err := common.ValidateNotEmptyString(ctx,
				"token",
				"username",
				"api-host",
				"region",
				"project-id",
				"gardener-domain",
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
			clusterList, err := client.ListShootClusters(ctx.String("gardener-domain"), ctx.String("region"), ctx.String("project-id"))
			if err != nil {
				re, ok := err.(*cleura.RequestAPIError)
				if ok {
					if re.StatusCode == 403 {
						return fmt.Errorf("error: invalid token")
					}
				}
				return err
			}
			if ctx.Bool("raw") {
				raw, err := json.MarshalIndent(clusterList, "", "  ")
				if err != nil {
					return err
				}
				fmt.Println(string(raw))
			} else {
				t := table.NewWriter()
				t.SetAutoIndex(true)
				t.Style().Format.Header = text.FormatTitle
				t.AppendHeader(table.Row{"Cluster name", "Kubernetes\nVersion", "Workers", "Hibernated?", "Status", "Last operation"})
				for _, cluster := range clusterList {
					var statuses string
					for _, condition := range cluster.Status.Conditions {
						statuses += fmt.Sprintf("%s : %s\n", condition.Type, condition.Status)
					}
					var workers string
					for _, worker := range cluster.Spec.Provider.Workers {
						workers += fmt.Sprintf("name: %s\ntype: %s\nimage: %s\nimage_version: %s\nmin_nodes: %d\nmax_nodes: %d\n\n", worker.Name, worker.Machine.Type, worker.Machine.Image.Name, worker.Machine.Image.Version, worker.Minimum, worker.Maximum)
					}
					lastOperation := fmt.Sprintf("progress: %d\nstate: %s\ntype: %s\n", cluster.Status.LastOperation.Progress, cluster.Status.LastOperation.State, cluster.Status.LastOperation.Type)
					t.AppendRow(table.Row{cluster.Metadata.Name, cluster.Spec.Kubernetes.Version, workers, cluster.Status.Hibernated, statuses, lastOperation})
				}
				fmt.Printf("Shoot clusters in:\n- Project: %s\n- Region: %s\n", ctx.String("project-id"), ctx.String("region"))
				fmt.Println(t.Render())
			}
			return nil
		},
	}
}
