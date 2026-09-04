package constants

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type Enum[T any] interface {
	HasValue[T]
	HasString
	HasValidator
}

type HasValidator interface {
	Valid() bool
}

type HasString interface {
	String() string
}

type HasValue[T any] interface {
	Value() T
}

// ResolveEnum attempts to parse a string into an enum value of type T
func ResolveEnum[T HasValidator](raw string) (T, error) {
	var zero T

	enumType := reflect.TypeFor[T]()
	enumName := enumType.Name()

	trimmed := strings.TrimSpace(raw)
	value := reflect.New(enumType).Elem()

	switch enumType.Kind() {
	case reflect.String:
		value.SetString(raw)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		parsed, err := strconv.ParseInt(trimmed, 10, enumType.Bits())
		if err != nil {
			return zero, err
		}
		value.SetInt(parsed)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		parsed, err := strconv.ParseUint(trimmed, 10, enumType.Bits())
		if err != nil {
			return zero, err
		}
		value.SetUint(parsed)
	default:
		return zero, fmt.Errorf("unsupported enum type '%s' for %s", enumType, enumName)
	}

	enum := value.Interface().(T)
	if !enum.Valid() {
		return zero, fmt.Errorf("invalid enum value '%q' for %s", trimmed, enumName)
	}
	return enum, nil
}
