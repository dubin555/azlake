package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-viper/mapstructure/v2"
)

// Strings is a []string that mapstructure can deserialize from a single string or from a list
// of strings.
type Strings []string

var (
	ourStringsType  = reflect.TypeFor[Strings]()
	stringType      = reflect.TypeFor[string]()
	stringSliceType = reflect.TypeFor[[]string]()

	ErrInvalidKeyValuePair = errors.New("invalid key-value pair")
)

// DecodeStrings is a mapstructure.HookFuncType that decodes a single string value or a slice
// of strings into Strings.
func DecodeStrings(fromValue reflect.Value, toValue reflect.Value) (any, error) {
	if toValue.Type() != ourStringsType {
		return fromValue.Interface(), nil
	}
	if fromValue.Type() == stringSliceType {
		return Strings(fromValue.Interface().([]string)), nil
	}
	if fromValue.Type() == stringType {
		return Strings(strings.Split(fromValue.String(), ",")), nil
	}
	return fromValue.Interface(), nil
}

type SecureString string

// String returns an elided version.  It is safe to call for logging.
func (SecureString) String() string {
	return "[SECRET]"
}

// SecureValue returns the actual value of s as a string.
func (s SecureString) SecureValue() string {
	return string(s)
}

func (s SecureString) MarshalText() ([]byte, error) {
	if string(s) == "" {
		return []byte(""), nil
	}
	return []byte("[SECRET]"), nil
}

// OnlyString is a string that can deserialize only from a string.
type OnlyString string

var (
	onlyStringType  = reflect.TypeFor[OnlyString]()
	ErrMustBeString = errors.New("must be a string")
)

func (o OnlyString) String() string {
	return string(o)
}

func DecodeOnlyString(fromValue reflect.Value, toValue reflect.Value) (any, error) {
	if toValue.Type() != onlyStringType {
		return fromValue.Interface(), nil
	}
	if fromValue.Type() != stringType {
		return nil, fmt.Errorf("%w, not a %s", ErrMustBeString, fromValue.Type().String())
	}
	return OnlyString(fromValue.Interface().(string)), nil
}

func DecodeStringToMap() mapstructure.DecodeHookFunc {
	return func(f reflect.Kind, t reflect.Kind, data any) (any, error) {
		if f != reflect.String || t != reflect.Map {
			return data, nil
		}
		if t != reflect.TypeFor[map[string]string]().Kind() {
			return data, nil
		}
		raw := data.(string)
		if raw == "" {
			return map[string]string{}, nil
		}
		const pairSep = ","
		const valueSep = "="
		pairs := strings.Split(raw, pairSep)
		m := make(map[string]string, len(pairs))
		for _, pair := range pairs {
			key, value, found := strings.Cut(pair, valueSep)
			if !found {
				return nil, fmt.Errorf("%w: %s", ErrInvalidKeyValuePair, pair)
			}
			m[strings.TrimSpace(key)] = strings.TrimSpace(value)
		}
		return m, nil
	}
}

func StringToSliceWithBracketHookFunc() mapstructure.DecodeHookFunc {
	return func(f reflect.Kind, t reflect.Kind, data any) (any, error) {
		if f != reflect.String || t != reflect.Slice {
			return data, nil
		}
		raw := data.(string)
		if raw == "" {
			return []string{}, nil
		}
		var result any
		err := json.Unmarshal([]byte(raw), &result)
		if err != nil {
			return data, nil
		}
		if reflect.TypeOf(result).Kind() != t {
			return data, nil
		}
		return result, nil
	}
}

func StringToStructHookFunc() mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data any) (any, error) {
		if f.Kind() != reflect.String ||
			(t.Kind() != reflect.Struct && !(t.Kind() == reflect.Pointer && t.Elem().Kind() == reflect.Struct)) {
			return data, nil
		}
		raw := data.(string)
		var val reflect.Value
		if t.Kind() == reflect.Struct {
			val = reflect.New(t)
		} else {
			val = reflect.New(t.Elem())
		}
		if raw == "" {
			return val, nil
		}
		var m map[string]any
		err := json.Unmarshal([]byte(raw), &m)
		if err != nil {
			return data, nil
		}
		return m, nil
	}
}
