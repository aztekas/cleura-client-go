package tokencmd

import (
	"fmt"
	"log"

	"github.com/aztekas/cleura-client-go/cmd/cleura/configcmd"
	"github.com/aztekas/cleura-client-go/cmd/cleura/utils"
	"github.com/aztekas/cleura-client-go/pkg/api/cleura"
	"github.com/urfave/cli/v2"
)

func revokeCommand() *cli.Command {
	return &cli.Command{
		Name:        "revoke",
		Description: "Revoke currently active token",
		Usage:       "Revoke currently active token",
		Before:      configcmd.TrySetConfigFromFile,
		Flags:       utils.CommonFlags(),
		Action: func(ctx *cli.Context) error {
			token := ctx.String("token")
			username := ctx.String("username")
			host := ctx.String("api-host")
			client, err := cleura.NewClientNoPassword(&host, &username, &token)
			if err != nil {
				return err
			}
			err = client.RevokeToken()
			if err != nil {
				re, ok := err.(*cleura.RequestAPIError)
				if ok {
					if re.StatusCode == 403 {
						return fmt.Errorf("error: invalid token")
					}
				}
				return err
			}
			log.Println("token successfully revoked")
			return nil

		},
	}

}
