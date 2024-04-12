package tokencmd

import (
	"fmt"

	"github.com/aztekas/cleura-client-go/cmd/cleura/common"
	"github.com/aztekas/cleura-client-go/cmd/cleura/configcmd"
	"github.com/aztekas/cleura-client-go/pkg/api/cleura"
	"github.com/urfave/cli/v2"
)

func validateCommand() *cli.Command {
	return &cli.Command{
		Name:        "validate",
		Description: "Validate currently active token",
		Usage:       "Validate currently active token",
		Before:      configcmd.TrySetConfigFromFile,
		Flags:       common.CleuraAuthFlags(),
		Action: func(ctx *cli.Context) error {
			logger := common.CliLogger(ctx.String("loglevel"))
			token := ctx.String("token")
			username := ctx.String("username")
			host := ctx.String("api-host")
			client, err := cleura.NewClientNoPassword(&host, &username, &token)
			if err != nil {
				return err
			}
			err = client.ValidateToken()
			if err != nil {
				re, ok := err.(*cleura.RequestAPIError)
				if ok {
					if re.StatusCode == 403 {
						return fmt.Errorf("error: token is invalid or not supplied")
					}
				}
				return err
			}
			logger.Info("token is valid")
			return nil
		},
	}

}
