package main

import (
	"errors"
	"log"
	"os"
	"slices"

	"github.com/aztekas/cleura-client-go/cmd/cleura/configcmd"
	"github.com/aztekas/cleura-client-go/cmd/cleura/domaincmd"
	"github.com/aztekas/cleura-client-go/cmd/cleura/projectcmd"
	"github.com/aztekas/cleura-client-go/cmd/cleura/shootcmd"
	"github.com/aztekas/cleura-client-go/cmd/cleura/tokencmd"
	"github.com/urfave/cli/v2"
)

var (
	version = "latest"
	commit  = "uncommitted"
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
	app.Usage = "A Cleura API CLI"
	app.Name = "cleura"
	app.Version = version + "-" + commit
	app.Commands = commands()
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:  "loglevel",
			Value: "info",
			Action: func(ctx *cli.Context, s string) error {
				pV := []string{"info", "warn", "error", "debug"}
				if !slices.Contains(pV, ctx.String("loglevel")) {
					return errors.New("--loglevel must be set to one of `info, warn, error, debug`")
				}
				return nil
			},
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
