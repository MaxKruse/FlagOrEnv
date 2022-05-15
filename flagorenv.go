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
		case reflect.Int, reflect.Int32:
			flag.IntVar(fieldInterface.(*int), fieldName, 0, "")
		case reflect.Float32, reflect.Float64:
			flag.Float64Var(fieldInterface.(*float64), fieldName, 0, "")
		case reflect.Bool:
			flag.BoolVar(fieldInterface.(*bool), fieldName, false, "")
		default:
			return res, fmt.Errorf("unsupported type %s", resref.Field(i).Kind())
		}
	}

	// consume flag
	flag.Parse()

	// iterate through every field of T. if a slice is encountered,
	// call the parseSlice function to parse the slice accordingly

	for i := 0; i < resref.NumField(); i++ {
		field := resref.Type().Field(i)

		if field.Type.Kind() != reflect.Slice {
			continue
		}

		// parse the slice
		slice, err := parseSlice(tempSliceStr, field.Type.Elem())
		if err != nil {
			return res, err
		}

		// set the slice to the field
		resref.Field(i).Set(slice)
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
		fieldName = strcase.UpperSnakeCase(c.Prefix + "_" + fieldName)

		// get the environment variable
		envVar := os.Getenv(fieldName)

		// parse the envVar to whatever the type of the current field is
		switch value.Field(i).Kind() {
		case reflect.String:
			value.Field(i).SetString(envVar)
		case reflect.Slice:
			slice, err := parseSlice(envVar, value.Field(i).Type().Elem())
			if err != nil {
				return res, err
			}
			value.Field(i).Set(slice)
		case reflect.Int64:
			// parse the string as an int
			value.Field(i).SetInt(parseInt64(envVar))
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

func parseSlice(str string, elem reflect.Type) (reflect.Value, error) {
	// split str by ","
	// iterate through the current field
	// parse the string as the current field type
	// append the parsed value to the slice
	// return the slice

	if str == "" {
		return reflect.MakeSlice(reflect.SliceOf(elem), 0, 0), nil
	}

	slice := reflect.MakeSlice(reflect.SliceOf(elem), 0, 0)
	for _, s := range strings.Split(str, ",") {

		// additionally, strip the string of whitespace
		s = strings.TrimSpace(s)

		// parse the string as the current field type
		// append the parsed value to the slice
		switch elem.Kind() {
		case reflect.String:
			slice = reflect.Append(slice, reflect.ValueOf(s))
		case reflect.Int64:
			slice = reflect.Append(slice, reflect.ValueOf(parseInt64(s)))
		case reflect.Int32, reflect.Int:
			slice = reflect.Append(slice, reflect.ValueOf(parseInt32(s)))
		case reflect.Float64, reflect.Float32:
			slice = reflect.Append(slice, reflect.ValueOf(parseFloat(s)))
		case reflect.Bool:
			slice = reflect.Append(slice, reflect.ValueOf(parseBool(s)))
		default:
			return slice, fmt.Errorf("unsupported type %s", elem.Kind())
		}
	}

	return slice, nil
}

func parseInt32(str string) int32 {
	i, err := strconv.ParseInt(str, 10, 32)
	if err != nil {
		panic(err)
	}

	return int32(i)
}

func parseInt64(s string) int64 {
	if s == "" {
		return 0
	}

	i, err := strconv.ParseInt(s, 10, 64)
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
