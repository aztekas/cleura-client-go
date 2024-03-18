package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
)

// Return default location within user home directory ($HOME/.config/cleura/config)
func GetDefaultConfigPath() (string, error) {
	homedir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homedir, ".config", "cleura", "config"), nil
}

// Choose path for the configuration file. Choose path` if supplied
// otherwise set path within default user directory
func ChoosePath(path string) (string, error) {
	var err error
	if path == "" {
		path, err = GetDefaultConfigPath()
		if err != nil {
			return "", fmt.Errorf("error: setting default path failed: %w", err)
		}
	}
	return path, nil
}

// Return a map of KV from an arbitrary struct
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
