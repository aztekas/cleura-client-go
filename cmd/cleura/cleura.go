package main

import (
	"log"
	"os"

	"github.com/aztekas/cleura-client-go/cmd/cleura/configcmd"
	"github.com/aztekas/cleura-client-go/cmd/cleura/domaincmd"
	"github.com/aztekas/cleura-client-go/cmd/cleura/projectcmd"
	"github.com/aztekas/cleura-client-go/cmd/cleura/shootcmd"
	"github.com/aztekas/cleura-client-go/cmd/cleura/tokencmd"
	"github.com/urfave/cli/v2"
)

var (
	version = "dev"
)

func commands() []*cli.Command {
	return append(
		[]*cli.Command{},
		configcmd.Command(),
		domaincmd.Command(),
		projectcmd.Command(),
		tokencmd.Command(),
		shootcmd.Command(),
	)
}

func main() {
	app := cli.NewApp()
	app.Name = "Cleura API CLI"
	app.Description = "Cleura API CLI application to work with Cleura API"
	app.Version = version
	app.Commands = commands()
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
