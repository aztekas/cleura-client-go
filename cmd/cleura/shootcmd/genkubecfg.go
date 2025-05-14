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

func genKubeConfigCommand() *cli.Command {
	commonFlags := append(common.CleuraAuthFlags(), common.LocationFlags()...)
	return &cli.Command{
		Name:        "generate-kubeconfig",
		Description: "Get and save kubeconfig for selected shoot cluster",
		Usage:       "Get and save kubeconfig for selected shoot cluster. NB: overwrites existing kubeconfig",
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
			&cli.Int64Flag{
				Name:    "config-duration",
				Aliases: []string{"d"},
				Usage:   "How long will the generated kubeconfig be valid in seconds. Defaults to 1 day (86400 seconds)",
				Value:   86400,
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
			body, err := client.GenerateKubeConfig(
				ctx.String("gardener-domain"),
				ctx.String("region"),
				ctx.String("project-id"),
				ctx.String("cluster-name"),
				ctx.Int64("config-duration"),
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
			var configContent interface{}
			err = json.Unmarshal(body, &configContent)
			if err != nil {
				return err
			}
			content, ok := configContent.(string)
			if !ok {
				return fmt.Errorf("error: cannot assert string")
			}
			if ctx.String("output-path") != "" {
				err := configfile.WriteByteToFile(ctx.String("output-path"), []byte(content))
				return err
			} else {
				fmt.Println(content)
			}
			return nil
		},
	}
}
