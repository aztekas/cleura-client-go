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
	app.Version = "v0.0.3"
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
							Name:     "token",
							Aliases:  []string{"t"},
							Usage:    "Token to revoke",
							Required: true,
							EnvVars:  []string{"CLEURA_API_TOKEN"},
						},
						&cli.StringFlag{
							Name:     "username",
							Aliases:  []string{"u"},
							Usage:    "Username token belongs to",
							Required: true,
							EnvVars:  []string{"CLEURA_API_USERNAME"},
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
							Name:     "token",
							Required: true,
							Aliases:  []string{"t"},
							Usage:    "Token to validate",
							EnvVars:  []string{"CLEURA_API_TOKEN"},
						},
						&cli.StringFlag{
							Name:     "username",
							Required: true,
							Aliases:  []string{"u"},
							Usage:    "Username token belongs to",
							EnvVars:  []string{"CLEURA_API_USERNAME"},
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
			Name: "domain",
			Subcommands: []*cli.Command{
				{
					Name:   "list",
					Action: domainsList,
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:     "token",
							Aliases:  []string{"t"},
							Usage:    "Token to validate",
							Required: true,
							EnvVars:  []string{"CLEURA_API_TOKEN"},
						},
						&cli.StringFlag{
							Name:     "username",
							Required: true,
							Aliases:  []string{"u"},
							Usage:    "Username token belongs to",
							EnvVars:  []string{"CLEURA_API_USERNAME"},
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
			Name: "project",
			Subcommands: []*cli.Command{
				{
					Name:   "list",
					Action: projectsList,
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:     "token",
							Required: true,
							Aliases:  []string{"t"},
							Usage:    "Token to validate",
							EnvVars:  []string{"CLEURA_API_TOKEN"},
						},
						&cli.StringFlag{
							Name:     "username",
							Required: true,
							Aliases:  []string{"u"},
							Usage:    "Username token belongs to",
							EnvVars:  []string{"CLEURA_API_USERNAME"},
						},
						&cli.StringFlag{
							Name:    "api-host",
							Aliases: []string{"host"},
							Usage:   "Cleura API host",
							Value:   "https://rest.cleura.cloud",
						},
						&cli.StringFlag{
							Name:     "domain-id",
							Required: true,
							Aliases:  []string{"d"},
							Usage:    "Openstack domain id. Try \"cleura domain list\" for the list of available domains",
							EnvVars:  []string{"CLEURA_API_DEFAULT_DOMAIN_ID"},
						},
					},
				},
			},
		},
		{
			Name: "shoot",
			Subcommands: []*cli.Command{
				{
					Name:        "generate-kubeconfig",
					Action:      getKubeconfig,
					Description: "Overwrites existing file",
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
						},
						&cli.StringFlag{
							Name:     "cluster-name",
							Required: true,
							Aliases:  []string{"n"},
							Usage:    "Shoot cluster name",
						},
						&cli.StringFlag{
							Name:     "cluster-region",
							Required: true,
							Aliases:  []string{"r"},
							Usage:    "Openstack cluster region. Try \"cleura domains list\" command for available regions in your domain",
							EnvVars:  []string{"CLEURA_API_DEFAULT_REGION"},
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
	host := c.String("api-host")
	client, err := cleura.NewClientNoPassword(&host, &username, &token)
	if err != nil {
		return err
	}
	err = client.ValidateToken()
	if err != nil {
		re, ok := err.(*cleura.RequestAPIError)
		if ok {
			if re.StatusCode == 403 {
				return fmt.Errorf("error: token is invalid")
			}
		}
		return err
	}
	log.Println("token is valid")
	return nil
}
func tokenRevoke(c *cli.Context) error {
	token := c.String("token")
	username := c.String("username")
	host := c.String("api-host")
	client, err := cleura.NewClientNoPassword(&host, &username, &token)
	if err != nil {
		return err
	}
	err = client.RevokeToken()
	if err != nil {
		return err
	}
	log.Println("token successfully revoked")
	return nil
}
func tokenPrint(c *cli.Context) error {
	token := os.Getenv("CLEURA_API_TOKEN")
	if token == "" {
		return errors.New("error: token is not set")
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

func getKubeconfig(c *cli.Context) error {

	token := c.String("token")
	username := c.String("username")
	host := c.String("api-host")

	clusterName := c.String("cluster-name")
	clusterRegion := c.String("cluster-region")
	clusterProjectId := c.String("project-id")
	configDuration := c.Int64("config-duration")
	outputPath := c.String("output-path")

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
}
