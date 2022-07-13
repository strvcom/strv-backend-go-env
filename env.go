package env

import (
	"encoding"
	"fmt"
	"os"
	"reflect"
	"strconv"
)

const (
	decadicBase = 10
	envTag      = "env"

	// envVarAppPrefix is the environment variable that holds the prefix for the
	// environment variables specified by the `env` tag.
	envVarAppPrefix = "APP_PREFIX"
)

// Apply applies the environment variables to the given target.
func Apply(target any) error {
	return ApplyWithPrefix(target, os.Getenv(envVarAppPrefix))
}

// MustApply calls Apply and panics on error.
func MustApply(target any) {
	if err := Apply(target); err != nil {
		panic(err)
	}
}

// ApplyWithPrefix applies the environment variables to the given target.
//
// The prefix is used to prefix the environment variables specified by the `env` tag.
func ApplyWithPrefix(target any, prefix string) error {
	rv := reflect.ValueOf(target)

	switch rv.Kind() { // nolint:exhaustive
	case reflect.Pointer:
		typ := rv.Elem().Kind()
		if typ != reflect.Struct || rv.IsNil() {
			return &InvalidTypeError{rv}
		}
	default:
		return &InvalidTypeError{rv}
	}

	return applyWithPrefix(rv.Elem(), prefix)
}

// MustApplyWithPrefix calls ApplyWithPrefix and panics on error.
func MustApplyWithPrefix(target any, prefix string) {
	if err := ApplyWithPrefix(target, prefix); err != nil {
		panic(err)
	}
}

func applyWithPrefix(rv reflect.Value, prefix string) error {
	rt := rv.Type()
L:
	for i := 0; i < rv.NumField(); i++ {
		rf := rv.Field(i)

		if !rt.Field(i).IsExported() {
			continue L
		}

		tagVal := rt.Field(i).Tag.Get(envTag)
		switch tagVal {
		case "":
			continue L
		case ",dive":
			k := rf.Kind()
			switch k { // nolint:exhaustive
			case reflect.Struct:
				if err := applyWithPrefix(rf, prefix); err != nil {
					return err
				}
			case reflect.Pointer:
				if rf.IsNil() {
					rf.Set(reflect.New(rf.Type().Elem()))
				}
				if err := applyWithPrefix(rf.Elem(), prefix); err != nil {
					return err
				}
			default:
				return fmt.Errorf("'dive' is not available to kind %q", k)
			}
		default:
			envKey := envKey(tagVal, prefix)
			envVal, envValSet := os.LookupEnv(envKey)
			if !envValSet {
				continue L
			}
			if err := setValue(envVal, rf); err != nil {
				return fmt.Errorf("set env key %q, env value %q: %w", envKey, envVal, err)
			}
		}
	}

	return nil
}

func envKey(envVar, prefix string) string {
	if prefix != "" {
		return fmt.Sprintf("%s_%s", prefix, envVar)
	}
	return envVar
}

func setValue(val string, rv reflect.Value) error {
	fieldKind := rv.Kind()
	fieldType := rv.Type()

	if rv.CanAddr() && fieldKind != reflect.Pointer {
		rf := rv.Addr()
		if um, ok := rf.Interface().(encoding.TextUnmarshaler); ok {
			if err := um.UnmarshalText([]byte(val)); err != nil {
				return fmt.Errorf("unmarshal value %q: %w", val, err)
			}
			return nil
		}
	}

	switch fieldKind { //nolint:exhaustive
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i, err := strconv.ParseInt(val, decadicBase, int(fieldType.Size()))
		if err != nil {
			return err
		}
		rv.SetInt(i)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		u, err := strconv.ParseUint(val, decadicBase, int(fieldType.Size()))
		if err != nil {
			return err
		}
		rv.SetUint(u)
	case reflect.Float32, reflect.Float64:
		f, err := strconv.ParseFloat(val, int(fieldType.Size()))
		if err != nil {
			return err
		}
		rv.SetFloat(f)
	case reflect.Bool:
		b, err := strconv.ParseBool(val)
		if err != nil {
			return err
		}
		rv.SetBool(b)
	case reflect.String:
		rv.SetString(val)
	case reflect.Pointer:
		typ := fieldType.Elem()
		ptr := reflect.New(typ)
		if err := setValue(val, ptr.Elem()); err != nil {
			return fmt.Errorf("set value %q: %w", val, err)
		}
		rv.Set(ptr)
	default:
		// This case is not recommended, but it is allowed.
		// This case is hit when a field type defines encoding.TextUnmashaler on value receiver.
		// It is considered resonsibility of the caller to make sure this is the intended use-case.
		um, ok := rv.Interface().(encoding.TextUnmarshaler)
		if !ok {
			return fmt.Errorf("field type %q does not implement encoding.TextUnmarshaler interface", fieldType)
		}
		if err := um.UnmarshalText([]byte(val)); err != nil {
			return fmt.Errorf("unmarshal value %q: %w", val, err)
		}
	}
	return nil

}
