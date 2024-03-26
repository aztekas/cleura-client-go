package shootcmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/aztekas/cleura-client-go/cmd/cleura/configcmd"
	"github.com/aztekas/cleura-client-go/cmd/cleura/utils"
	"github.com/aztekas/cleura-client-go/pkg/api/cleura"
	"github.com/urfave/cli/v2"
)

func genKubeConfigCommand() *cli.Command {
	return &cli.Command{
		Name:        "generate-kubeconfig",
		Description: "Get and save kubeconfig for selected shoot cluster",
		Usage:       "Get and save kubeconfig for selected shoot cluster. NB: overwrites existing kubeconfig",
		Before:      configcmd.TrySetConfigFromFile,
		Flags: append(
			utils.CommonFlags(),
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
			&cli.StringFlag{
				Name:    "region",
				Aliases: []string{"r"},
				Usage:   "Openstack region. Try \"cleura domains list\" command for available regions in your domain",
				EnvVars: []string{"CLEURA_API_DEFAULT_REGION"},
			},
			&cli.StringFlag{
				Name:    "project-id",
				Aliases: []string{"p"},
				Usage:   "Openstack project id. Try \"cleura project list\" command for the list of available projects",
				EnvVars: []string{"CLEURA_API_DEFAULT_PROJECT_ID"},
			},
			&cli.Int64Flag{
				Name:    "config-duration",
				Aliases: []string{"d"},
				Usage:   "How long will the generated kubeconfig be valid in seconds. Defaults to 1 day (86400 seconds)",
				Value:   86400,
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
			body, err := client.GenerateKubeConfig(
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
			if ctx.String("output-path") != "" {
				f, err := os.OpenFile(ctx.String("output-path"), os.O_RDWR|os.O_CREATE, 0600)
				if err != nil {
					return err
				}
				if _, err := f.Write([]byte(configContent.(string))); err != nil {
					f.Close() // ignore error; Write error takes precedence
					return err
				}
				if err := f.Close(); err != nil {
					return err
				}
			} else {
				content, ok := configContent.(string)
				if !ok {
					return fmt.Errorf("error: cannot assert string")
				}
				fmt.Println(content)
			}
			return nil
		},
	}
}
