package testutils

import (
	"errors"
	"fmt"
	"reflect"
)

// Diff compares two values of the same type. If they are different types, an
// error is returned. If they are the same type, it returns a map of names of
// field struct fields that did not match as keys and a slice containing the
// unequal field values as values
func Diff(a interface{}, b interface{}) (map[string][]interface{}, error) {
	// Types must match to compare
	aType := reflect.TypeOf(a)
	bType := reflect.TypeOf(b)
	if aType != bType {
		return nil, fmt.Errorf(
			"types don't match -- %s/%s",
			aType.String(),
			bType.String(),
		)
	}

	// Types must be structs
	aValue := reflect.Indirect(reflect.ValueOf(a))
	if aValue.Kind() != reflect.Struct {
		return nil, errors.New(aType.Name() + "is not a struct")
	}

	// Get type of underlying struct, in case a and b are Ptr
	structType := aValue.Type()

	diff := make(map[string][]interface{})
	bValue := reflect.Indirect(reflect.ValueOf(b))
	fieldCount := aValue.NumField()
	for i := 0; i < fieldCount; i++ {
		if structType.Field(i).PkgPath != "" {
			continue
		}
		aField := reflect.Indirect(aValue.Field(i)).Interface()
		bField := reflect.Indirect(bValue.Field(i)).Interface()
		if aField != bField {
			diff[structType.Field(i).Name] = []interface{}{aField, bField}
		}
	}
	return diff, nil
}
