package shootcmd

import (
	"encoding/json"
	"fmt"

	"github.com/aztekas/cleura-client-go/cmd/cleura/configcmd"
	"github.com/aztekas/cleura-client-go/cmd/cleura/utils"
	"github.com/aztekas/cleura-client-go/pkg/api/cleura"
	"github.com/urfave/cli/v2"
)

func createCommand() *cli.Command {
	return &cli.Command{
		Name:        "create",
		Description: "Create shoot cluster or add a workergroup",
		Usage:       "Create shoot cluster or add a workergroup",
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
				Usage:    "Specify Cleura project id",
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
			&cli.StringFlag{
				Name:     "k8s-version",
				Category: "Basic cluster settings",
				Usage:    "Supported Kubernetes version",
				Value:    "1.28.8",
			},
			&cli.IntFlag{
				Name:     "wg-min",
				Category: "Workergroup settings",
				Usage:    "Autoscaler min nodes",
				Value:    2,
			},
			&cli.IntFlag{
				Name:     "wg-max",
				Category: "Workergroup settings",
				Usage:    "Autoscaler max nodes",
				Value:    3,
			},
			&cli.StringFlag{
				Name:     "wg-type",
				Category: "Workergroup settings",
				Usage:    "Workergroup machine type",
				Value:    "b.2c4gb",
			},
			&cli.StringFlag{
				Name:     "wg-image-name",
				Category: "Workergroup settings",
				Usage:    "Workergroup image name",
				Value:    "gardenlinux",
			},
			&cli.StringFlag{
				Name:     "wg-image-version",
				Category: "Workergroup settings",
				Usage:    "Workergroup image version",
				Value:    "1312.3.0",
			},
			&cli.StringFlag{
				Name:     "wg-volume-size",
				Category: "Workergroup settings",
				Usage:    "Workergroup machine volume size",
				Value:    "50Gi",
			},
			&cli.StringFlag{
				Name:     "hibernation-start",
				Category: "Hibernation settings",
				Usage:    "Hibernation schedule, Start in cron format (ex: \"00 18 * * 1,2,3,4,5\")",
				Action: func(ctx *cli.Context, s string) error {
					if ctx.String("hibernation-end") == "" {
						return fmt.Errorf("error: both hibernation -start and -end flags must be supplied")
					}
					return nil
				},
			},
			&cli.StringFlag{
				Name:     "hibernation-end",
				Category: "Hibernation settings",
				Usage:    "Hibernation schedule, End in cron format (ex: \"00 08 * * 1,2,3,4,5\")",
				Action: func(ctx *cli.Context, s string) error {
					if ctx.String("hibernation-start") == "" {
						return fmt.Errorf("error: both hibernation -start and -end flags must be supplied")
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
				clusterReq := generateShootClusterRequest(ctx)
				body, _ := json.MarshalIndent(clusterReq, "", " ")
				fmt.Printf("%s", string(body))
				_, err := client.CreateShootCluster(ctx.String("region"), ctx.String("project-id"), clusterReq)
				if err != nil {
					re, ok := err.(*cleura.RequestAPIError)
					if ok {
						if re.StatusCode == 403 {
							return fmt.Errorf("error: invalid token")
						}
					}
					return err
				}
				fmt.Printf("Cluster: `%s` is being created.\nPlease check status with `cleura shoot list` command\n", ctx.String("cluster-name"))

			}
			if ctx.Bool("workergroup") {
				wgReq := generateWorkerGroupRequest(ctx)
				// body, err := json.MarshalIndent(wgReq, "", " ")
				// if err != nil {
				// 	return err
				// }
				// fmt.Printf("%s", string(body))
				resp, err := client.AddWorkerGroup(ctx.String("cluster-name"), ctx.String("region"), ctx.String("project-id"), wgReq)
				if err != nil {
					re, ok := err.(*cleura.RequestAPIError)
					if ok {
						if re.StatusCode == 403 {
							return fmt.Errorf("error: invalid token")
						}
					}
					return err
				}
				fmt.Printf("New workgroup is being added to the cluster `%s`.\nPlease check status with `cleura shoot list` command\n", resp.Metadata.Name)
			}
			return nil
		},
	}
}

func generateShootClusterRequest(ctx *cli.Context) cleura.ShootClusterRequest {

	clusterReq := cleura.ShootClusterRequest{
		Shoot: cleura.ShootClusterRequestConfig{
			Name: ctx.String("cluster-name"),
			KubernetesVersion: &cleura.K8sVersion{
				Version: ctx.String("k8s-version"),
			},
			Provider: &cleura.ProviderDetails{
				InfrastructureConfig: cleura.InfrastructureConfigDetails{
					FloatingPoolName: "ext-net",
				},
				Workers: []cleura.Worker{
					{
						Minimum: int16(ctx.Int("wg-min")),
						Maximum: int16(ctx.Int("wg-max")),
						Machine: cleura.MachineDetails{
							Type: ctx.String("wg-type"),
							Image: cleura.ImageDetails{
								Name:    ctx.String("wg-image-name"),
								Version: ctx.String("wg-image-version"),
							},
						},
						Volume: cleura.VolumeDetails{
							Size: ctx.String("wg-volume-size"),
						},
					},
				},
			},
		},
	}
	// Give name to a worker group if provided, otherwise it will be generated automatically
	if ctx.String("wg-name") != "" {
		clusterReq.Shoot.Provider.Workers[0].Name = ctx.String("wg-name")
	}
	if ctx.String("hibernation-start") != "" && ctx.String("hibernation-end") != "" {
		clusterReq.Shoot.Hibernation = &cleura.HibernationSchedules{
			HibernationSchedules: []cleura.HibernationSchedule{
				{
					Start: ctx.String("hibernation-start"),
					End:   ctx.String("hibernation-end"),
				},
			},
		}
	}
	return clusterReq
}

func generateWorkerGroupRequest(ctx *cli.Context) cleura.WorkerGroupRequest {
	wgReq := cleura.WorkerGroupRequest{
		Worker: cleura.Worker{
			Minimum: int16(ctx.Int("wg-min")),
			Maximum: int16(ctx.Int("wg-max")),
			Machine: cleura.MachineDetails{
				Type: ctx.String("wg-type"),
				Image: cleura.ImageDetails{
					Name:    ctx.String("wg-image-name"),
					Version: ctx.String("wg-image-version"),
				},
			},
			Volume: cleura.VolumeDetails{
				Size: ctx.String("wg-volume-size"),
			},
		},
	}
	if ctx.String("wg-name") != "" {
		wgReq.Worker.Name = ctx.String("wg-name")
	}
	return wgReq
}
