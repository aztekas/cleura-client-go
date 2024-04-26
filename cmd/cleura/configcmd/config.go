package configcmd

import (
	"fmt"
	"slices"

	"github.com/aztekas/cleura-client-go/cmd/cleura/common"
	"github.com/aztekas/cleura-client-go/internal/configfile"
	"github.com/urfave/cli/v2"
)

// Get available configuration settings from supplied or default configuration file.
func TrySetConfigFromFile(c *cli.Context) error {
	logger := common.CliLogger(c.String("loglevel"))
	var config *configfile.Configuration
	// ignore api-host as it always has a default value
	// ignore path as it always has a default value (%HOME/.config/cleura/config)
	var ignoreFlags = []string{"help"}
	config, err := configfile.InitConfiguration(c.String("config-path"))
	if err != nil {
		logger.Warn(fmt.Sprintf("Failed to read configuration file: %s, proceed with explicit configuration via flags and environment variables", err))
		return nil
	}
	// Get a map out of struct fields
	profileMap, err := config.GetProfileMap(config.GetActiveProfile())
	if err != nil {
		return err
	}

	// Iterate over all defined flags (set and unset)
	for _, flag := range c.Command.Flags {
		ok := c.IsSet(flag.Names()[0])
		// Only set flags that are not already set directly or via environment variables
		_, inConfig := profileMap[flag.Names()[0]]
		if !ok && !slices.Contains(ignoreFlags, flag.Names()[0]) && inConfig {
			flagValueFromConfig, ok := profileMap[flag.Names()[0]].(string)
			if !ok {
				return fmt.Errorf("error: cannot assert string type assertion")
			}
			err = c.Set(flag.Names()[0], flagValueFromConfig)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func Command() *cli.Command {
	return &cli.Command{
		Name:        "config",
		Description: "Command used for working with configuration file for the cleura cli",
		Usage:       "Command used for working with configuration file for the cleura cli",
		Subcommands: []*cli.Command{
			setCommand(),
			listCommand(),
			generateConfigTemplateCommand(),
			showCommand(),
		},
	}
}
