package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/aztekas/cleura-client-go/pkg/api/cleura"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/urfave/cli/v2"
)

func main() {
	app := cli.NewApp()
	app.Name = "Cleura API CLI"
	//var commonOutput string
	app.Version = "v0.0.2"
	app.Commands = []*cli.Command{
		{
			Name: "token",
			Subcommands: []*cli.Command{
				{
					Name:   "get",
					Action: tokenGet,
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:    "username",
							Aliases: []string{"u"},
							Usage:   "Username for token request",
							EnvVars: []string{"CLEURA_API_USERNAME"},
						},
						&cli.StringFlag{
							Name:    "password",
							Aliases: []string{"p"},
							Usage:   "Password for token request",
							EnvVars: []string{"CLEURA_API_PASSWORD"},
						},
						&cli.StringFlag{
							Name:    "credentials-file",
							Aliases: []string{"c"},
							Usage:   "Path to credentials json file",
						},
						&cli.StringFlag{
							Name:    "api-host",
							Aliases: []string{"host"},
							Usage:   "Cleura API host",
							Value:   "https://rest.cleura.cloud",
						},
						//Add interactive mode
						//Add two factor mode
					},
				},
				{
					Name:   "revoke",
					Action: tokenRevoke,
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:    "token",
							Aliases: []string{"t"},
							Usage:   "Token to revoke",
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
						},
					},
				},
				{
					Name:   "print",
					Action: tokenPrint,
				},
				{
					Name:   "validate",
					Action: tokenValidate,
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
						},
					},
				},
			},
		},
		{
			Name: "domains",
			Subcommands: []*cli.Command{
				{
					Name:   "list",
					Action: domainsList,
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
						},
					},
				},
			},
		},
		{
			Name: "projects",
			Subcommands: []*cli.Command{
				{
					Name:   "list",
					Action: projectsList,
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
						},
						&cli.StringFlag{
							Name:    "domain-id",
							Aliases: []string{"d"},
							Usage:   "Openstack domain id. Try `cleura domains list` for the list of available domains",
						},
					},
				},
			},
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
func tokenValidate(c *cli.Context) error {
	token := c.String("token")
	username := c.String("username")
	if token == "" {
		return errors.New("error: token is not provided")
	}
	if username == "" {
		return errors.New("error: username is not provided")
	}
	host := c.String("api-host")
	client, err := cleura.NewClientNoPassword(&host, &username, &token)
	if err != nil {
		return err
	}
	err = client.ValidateToken()
	if err != nil {
		return err
	}
	fmt.Println("token is valid")
	return nil
}
func tokenRevoke(c *cli.Context) error {
	token := c.String("token")
	username := c.String("username")
	if token == "" {
		return errors.New("error: token is not provided")
	}
	if username == "" {
		return errors.New("error: username is not provided")
	}
	host := c.String("api-host")
	client, err := cleura.NewClientNoPassword(&host, &username, &token)
	if err != nil {
		return err
	}
	err = client.RevokeToken()
	if err != nil {
		return err
	}
	fmt.Println("token successfully revoked")
	return nil
}
func tokenPrint(c *cli.Context) error {
	token := os.Getenv("CLEURA_API_TOKEN")
	if token == "" {
		return errors.New("CLEURA_API_TOKEN is not set")
	}
	fmt.Println(token)
	return nil
}
func tokenGet(c *cli.Context) error {
	var host, username, password string
	if c.String("credentials-file") != "" {
		p := c.String("credentials-file")
		credentials := struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}{}
		file, err := os.Open(filepath.Join(p))
		if err != nil {
			return err
		}
		defer file.Close()
		jsonByte, err := io.ReadAll(file)
		if err != nil {
			return err
		}
		err = json.Unmarshal(jsonByte, &credentials)
		if err != nil {
			return err
		}
		username = credentials.Username
		password = credentials.Password

	} else {
		username = c.String("username")
		password = c.String("password")
	}
	//Add option to supply u&p via console
	if username == "" || password == "" {
		return errors.New("error: password and username must be supplied")
	}
	host = c.String("api-host")

	client, err := cleura.NewClient(&host, &username, &password)
	if err != nil {
		return err
	}
	fmt.Printf("export CLEURA_API_TOKEN=%v\nexport CLEURA_API_USERNAME=%v\nexport CLEURA_API_HOST=%v\n", client.Token, client.Auth.Username, c.String("api-host"))
	return nil
}

func domainsList(c *cli.Context) error {

	token := c.String("token")
	username := c.String("username")
	if token == "" {
		return errors.New("error: token is not provided")
	}
	if username == "" {
		return errors.New("error: username is not provided")
	}
	host := c.String("api-host")
	client, err := cleura.NewClientNoPassword(&host, &username, &token)
	if err != nil {
		return err
	}
	domains, err := client.ListDomains()
	if err != nil {
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

}

func projectsList(c *cli.Context) error {

	token := c.String("token")
	username := c.String("username")
	domain_id := c.String("domain-id")
	if token == "" {
		return errors.New("error: token is not provided")
	}
	if username == "" {
		return errors.New("error: username is not provided")
	}
	if domain_id == "" {
		return errors.New("error: domain-id is not provided")
	}
	host := c.String("api-host")
	client, err := cleura.NewClientNoPassword(&host, &username, &token)
	if err != nil {
		return err
	}
	projects, err := client.ListProjects(domain_id)
	if err != nil {
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
}
