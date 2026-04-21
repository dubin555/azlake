package config

import (
	"reflect"
	"strings"
)

const sep = "."

// GetStructKeys returns all keys in a nested struct type.
func GetStructKeys(typ reflect.Type, tag, squashValue string) []string {
	return appendStructKeys(typ, tag, ","+squashValue, nil, nil)
}

func appendStructKeys(typ reflect.Type, tag, squashValue string, prefix []string, keys []string) []string {
	for ; typ.Kind() == reflect.Pointer; typ = typ.Elem() {
	}
	if typ.Kind() != reflect.Struct {
		return append(keys, strings.Join(prefix, sep))
	}
	for i := 0; i < typ.NumField(); i++ {
		fieldType := typ.Field(i)
		var (
			fieldName string
			squash    bool
			ok        bool
		)
		if fieldName, ok = fieldType.Tag.Lookup(tag); ok {
			if strings.HasSuffix(fieldName, squashValue) {
				squash = true
				fieldName = strings.TrimSuffix(fieldName, squashValue)
			}
		} else {
			fieldName = strings.ToLower(fieldType.Name)
		}
		if !squash {
			prefix = append(prefix, fieldName)
		}
		keys = appendStructKeys(fieldType.Type, tag, squashValue, prefix, keys)
		if !squash {
			prefix = prefix[:len(prefix)-1]
		}
	}
	return keys
}

// ValidateMissingRequiredKeys returns all keys that have a required tag but are unset.
func ValidateMissingRequiredKeys(value any, tag, squashValue string) []string {
	return appendStructKeysIfZero(reflect.ValueOf(value), tag, ","+squashValue, "validate", "required", nil, nil)
}

func isScalar(kind reflect.Kind) bool {
	switch kind {
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Pointer, reflect.Slice:
		return false
	}
	return true
}

func appendStructKeysIfZero(value reflect.Value, tag, squashValue, validateTag, requiredValue string, prefix []string, keys []string) []string {
	for value.Kind() == reflect.Pointer {
		if value.IsZero() {
			return keys
		}
		value = value.Elem()
	}
	if !isScalar(value.Kind()) {
		if !isScalar(value.Type().Elem().Kind()) {
			panic("No support for detecting required keys inside " + value.Kind().String() + " of structs")
		}
	}
	if value.Kind() != reflect.Struct {
		return keys
	}
	for i := 0; i < value.NumField(); i++ {
		fieldType := value.Type().Field(i)
		fieldValue := value.Field(i)
		var (
			fieldName string
			squash    bool
			ok        bool
		)
		if fieldName, ok = fieldType.Tag.Lookup(tag); ok {
			if strings.HasSuffix(fieldName, squashValue) {
				squash = true
				fieldName = strings.TrimSuffix(fieldName, squashValue)
			}
		} else {
			fieldName = strings.ToLower(fieldType.Name)
		}
		if validationsString, ok := fieldType.Tag.Lookup(validateTag); ok {
			for validation := range strings.SplitSeq(validationsString, ",") {
				if validation == requiredValue && fieldValue.IsZero() {
					keys = append(keys, strings.Join(append(prefix, fieldName), sep))
				}
			}
		}
		if !squash {
			prefix = append(prefix, fieldName)
		}
		keys = appendStructKeysIfZero(fieldValue, tag, squashValue, validateTag, requiredValue, prefix, keys)
		if !squash {
			prefix = prefix[:len(prefix)-1]
		}
	}
	return keys
}
