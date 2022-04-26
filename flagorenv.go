package flagorenv

import (
	"flag"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/mcuadros/go-defaults"
	"github.com/stoewer/go-strcase"
)

type Config struct {
	Prefix     string
	PreferFlag bool
}

const (
	prefix = "FLAGENV"
)

func LoadFlagsOrEnv[T any](c *Config) (T, error) {
	var res T
	if c.Prefix == "" {
		c.Prefix = prefix
	}

	envVars, err := loadEnv[T](c)
	if err != nil {
		return res, err
	}
	flagVars, err := loadFlags[T](c)
	if err != nil {
		return res, err
	}

	// merge envVars and flagVars, giving precedence to envVars
	res, err = merge(c, envVars, flagVars)

	return res, err
}

func merge[T any](c *Config, envVars T, flagVars T) (T, error) {
	var res T

	// set defaults according to the 'default:' tag
	defaults.SetDefaults(&res)

	// for every pair of fields in T,
	// compare if either is a default value
	// if so, set the res to default
	// otherwise, set the res to the value from the flag
	value := reflect.ValueOf(&res).Elem()
	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		envField := reflect.ValueOf(envVars).Field(i)
		flagField := reflect.ValueOf(flagVars).Field(i)

		// if both the envField and flagField have a value that is not Zero,
		// set the res to the value according to the config.PreferFlag option
		// If only one of them has a value that is not Zero,
		// set the field to that value

		if !envField.IsZero() && !flagField.IsZero() {
			if c.PreferFlag {
				field.Set(flagField)
			} else {
				field.Set(envField)
			}
		}

		if !envField.IsZero() && flagField.IsZero() {
			field.Set(envField)
		}

		if envField.IsZero() && !flagField.IsZero() {
			field.Set(flagField)
		}
	}

	return res, nil
}

func loadFlags[T any](c *Config) (T, error) {
	var res T

	// use reflection to iterate through every field of T
	// and load the value from the flag with the same name as the field
	// String Slices are separated by ","
	// Ints and Floats are parsed as is
	// Bools are parsed as "true" and "false"

	var tempSliceStr string

	resref := reflect.ValueOf(&res).Elem()
	for i := 0; i < resref.NumField(); i++ {
		field := resref.Type().Field(i)
		fieldName := strcase.KebabCase(c.Prefix + field.Name)
		fieldInterface := resref.Field(i).Addr().Interface()

		// depending on the field type, set the flag module switch
		switch resref.Field(i).Kind() {
		case reflect.String:
			flag.StringVar(fieldInterface.(*string), fieldName, "", "")
		case reflect.Slice:
			flag.StringVar(&tempSliceStr, fieldName, "", "")
		case reflect.Int64:
			flag.Int64Var(fieldInterface.(*int64), fieldName, 0, "")
		case reflect.Float64:
			flag.Float64Var(fieldInterface.(*float64), fieldName, 0, "")
		case reflect.Bool:
			flag.BoolVar(fieldInterface.(*bool), fieldName, false, "")
		default:
			return res, fmt.Errorf("unsupported type %s", resref.Field(i).Kind())
		}
	}

	// consume flag
	flag.Parse()

	// split the string by ","
	// and set the field to the resulting slice
	for i := 0; i < resref.NumField(); i++ {

		// depending on the field type, set the flag module switch
		switch resref.Field(i).Kind() {
		case reflect.Slice:
			splitStrs := reflect.ValueOf(strings.Split(tempSliceStr, ","))
			resref.Field(i).Set(splitStrs)
		}
	}

	return res, nil
}

func loadEnv[T any](c *Config) (T, error) {
	var res T

	// use reflection to iterate through every field of T
	// and load the value from the environment variable
	// with the same name as the field
	// String Slices are separated by ","
	// Ints and Floats are parsed as is
	// Bools are parsed as "true" and "false"

	value := reflect.ValueOf(&res).Elem()
	for i := 0; i < value.NumField(); i++ {
		fieldName := value.Type().Field(i).Name
		// if the fieldname has CamelCase, convert it to snake_case
		fieldName = strcase.UpperSnakeCase(c.Prefix + fieldName)

		// get the environment variable
		envVar := os.Getenv(fieldName)

		// parse the envVar to whatever the type of the current field is
		switch value.Field(i).Kind() {
		case reflect.String:
			value.Field(i).SetString(envVar)
		case reflect.Slice:
			// split the string by ","
			// and set the field to the resulting slice
			value.Field(i).Set(reflect.ValueOf(strings.Split(envVar, ",")))
		case reflect.Int64:
			// parse the string as an int
			value.Field(i).SetInt(int64(parseInt(envVar)))
		case reflect.Float64:
			// parse the string as a float
			value.Field(i).SetFloat(parseFloat(envVar))
		case reflect.Bool:
			// parse the string as a bool
			value.Field(i).SetBool(parseBool(envVar))
		default:
			return res, fmt.Errorf("unsupported type %s", value.Field(i).Kind())
		}
	}

	return res, nil
}

func parseInt(s string) int {
	if s == "" {
		return 0
	}

	i, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return i
}

func parseFloat(s string) float64 {
	if s == "" {
		return 0.0
	}

	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		panic(err)
	}
	return f
}

func parseBool(s string) bool {
	return s == "true"
}
