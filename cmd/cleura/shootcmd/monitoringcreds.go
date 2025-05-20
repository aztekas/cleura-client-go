package shootcmd

import (
	"encoding/json"
	"fmt"

	"github.com/aztekas/cleura-client-go/cmd/cleura/common"
	"github.com/aztekas/cleura-client-go/cmd/cleura/configcmd"
	"github.com/aztekas/cleura-client-go/pkg/api/cleura"
	"github.com/aztekas/cleura-client-go/pkg/configfile"
	"github.com/urfave/cli/v2"
)

func getMonitoringCredentialsCommand() *cli.Command {
	commonFlags := append(common.CleuraAuthFlags(), common.LocationFlags()...)
	return &cli.Command{
		Name:        "get-monitoring-creds",
		Description: "Get monitoring credentials for selected shoot cluster",
		Usage:       "Get monitoring credentials for selected shoot cluster. NB: overwrites existing output file",
		Before:      configcmd.TrySetConfigFromFile,
		Flags: append(
			commonFlags,
			&cli.StringFlag{
				Name:    "output-path",
				Aliases: []string{"o"},
				Usage:   "Specify path with filename to store kubeconfig. Print to stdout if not set",
			},
			&cli.StringFlag{
				Name:     "cluster-name",
				Required: true,
				Aliases:  []string{"n"},
				Usage:    "Shoot cluster name",
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
			body, err := client.GetMonitoringCredentials(
				ctx.String("gardener-domain"),
				ctx.String("region"),
				ctx.String("project-id"),
				ctx.String("cluster-name"),
			)
			if err != nil {
				re, ok := err.(*cleura.RequestAPIError)
				if ok {
					if re.StatusCode == 403 {
						return fmt.Errorf("error: invalid token")
					}
				}
				return err
			}
			var content interface{}
			err = json.Unmarshal(body, &content)
			if err != nil {
				return err
			}

			prettyJson, err := json.MarshalIndent(content, "", "  ")
			if err != nil {
				return err
			}

			if ctx.String("output-path") != "" {
				err := configfile.WriteByteToFile(ctx.String("output-path"), prettyJson)
				return err
			} else {
				fmt.Println(string(prettyJson))
			}
			return nil
		},
	}
}
