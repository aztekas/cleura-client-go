package shootcmd

import (
	"fmt"
	"strings"

	"github.com/aztekas/cleura-client-go/cmd/cleura/common"
	"github.com/aztekas/cleura-client-go/cmd/cleura/configcmd"
	"github.com/aztekas/cleura-client-go/pkg/api/cleura"
	"github.com/urfave/cli/v2"
)

func createCommand() *cli.Command {
	commonFlags := append(common.CleuraAuthFlags(), common.LocationFlags()...)
	return &cli.Command{
		Name:        "create",
		Description: "Create shoot cluster or add a workergroup",
		Usage:       "Create shoot cluster or add a workergroup",
		Before:      configcmd.TrySetConfigFromFile,
		Flags: append(
			commonFlags,
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
			&cli.BoolFlag{
				Name:     "enable-ha-control-plane",
				Category: "Basic cluster settings",
				Usage:    "Enable HA for control plane",
				Value:    false,
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
			&cli.StringSliceFlag{
				Name:     "wg-annotation",
				Category: "Workergroup settings",
				Usage:    "Custom annotations for workergroup, can be set multiple times. supplied as key=value",
				Action: func(ctx *cli.Context, s []string) error {
					for _, param := range s {
						if !strings.Contains(param, "=") {
							return fmt.Errorf("error: Annotations must be supplied as key=value")
						}
					}
					return nil
				},
			},
			&cli.StringSliceFlag{
				Name:     "wg-label",
				Category: "Workergroup settings",
				Usage:    "Custom labels for workergroup, can be set multiple times. supplied as key=value",
				Action: func(ctx *cli.Context, s []string) error {
					for _, param := range s {
						if !strings.Contains(param, "=") {
							return fmt.Errorf("error: Labels must be supplied as key=value")
						}
					}
					return nil
				},
			},
			&cli.StringSliceFlag{
				Name:     "wg-taint",
				Category: "Workergroup settings",
				Usage:    "Custom taints for workergroup, can be set multiple times. supplied as key=value",
				Action: func(ctx *cli.Context, s []string) error {
					for _, param := range s {
						if !strings.Contains(param, "=") || !strings.Contains(param, ":") {
							return fmt.Errorf("error: Taints must be supplied as key=value:effect")
						}
					}
					return nil
				},
			},
			&cli.StringSliceFlag{
				Name:     "wg-zone",
				Category: "Workergroup settings",
				Usage:    "Set compute zone for workergroup",
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
			&cli.StringFlag{
				Name:     "maintenance-start",
				Category: "Maintenance settings",
				Usage:    "Maintenance schedule, Start in format 010000+0000, (e.g. 04:00:00 UTC)",
				Action: func(ctx *cli.Context, s string) error {
					if ctx.String("maintenance-end") == "" {
						return fmt.Errorf("error: both maintenance -start and -end flags must be supplied")
					}
					return nil
				},
			},
			&cli.StringFlag{
				Name:     "maintenance-end",
				Category: "Maintenance settings",
				Usage:    "Maintenance schedule, End in format 040000+0000, (e.g. 04:00:00 UTC)",
				Action: func(ctx *cli.Context, s string) error {
					if ctx.String("maintenance-start") == "" {
						return fmt.Errorf("error: both maintenance -start and -end flags must be supplied")
					}
					return nil
				},
			},
			&cli.BoolFlag{
				Name:     "allow-k8s-autoupdate",
				Category: "Maintenance settings",
				Usage:    "Toggle if automatic updates of kubernetes is allowed",
				Value:    true,
			},
			&cli.BoolFlag{
				Name:     "allow-worker-image-autoupdate",
				Category: "Maintenance settings",
				Usage:    "Toggle if automatic updates of kubernetes is allowed",
				Value:    true,
			},
			&cli.StringFlag{
				Name:     "network-id",
				Category: "Network settings",
				Usage:    "ID of an existing OpenStack network to attach workers on",
				Action: func(ctx *cli.Context, s string) error {
					if ctx.String("router-id") == "" {
						return fmt.Errorf("error: both network-id and router-id flags must be set")
					}
					return nil
				},
			},
			&cli.StringFlag{
				Name:     "router-id",
				Category: "Network settings",
				Usage:    "ID of an existing OpenStack router managing the worker node subnet",
				Action: func(ctx *cli.Context, s string) error {
					if ctx.String("network-id") == "" {
						return fmt.Errorf("error: both network-id and router-id flags must be set")
					}
					return nil
				},
			},
			&cli.StringFlag{
				Name:     "subnet-cidr",
				Category: "Network settings",
				Usage:    "Custom subnet CIDR to use for worker nodes",
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
				_, err := client.CreateShootCluster(ctx.String("gardener-domain"), ctx.String("region"), ctx.String("project-id"), clusterReq)
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
				resp, err := client.AddWorkerGroup(ctx.String("gardener-domain"), ctx.String("cluster-name"), ctx.String("region"), ctx.String("project-id"), wgReq)
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
			EnableHaControlPlane: ctx.Bool("enable-ha-control-plane"),
			Maintenance: &cleura.MaintenanceDetails{
				AutoUpdate: &cleura.AutoUpdateDetails{
					KubernetesVersion:   ctx.Bool("allow-k8s-autoupdate"),
					MachineImageVersion: ctx.Bool("allow-worker-image-autoupdate"),
				},
				TimeWindow: &cleura.TimeWindowDetails{
					Begin: "000000+0000",
					End:   "010000+0000",
				},
			},
			Provider: &cleura.ProviderDetailsRequest{
				InfrastructureConfig: cleura.InfrastructureConfigDetails{
					FloatingPoolName: "ext-net",
				},
				Workers: []cleura.WorkerRequest{
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
						Annotations: stringSliceToKeyValueSlice(ctx.StringSlice("wg-annotation")),
						Labels:      stringSliceToKeyValueSlice(ctx.StringSlice("wg-labels")),
						Taints:      taintStringSliceToKeyValueSlice(ctx.StringSlice("wg-taints")),
						Zones:       ctx.StringSlice("wg-labels"),
					},
				},
			},
		},
	}
	// Give name to a worker group if provided, otherwise it will be generated automatically
	if ctx.String("wg-name") != "" {
		clusterReq.Shoot.Provider.Workers[0].Name = ctx.String("wg-name")
	}

	if ctx.String("maintenance-start") != "" && ctx.String("maintenance-end") != "" {
		clusterReq.Shoot.Maintenance.TimeWindow = &cleura.TimeWindowDetails{
			Begin: ctx.String("maintenance-start"),
			End:   ctx.String("maintenance-end"),
		}
	}

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

	if ctx.String("network-id") != "" && ctx.String("router-id") != "" {
		clusterReq.Shoot.Provider.InfrastructureConfig.Networks = &cleura.WorkerNetwork{
			Id: ctx.String("network-id"),
			Router: cleura.Router{
				Id: ctx.String("router-id"),
			},
		}
	}

	if cidr := ctx.String("subnet-cidr"); cidr != "" {
		if clusterReq.Shoot.Provider.InfrastructureConfig.Networks == nil {
			clusterReq.Shoot.Provider.InfrastructureConfig.Networks = &cleura.WorkerNetwork{}
		}
		clusterReq.Shoot.Provider.InfrastructureConfig.Networks.WorkersCIDR = cidr
	}

	return clusterReq
}

func stringSliceToKeyValueSlice(stringslice []string) []cleura.KeyValuePair {
	kv := make([]cleura.KeyValuePair, 0)
	if len(stringslice) > 0 {
		for _, annotation := range stringslice {
			keyValue := strings.SplitN(annotation, "=", 2)
			if len(keyValue) != 2 {
				panic(fmt.Errorf("expected annotation have a equal sign as delimited, got %s", annotation))
			}
			kv = append(kv, cleura.KeyValuePair{
				Key:   keyValue[0],
				Value: keyValue[1],
			})
		}
	}
	return kv
}

func taintStringSliceToKeyValueSlice(stringslice []string) []cleura.Taint {
	kv := make([]cleura.Taint, 0)
	if len(stringslice) > 0 {
		for _, annotation := range stringslice {
			keyValue := strings.SplitN(annotation, "=", 2)
			if len(keyValue) != 2 {
				panic(fmt.Errorf("expected annotation have a equal sign as delimiter for key=value, got %s", annotation))
			}
			valueEffect := strings.SplitN(keyValue[1], ":", 2)
			if len(keyValue) != 2 {
				panic(fmt.Errorf("expected annotation have a colon sign as delimiter for value:effect, got %s", annotation))
			}
			kv = append(kv, cleura.Taint{
				Key:    keyValue[0],
				Value:  valueEffect[0],
				Effect: valueEffect[1],
			})
		}
	}
	return kv
}

func generateWorkerGroupRequest(ctx *cli.Context) cleura.WorkerGroupRequest {

	wgReq := cleura.WorkerGroupRequest{
		Worker: cleura.WorkerRequest{
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
			Annotations: stringSliceToKeyValueSlice(ctx.StringSlice("wg-annotation")),
			Labels:      stringSliceToKeyValueSlice(ctx.StringSlice("wg-label")),
			Taints:      []cleura.Taint{},
			Zones:       ctx.StringSlice("wg-zones"),
		},
	}
	if ctx.String("wg-name") != "" {
		wgReq.Worker.Name = ctx.String("wg-name")
	}
	fmt.Println(wgReq)
	return wgReq
}
