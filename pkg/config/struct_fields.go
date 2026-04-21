package config

import (
	"reflect"
	"strings"

	"github.com/dubin555/azlake/pkg/logging"
)

const (
	FieldMaskedValue   = "******"
	FieldMaskedNoValue = "------"
)

// MapLoggingFields returns all logging.Fields formatted based on configuration keys.
func MapLoggingFields(value any) logging.Fields {
	fields := make(logging.Fields)
	structFieldsFunc(reflect.ValueOf(value), "mapstructure", ",squash", nil, func(key string, value any) {
		fields[key] = value
	})
	return fields
}

func structFieldsFunc(value reflect.Value, tag, squashValue string, prefix []string, cb func(key string, value any)) {
	for value.Kind() == reflect.Pointer {
		if value.IsZero() {
			return
		}
		value = value.Elem()
	}
	if value.Kind() != reflect.Struct {
		key := strings.Join(prefix, sep)
		cb(key, value)
		return
	}
	for i := 0; i < value.NumField(); i++ {
		fieldType := value.Type().Field(i)
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
		fieldValue := value.Field(i)
		switch fieldValue.Interface().(type) {
		case SecureString:
			key := strings.Join(prefix, sep)
			val := FieldMaskedValue
			if fieldValue.IsZero() {
				val = FieldMaskedNoValue
			}
			cb(key, val)
		default:
			structFieldsFunc(fieldValue, tag, squashValue, prefix, cb)
		}
		if !squash {
			prefix = prefix[:len(prefix)-1]
		}
	}
}

func GetSecureStringKeyPaths(value any) []string {
	keys := []string{}
	getSecureStringKeys(reflect.ValueOf(value), "_", "mapstructure", ",squash", nil, func(key string) {
		keys = append(keys, strings.ToUpper(key))
	})
	return keys
}

func getSecureStringKeys(value reflect.Value, separator, tag, squashValue string, prefix []string, cb func(key string)) {
	for value.Kind() == reflect.Pointer {
		if value.IsZero() {
			return
		}
		value = value.Elem()
	}
	if value.Kind() != reflect.Struct {
		if value.Type() == reflect.TypeFor[SecureString]() {
			key := strings.Join(prefix, separator)
			cb(key)
		}
		return
	}
	for i := 0; i < value.NumField(); i++ {
		fieldType := value.Type().Field(i)
		var (
			fieldName string
			squash    bool
		)
		fieldName, squash = parseTag(fieldType, tag, squashValue)
		if !squash {
			prefix = append(prefix, fieldName)
		}
		fieldValue := value.Field(i)
		getSecureStringKeys(fieldValue, separator, tag, squashValue, prefix, cb)
		if !squash {
			prefix = prefix[:len(prefix)-1]
		}
	}
}

func parseTag(field reflect.StructField, tag, squashValue string) (string, bool) {
	if tagValue, ok := field.Tag.Lookup(tag); ok {
		if before, found := strings.CutSuffix(tagValue, squashValue); found {
			return before, true
		}
		return tagValue, false
	}
	return strings.ToLower(field.Name), false
}
