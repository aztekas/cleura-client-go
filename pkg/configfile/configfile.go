package configfile

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/aztekas/cleura-client-go/pkg/util"
	"github.com/gofrs/flock"
	"golang.org/x/exp/maps"
	"gopkg.in/yaml.v2"
)

type Configuration struct {
	configFile configFile
	Location   string
}

type configFile struct {
	ActiveProfile string             `yaml:"active_profile"`
	Profiles      map[string]profile `yaml:"profiles"`
}
type profile struct {
	Username         string `yaml:"username,omitempty"`
	Token            string `yaml:"token,omitempty"`
	DefaultDomainID  string `yaml:"domain-id,omitempty"`
	DefaultRegion    string `yaml:"region,omitempty"`
	DefaultProjectID string `yaml:"project-id,omitempty"`
	ApiUrl           string `yaml:"api-url,omitempty"`
}

// Validate configuration file for active profile and profile data.
func (cf *configFile) validateConfigurationFile() error {
	if cf.ActiveProfile == "" {
		return fmt.Errorf("active profile is not set")
	}
	if len(cf.Profiles) == 0 {
		return fmt.Errorf("no profile data found")
	}
	if _, ok := cf.Profiles[cf.ActiveProfile]; !ok {
		return fmt.Errorf("no profile data exist for chosen active profile `%s`", cf.ActiveProfile)
	}
	return nil
}

// Return current active profile name.
func (c *Configuration) GetActiveProfile() string {
	return c.configFile.ActiveProfile
}

// Add new profile to the configuration file.
func (c *Configuration) AddProfile(name string, username string, token string, domainid string, region string, projectid string, apiurl string) error {
	c.configFile.Profiles[name] = profile{
		Username:         username,
		Token:            token,
		DefaultDomainID:  domainid,
		DefaultRegion:    region,
		DefaultProjectID: projectid,
		ApiUrl:           apiurl,
	}
	err := writeConfigFile(c.Location, c.configFile)
	if err != nil {
		return err
	}
	return nil
}

// Set active profile from existing profiles in the file.
func (c *Configuration) SetActiveProfile(profile string) error {
	if profile == "" {
		return errors.New("error: profile name is empty")
	}
	_, ok := c.configFile.Profiles[profile]
	if !ok {
		return errors.New("error: specified profile name is not present in configuration file")
	}
	c.configFile.ActiveProfile = profile
	err := writeConfigFile(c.Location, c.configFile)
	if err != nil {
		return err
	}
	return nil
}

// Set arbitrary (among available profile fields) profile field value.
func (c *Configuration) SetProfileField(key string, value string) error {
	tagFieldMap, err := util.ToTagFieldNameMap(&profile{}, "yaml", 0)
	if err != nil {
		return err
	}
	if _, ok := tagFieldMap[key]; !ok {
		return fmt.Errorf("key: %s is not supported", key)
	}
	activeProfileData := c.configFile.Profiles[c.configFile.ActiveProfile]
	is_set, err := util.SetStructField(&activeProfileData, tagFieldMap[key], value)
	if err != nil {
		return err
	}
	if is_set {
		c.configFile.Profiles[c.configFile.ActiveProfile] = activeProfileData
		err = writeConfigFile(c.Location, c.configFile)
	}
	return err
}

// Return profile map[FieldTag]FieldValue.
func (c *Configuration) GetProfileMap(profile string) (map[string]interface{}, error) {
	var err error
	var activeProfileMap map[string]interface{}
	if profile == "" || profile == c.configFile.ActiveProfile {
		activeProfileMap, err = util.ToMap(c.configFile.Profiles[c.configFile.ActiveProfile], "yaml")
		if err != nil {
			return nil, err
		}
	} else {
		activeProfileMap, err = util.ToMap(c.configFile.Profiles[profile], "yaml")
		if err != nil {
			return nil, err
		}
	}
	return activeProfileMap, nil
}

// Return profile names slice.
func (c *Configuration) ProfilesSlice() []string {
	return maps.Keys(c.configFile.Profiles)
}

