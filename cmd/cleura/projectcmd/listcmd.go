package projectcmd

import (
	"fmt"
	"strconv"

	"github.com/aztekas/cleura-client-go/cmd/cleura/configcmd"
	"github.com/aztekas/cleura-client-go/cmd/cleura/utils"
	"github.com/aztekas/cleura-client-go/pkg/api/cleura"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/urfave/cli/v2"
)

func listCommand() *cli.Command {
	return &cli.Command{
		Name:   "list",
		Usage:  "List projects in the defined domain",
		Before: configcmd.TrySetConfigFromFile,
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
				Name:    "domain-id",
				Aliases: []string{"d"},
				Usage:   "Openstack domain id. Try \"cleura domain list\" for the list of available domains",
				EnvVars: []string{"CLEURA_API_DEFAULT_DOMAIN_ID"},
			},
			&cli.StringFlag{
				Name:    "config-path",
				Aliases: []string{"p"},
				Usage:   "Path to configuration file. $HOME/.config/cleura/config if not set",
			},
		},
		Action: func(ctx *cli.Context) error {
			err := utils.ValidateNotEmpty(ctx,
				"token",
				"username",
				"api-host",
				"domain-id",
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
			projects, err := client.ListProjects(ctx.String("domain-id"))
			if err != nil {
				re, ok := err.(*cleura.RequestAPIError)
				if ok {
					if re.StatusCode == 403 {
						return fmt.Errorf("invalid token")
					}
				}
				return err
			}
			t := table.NewWriter()
			t.SetAutoIndex(true)
			t.Style().Format.Header = text.FormatTitle
			t.AppendHeader(table.Row{"Name", "Project Id", "Status", "Domain Id", "Description"})
			for _, project := range *projects {
				t.AppendRow(table.Row{project.Name, project.Id, fmt.Sprintf("default:%s\nenabled:%s", strconv.FormatBool(project.Default), strconv.FormatBool(project.Enabled)), project.DomainId, project.Description})
			}
			fmt.Println(t.Render())
			return nil
		},
	}
}
