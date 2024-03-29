package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"

	"github.com/urfave/cli/v2"
	"golang.org/x/term"
)

// Return default location within user home directory ($HOME/.config/cleura/config).
func GetDefaultConfigPath() string {
	homedir, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(homedir, ".config", "cleura", "config")
}

// Choose path for the configuration file. Choose path` if supplied
// otherwise set path within default user directory.
func ChoosePath(path string) (string, error) {
	if path == "" {
		path = GetDefaultConfigPath()
		if path == "" {
			return "", fmt.Errorf("error: setting default path failed")
		}
	}
	return path, nil
}

// Return a map of KV from an arbitrary struct.
func ToMap(in interface{}, tag string) (map[string]interface{}, error) {
	out := make(map[string]interface{})

	v := reflect.ValueOf(in)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// we only accept structs
	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("ToMap only accepts structs; got %T", v)
	}

	typ := v.Type()
	for i := 0; i < v.NumField(); i++ {
		// gets us a StructField
		fi := typ.Field(i)
		if tagv := fi.Tag.Get(tag); tagv != "" {
			// set key of map to value in struct field
			out[tagv] = v.Field(i).Interface()
		}
	}
	return out, nil
}

func ValidateNotEmptyString(ctx *cli.Context, flags ...string) error {
	for _, flag := range flags {
		if ctx.String(flag) == "" {
			return fmt.Errorf("error: required flag: `%s` is not set or empty", flag)
		}
	}
	return nil
}

func GetUserInput(asking string, masked bool) (userInput string, err error) {
	fmt.Printf("Enter %s: ", asking)
	if masked {

		secretInput, err := term.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			return "", err
		}
		fmt.Println("")
		return string(secretInput), nil
	}
	var input string
	fmt.Scanln(&input)
	return input, nil
}