// Load configuration profiles from a yaml file on specified path.
func InitConfiguration(path string) (*Configuration, error) {
	var cFile configFile
	filename, err := setConfigPath(path)
	if err != nil {
		return nil, fmt.Errorf("provided path is empty, setting default path failed with: %w", err)
	}
	if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
		return nil, err
	} else {
		// file exists
		cFile, err = readConfigFile(filename)
		if err != nil {
			return nil, fmt.Errorf("reading file: %s failed with: %w", filename, err)
		}
	}
	err = cFile.validateConfigurationFile()
	if err != nil {
		return nil, fmt.Errorf("configuration file: `%s` is not valid, %w", filename, err)
	}
	return &Configuration{
		configFile: cFile,
		Location:   filename,
	}, nil
}

// Create configuration template file on specified path.
func CreateConfigTemplateFile(filename string) error {
	err := os.MkdirAll(filepath.Dir(filename), os.ModePerm)
	if err != nil {
		return fmt.Errorf("can not create path to the specified file, %w", err)
	}
	template := &configFile{
		ActiveProfile: "default",
		Profiles: map[string]profile{
			"default": {},
		},
	}
	err = writeConfigFile(filename, *template)
	if err != nil {
		return fmt.Errorf("creating configuration file failed: %w", err)
	}
	return nil
}

// Choose path for the configuration file. Choose `path` if supplied
// otherwise set path within default user directory.
func setConfigPath(path string) (string, error) {
	if path == "" {
		homedir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("error: setting default path failed")
		}
		return filepath.Join(homedir, ".config", "cleura", "config"), nil

	}
	return path, nil
}

// Read configuration file from disk.
func readConfigFile(filename string) (configFile, error) {
	lock := flock.New(filename)
	// wait up to a second for the file to lock
	config := configFile{
		Profiles: make(map[string]profile),
	}
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second)
	defer cancel()
	ok, err := lock.TryRLockContext(ctx, 250*time.Millisecond) // try to lock every 1/4 second
	if !ok {
		// unable to lock the cache, something is wrong, refuse to use it.
		return config, fmt.Errorf("unable to read lock file %s: %v", filename, err)
	}
	data, err := os.ReadFile(filename)
	if err != nil {
		return config, fmt.Errorf("failed to read cache file: %w", err)
	}
	if err := yaml.Unmarshal(data, &config); err != nil {
		return config, fmt.Errorf("unable to parse file %s: %w", filename, err)
	}
	err = lock.Unlock()
	if err != nil {
		return config, fmt.Errorf("unable to unlock file %s: %w", filename, err)
	}
	return config, nil
}

// Write configuration file to disk.
func writeConfigFile(filename string, file interface{}) error {
	lock := flock.New(filename)
	// wait up to a second for the file to lock
	ctx, cancel := context.WithTimeout(context.TODO(), 2*time.Second)
	defer cancel()
	ok, err := lock.TryRLockContext(ctx, 250*time.Millisecond) // try to lock every 1/4 second
	if !ok {
		// unable to lock the config file, something is wrong, refuse to use it.
		return fmt.Errorf("unable to read lock file %s: %v", filename, err)
	}
	data, err := yaml.Marshal(file)
	if err != nil {
		return err
	}
	err = os.WriteFile(filename, data, 0600)
	if err != nil {
		return err
	}
	err = lock.Unlock()
	if err != nil {
		return fmt.Errorf("unable to unlock file %s: %w", filename, err)
	}
	return nil
}

func WriteByteToFile(filename string, fileData []byte) error {
	lock := flock.New(filename)
	// wait up to a second for the file to lock
	ctx, cancel := context.WithTimeout(context.TODO(), 2*time.Second)
	defer cancel()
	ok, err := lock.TryRLockContext(ctx, 250*time.Millisecond) // try to lock every 1/4 second
	if !ok {
		// unable to lock the config file, something is wrong, refuse to use it.
		return fmt.Errorf("unable to read lock file %s: %v", filename, err)
	}
	err = os.WriteFile(filename, fileData, 0600)
	if err != nil {
		return err
	}
	err = lock.Unlock()
	if err != nil {
		return fmt.Errorf("unable to unlock file %s: %w", filename, err)
	}
	return nil
}
