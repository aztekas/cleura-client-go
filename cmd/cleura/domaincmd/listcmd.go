package domaincmd

import (
	"fmt"
	"strconv"

	"github.com/aztekas/cleura-client-go/cmd/cleura/common"
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
		Flags:       common.CleuraAuthFlags(),
		Action: func(ctx *cli.Context) error {
			err := common.ValidateNotEmptyString(ctx,
				"token",
				"username",
				"api-host",
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
