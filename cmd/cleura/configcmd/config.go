package configcmd

import (
	"bytes"
	"errors"
	"fmt"
	"slices"

	"github.com/aztekas/cleura-client-go/cmd/cleura/utils"
	"github.com/spf13/viper"
	"github.com/urfave/cli/v2"
	"golang.org/x/exp/maps"
	"gopkg.in/yaml.v2"
)

type Config struct {
	ActiveConfig   string                   `mapstructure:"active_configuration" yaml:"active_configuration"`
	Configurations map[string]Configuration `mapstructure:"configurations"`
	location       string
}
type Configuration struct {
	Username         string `mapstructure:"username" yaml:"username"`
	Token            string `mapstructure:"token" yaml:"token"`
	DefaultDomainID  string `mapstructure:"domain-id" yaml:"domain-id"`
	DefaultRegion    string `mapstructure:"region" yaml:"region"`
	DefaultProjectID string `mapstructure:"project-id" yaml:"project-id"`
	//APIHost          string `mapstructure:"api-host" yaml:"api-host"`
}

func (c *Config) SetActive(name string) error {
	viper.SetConfigFile(c.location)
	viper.SetConfigType("yaml")

	_, ok := c.Configurations[name]
	if !ok {
		return errors.New("error: specified configuration name is not present in configuration file")
	}
	c.ActiveConfig = name
	configByte, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	err = viper.ReadConfig(bytes.NewBuffer(configByte))
	if err != nil {
		return err
	}
	err = viper.WriteConfig()
	if err != nil {
		return err
	}
	return nil
}

func (c *Config) UpdateConfiguration(key string, value string) error {

	_, ok := c.Configurations[c.ActiveConfig]
	if !ok {
		return errors.New("error: no active configuration in configuration file")
	}
	viper.SetConfigFile(c.location)
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}
	viper.Set("configurations."+c.ActiveConfig+"."+key, value)
	err = viper.WriteConfig()
	if err != nil {
		return err
	}
	return nil
}

func (c *Config) PrintConfigurations() {
	if len(c.Configurations) > 0 {
		fmt.Printf("Available configurations: (in %s)\n", c.location)
		for i, conf := range maps.Keys(c.Configurations) {
			if conf == c.ActiveConfig {
				conf = conf + " (active)"
			}
			fmt.Printf("%d. %s\n", i, conf)
		}
	} else {
		fmt.Printf("no configurations found in configuration file: %s\n", c.location)
	}
}

// Print specified configuration to stdout
func (c *Config) PrintConfigurationContent(name string) error {
	var confName string
	if name == "" {
		confName = c.ActiveConfig
	} else {
		confName = name
	}
	_, ok := c.Configurations[confName]
	if !ok {
		return fmt.Errorf("error: configuration name: %s not found in configuration file %s", confName, c.location)
	}
	configurationMap, err := utils.ToMap(c.Configurations[confName], "yaml")
	if err != nil {
		return err
	}
	fmt.Printf("Configuration name: `%s` (%s)\n\n", confName, c.location)
	for _, key := range maps.Keys(configurationMap) {
		if key == "token" && configurationMap[key] != "" {
			fmt.Printf("%-10s: ****hidden****\n", "token")
		} else {
			fmt.Printf("%-10s: %s\n", key, configurationMap[key])
		}

	}
	return nil
}

// Loads configuration settings from specified file path
func LoadConfiguration(path string) (*Config, error) {
	var config Config
	path, err := utils.ChoosePath(path)
	if err != nil {
		return nil, err
	}
	viper.SetConfigFile(path)
	viper.SetConfigType("yaml")
	err = viper.ReadInConfig() // Find and read the config file
	if err != nil {            // Handle errors reading the config file
		return nil, fmt.Errorf("LoadConfiguration: error reading config file: %w", err)
	}
	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, err
	}
	config.location = path
	return &config, nil
}

// Get available configuration settings from supplied or default configuration file
func TrySetConfigFromFile(c *cli.Context) error {
	var config *Config
	//ignore api-host as it always has default value
	//ignore path as it always has default value (%HOME/.config/cleura/config)
	var ignoreFlags = []string{"help"}
	config, err := LoadConfiguration(c.String("config-path"))
	if err != nil {
		fmt.Printf("Failed to read configuration file: %s, proceed with explicit configuration via flags and environment variables\n", err)
		return nil
	}
	activeConfig := config.ActiveConfig
	if activeConfig == "" {
		fmt.Printf("No active configuration found in %s, proceed with explicit configuration via flags and environment variables\n", config.location)
		return nil
	}
	// Get a map put of struct fields with yaml tag
	configFromFile, err := utils.ToMap(config.Configurations[activeConfig], "yaml")
	if err != nil {
		return err
	}

	//Iterate over all defined flags (set and unset)
	for _, flag := range c.Command.Flags {
		ok := c.IsSet(flag.Names()[0])
		// Only set flags that are not already set directly of via environment variables
		_, inConfig := configFromFile[flag.Names()[0]]
		if !ok && !slices.Contains(ignoreFlags, flag.Names()[0]) && inConfig {
			c.Set(flag.Names()[0], configFromFile[flag.Names()[0]].(string))
		}
	}
	//return fmt.Errorf("error: expected error")
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
			showActiveCommand(),
			generateConfigTemplateCommand(),
			showCommand(),
		},
	}
}
