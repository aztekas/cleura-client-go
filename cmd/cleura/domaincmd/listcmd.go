package domaincmd

import (
	"fmt"
	"strconv"

	"github.com/aztekas/cleura-client-go/cmd/cleura/configcmd"
	"github.com/aztekas/cleura-client-go/pkg/api/cleura"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/urfave/cli/v2"
)

func listCommand() *cli.Command {
	return &cli.Command{
		Name:        "list",
		Description: "List available domains",
		Usage:       "List domains available to current user",
		Before:      configcmd.TrySetConfigFromFile,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "token",
				Aliases: []string{"t"},
				Usage:   "Token to validate",
				EnvVars: []string{"CLEURA_API_TOKEN"},
			},
			&cli.StringFlag{
				Name:    "username",
				Aliases: []string{"u"},
				Usage:   "Username token belongs to",
				EnvVars: []string{"CLEURA_API_USERNAME"},
			},
			&cli.StringFlag{
				Name:    "api-host",
				Aliases: []string{"host"},
				Usage:   "Cleura API host",
				Value:   "https://rest.cleura.cloud",
				EnvVars: []string{"CLEURA_API_HOST"},
			},
			&cli.StringFlag{
				Name:    "config-path",
				Aliases: []string{"p"},
				Usage:   "Path to configuration file. $HOME/.config/cleura/config if not set",
			},
		},
		Action: func(ctx *cli.Context) error {
			token := ctx.String("token")
			username := ctx.String("username")
			host := ctx.String("api-host")
			client, err := cleura.NewClientNoPassword(&host, &username, &token)
			if err != nil {
				return err
			}
			domains, err := client.ListDomains()
			if err != nil {
				re, ok := err.(*cleura.RequestAPIError)
				if ok {
					if re.StatusCode == 403 {
						return fmt.Errorf("error: invalid token")
					}
				}
				return err
			}
			t := table.NewWriter()
			t.SetAutoIndex(true)
			t.Style().Format.Header = text.FormatTitle
			t.AppendHeader(table.Row{"Name", "Domain Id", "Regions"})
			for _, domain := range *domains {
				var regs string
				for _, region := range domain.Area.Regions {
					regs += fmt.Sprintf("%s:%s\n", region.Region, region.Status)
				}
				t.AppendRow(table.Row{fmt.Sprintf("%s(enabled:%s,status:%s)", domain.Name, strconv.FormatBool(domain.Enabled), domain.Status), domain.Id, regs})
			}
			fmt.Println(t.Render())
			return nil
		},
	}
}
