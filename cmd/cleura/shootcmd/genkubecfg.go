package shootcmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/aztekas/cleura-client-go/cmd/cleura/configcmd"
	"github.com/aztekas/cleura-client-go/pkg/api/cleura"
	"github.com/urfave/cli/v2"
)

func genKubeConfigCommand() *cli.Command {
	return &cli.Command{
		Name:        "generate-kubeconfig",
		Description: "Get and save kubeconfig for selected shoot cluster",
		Usage:       "Get and save kubeconfig for selected shoot cluster. NB: overwrites existing kubeconfig",
		Before:      configcmd.TrySetConfigFromFile,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "output-path",
				Aliases: []string{"o"},
				Usage:   "Specify path with filename to store kubeconfig",
			},
			&cli.StringFlag{
				Name:     "username",
				Required: true,
				Aliases:  []string{"u"},
				Usage:    "Username token belongs to",
				EnvVars:  []string{"CLEURA_API_USERNAME"},
			},
			&cli.StringFlag{
				Name:     "token",
				Required: true,
				Aliases:  []string{"t"},
				Usage:    "Token to validate",
				EnvVars:  []string{"CLEURA_API_TOKEN"},
			},
			&cli.StringFlag{
				Name:    "api-host",
				Aliases: []string{"host"},
				Usage:   "Cleura API host",
				Value:   "https://rest.cleura.cloud",
				EnvVars: []string{"CLEURA_API_HOST"},
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
				Name:     "project-id",
				Required: true,
				Aliases:  []string{"p"},
				Usage:    "Openstack project id. Try \"cleura project list\" command for the list of available projects",
				EnvVars:  []string{"CLEURA_API_DEFAULT_PROJECT_ID"},
			},
			&cli.Int64Flag{
				Name:    "config-duration",
				Aliases: []string{"d"},
				Usage:   "How long will the generated kubeconfig be valid in seconds. Defaults to 1 day (86400 seconds)",
				Value:   86400,
			},
		},
		Action: func(ctx *cli.Context) error {
			token := ctx.String("token")
			username := ctx.String("username")
			host := ctx.String("api-host")
			clusterName := ctx.String("cluster-name")
			clusterRegion := ctx.String("region")
			clusterProjectId := ctx.String("project-id")
			configDuration := ctx.Int64("config-duration")
			outputPath := ctx.String("output-path")
			client, err := cleura.NewClientNoPassword(&host, &username, &token)
			if err != nil {
				return err
			}
			body, err := client.GenerateKubeConfig(clusterRegion, clusterProjectId, clusterName, configDuration)
			if err != nil {
				return err
			}
			var configContent interface{}
			err = json.Unmarshal(body, &configContent)
			if err != nil {
				return err
			}
			if outputPath != "" {
				f, err := os.OpenFile(outputPath, os.O_RDWR|os.O_CREATE, 0600)
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
				fmt.Println(configContent.(string))
			}
			return nil
		},
	}
}
