package util

import (
	"fmt"
	"os"
	"reflect"
	"strings"

	"golang.org/x/term"
)

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
		if tagv := strings.Split(fi.Tag.Get(tag), ",")[0]; tagv != "" {
			// set key of map to value in struct field
			out[tagv] = v.Field(i).Interface()
		}
	}
	return out, nil
}

// Return map[tag]FieldName of a given struct.
// TODO: Add check if tag is not set.
func ToTagFieldNameMap(obj any, tag string, tagNum int) (map[string]string, error) {
	output := make(map[string]string)
	ref := reflect.ValueOf(obj)
	// if its a pointer, resolve its value
	if ref.Kind() == reflect.Ptr {
		ref = reflect.Indirect(ref)
	}
	if ref.Kind() == reflect.Interface {
		ref = ref.Elem()
	}
	// should double check we now have a struct (could still be anything)
	if ref.Kind() != reflect.Struct {
		return nil, fmt.Errorf("only structs are supported, got %T", ref)
	}
	typ := ref.Type()
	for i := 0; i < ref.NumField(); i++ {
		nameTag := strings.Split(typ.Field(i).Tag.Get(tag), ",")[tagNum]
		output[nameTag] = typ.Field(i).Name
	}
	return output, nil
}

// Set arbitrary struct field.
func SetStructField(obj any, key string, value string) (bool, error) {
	ref := reflect.ValueOf(obj)
	// if its a pointer, resolve its value
	if ref.Kind() == reflect.Ptr {
		ref = reflect.Indirect(ref)
	}
	if ref.Kind() == reflect.Interface {
		ref = ref.Elem()
	}
	// should double check we now have a struct (could still be anything)
	if ref.Kind() != reflect.Struct {
		return false, fmt.Errorf("only structs are supported, got %T", ref)
	}
	currentValue := ref.FieldByName(key)
	if value == currentValue.String() {
		return false, nil
	}
	currentValue.SetString(value)
	return true, nil
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
