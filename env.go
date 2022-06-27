package env

import (
	"encoding"
	"fmt"
	"os"
	"reflect"
	"strconv"
)

const tagKey = "env"

func FillConfigFromEnv(config any, appPrefix string) error {
	v := reflect.ValueOf(config).Elem()

	return processValue(v, appPrefix)
}

func processValue(v reflect.Value, appPrefix string) error {
	for i := 0; i < v.NumField(); i++ {
		if !v.Type().Field(i).IsExported() {
			continue
		}

		tagValue := v.Type().Field(i).Tag.Get(tagKey)
		envKey := appPrefix + "_" + tagValue
		envVal, envValSet := os.LookupEnv(envKey)

		fieldKind := v.Field(i).Kind()
		fieldType := v.Field(i).Type()

		if fieldKind == reflect.Struct {
			if tagValue == "" {
				if err := processValue(v.Field(i), appPrefix); err != nil {
					return err
				}
				continue
			}
			if !envValSet {
				continue // Nothing to be overriden or set
			}
			um, ok := v.Field(i).Interface().(encoding.TextUnmarshaler)
			if !ok {
				return fmt.Errorf("type does not implement TextUnmarshaler interface: %s", fieldType)
			}
			if err := um.UnmarshalText([]byte(envVal)); err != nil {
				return fmt.Errorf("filling field with env tagValue: %s, %w", tagValue, err)
			}
			v.Field(i).Set(reflect.ValueOf(um).Elem())
			continue
		}

		if tagValue != "" && envValSet {
			if err := setPrimitive(envVal, v.Field(i), fieldKind); err != nil {
				return err
			}
		}

	}

	return nil
}

func setPrimitive(str string, v reflect.Value, k reflect.Kind) error {
	fieldType := v.Type()
	switch k {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i, err := strconv.ParseInt(str, 10, int(fieldType.Size()))
		if err != nil {
			return err
		}
		v.SetInt(i)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		u, err := strconv.ParseUint(str, 10, int(fieldType.Size()))
		if err != nil {
			return err
		}
		v.SetUint(u)
	case reflect.Float32, reflect.Float64:
		f, err := strconv.ParseFloat(str, int(fieldType.Size()))
		if err != nil {
			return err
		}
		v.SetFloat(f)
	case reflect.Bool:
		b, err := strconv.ParseBool(str)
		if err != nil {
			return err
		}
		v.SetBool(b)
	case reflect.String:
		v.SetString(str)
	default:
		return fmt.Errorf("unsupported kind: %s", k)
	}
	return nil

}
