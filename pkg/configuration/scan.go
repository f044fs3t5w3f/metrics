package configuration

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"time"
)

var ErrNotSupportedType = errors.New("not supported type")

func ScanConfig(ptr any, args []string) error {
	if args == nil {
		args = os.Args[1:]
	}
	v := reflect.ValueOf(ptr)

	fs := flag.NewFlagSet("test", flag.ContinueOnError)

	configPath := os.Getenv("CONFIG")
	if configPath == "" {
		fs.StringVar(&configPath, "c", "", "path to config file")
		fs.StringVar(&configPath, "config", "", "path to config file")
	}

	if v.Kind() != reflect.Ptr || v.IsNil() {
		return fmt.Errorf("nil")
	}

	v = v.Elem()

	if v.Kind() != reflect.Struct {
		return fmt.Errorf("expected pointer to struct")
	}

	t := v.Type()

	flags := make(map[string]*string)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		flagName := field.Tag.Get("flag")

		description := field.Tag.Get("description")
		if flagName != "" {
			flags[flagName] = fs.String(flagName, "", description)
		}
	}
	if err := fs.Parse(args); err != nil {
		return err
	}
	jsonConfig, err := loadJSONConfig(configPath)
	if err != nil {
		return fmt.Errorf("loadJSONConfig: %w", err)
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		if !field.IsExported() {
			continue
		}

		envName := field.Tag.Get("env")
		jsonName := field.Tag.Get("jsonConfig")
		flagName := field.Tag.Get("flag")
		defaultValue := field.Tag.Get("default")
		value := selectValue(envName, flags[flagName], jsonConfig[jsonName], defaultValue)
		switch field.Type.Kind() {
		case reflect.String:
			v.FieldByIndex(field.Index).Set(reflect.ValueOf(value))
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if field.Type == reflect.TypeOf(time.Duration(0)) {
				parsedValue := time.Duration(0)
				if value != "" {
					parsedValue, err = time.ParseDuration(value)
					if err != nil {
						return fmt.Errorf("cannot set field %s. Value %s. %w ", field.Name, value, err)
					}
				}

				v.FieldByIndex(field.Index).Set(reflect.ValueOf(parsedValue))
			} else {
				parsedValue, err := strconv.ParseInt(value, 10, 64)
				if err != nil {
					return fmt.Errorf("cannot set field %s. Value %s. %w ", field.Name, value, err)
				}
				fieldType := field.Type
				v.FieldByIndex(field.Index).Set(reflect.ValueOf(parsedValue).Convert(fieldType))
			}
		case reflect.Bool:
			parsedValue := false
			if value != "" {
				parsedValue, err = strconv.ParseBool(value)
			}
			if err != nil {
				return fmt.Errorf("cannot set field %s. Value %s. %w ", field.Name, value, err)
			}
			v.FieldByIndex(field.Index).Set(reflect.ValueOf(parsedValue))
		default:
			return fmt.Errorf("cannot set field %s. %w", field.Name, ErrNotSupportedType)
		}
	}
	return nil
}

func selectValue(envName string, flagValuePtr *string, jsonConfigPtr *string, defaultValue string) string {
	if v := os.Getenv(envName); v != "" {
		return v
	}
	if flagValuePtr != nil && *flagValuePtr != "" {
		return *flagValuePtr
	}
	if jsonConfigPtr != nil && *jsonConfigPtr != "" {
		return *jsonConfigPtr
	}
	return defaultValue
}
