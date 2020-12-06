package testutils

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
)

// MARK: Public Functions

// Diff recursively compares two values of the same type. If they are different
// types, an error is returned. If they are the same type, it returns a map of
// names of field struct fields that did not match as keys and a slice
// containing the unequal field values as values
func Diff(a, b interface{}) (interface{}, error) {
	// Types must match to compare
	aType := reflect.TypeOf(a)
	bType := reflect.TypeOf(b)
	if aType != bType {
		return map[string][]string{
				"mismatched types": {aType.String(), bType.String()},
			}, fmt.Errorf(
				"types don't match -- %s/%s",
				aType.String(),
				bType.String(),
			)
	}

	// Get type of underlying struct, in case a and b are Ptr
	// structType := aValue.Type()
	aValue := reflect.Indirect(reflect.ValueOf(a))
	aInterface := aValue.Interface()
	bInterface := reflect.Indirect(reflect.ValueOf(b)).Interface()
	compareType := aValue.Type()
	for _, t := range directlyComparable {
		if compareType.String() == t {
			if aInterface == bInterface {
				return nil, nil
			}
			return []interface{}{aInterface, bInterface}, nil
		}
	}

	switch compareType.Kind() {
	case reflect.Array:
		fallthrough
	case reflect.Slice:
		return diffSlice(a, b)
	case reflect.Map:
		return diffMap(a, b)
	case reflect.Struct:
		return diffStruct(a, b)
	default:
		if aInterface == bInterface {
			return nil, nil
		}
		return []interface{}{aInterface, bInterface}, nil
	}
}

// DiffJSON returns the results of Diff(a, b) marshalled to JSON
func DiffJSON(a, b interface{}) (string, error) {
	diff, diffErr := Diff(a, b)
	bytes, err := json.Marshal(diff)
	if err != nil {
		return "", err
	}
	return string(bytes), diffErr
}

// DiffJSON returns the results of Diff(a, b) marshalled to JSON with the
// provided prefix and indent strings
func DiffJSONIndent(a, b interface{}, prefix, indent string) (string, error) {
	diff, diffErr := Diff(a, b)
	bytes, err := json.MarshalIndent(diff, prefix, indent)
	if err != nil {
		return "", err
	}
	return string(bytes), diffErr
}

// MARK: Private Functions

func diffSlice(a, b interface{}) (interface{}, error) {
	diff := make(map[int]interface{})
	aValue := reflect.Indirect(reflect.ValueOf(a))
	bValue := reflect.Indirect(reflect.ValueOf(b))
	aLen := aValue.Len()
	bLen := bValue.Len()
	var minLen, maxLen int
	if aLen > bLen {
		minLen = bLen; maxLen = aLen
	} else {
		minLen = aLen; maxLen = bLen
	}
	for i := 0; i < maxLen; i++ {
		if i >= minLen {
			if aLen > bLen {
				diff[i] = []interface{}{aValue.Index(i).Interface(), nil}
			} else {
				diff[i] = []interface{}{nil, bValue.Index(i).Interface()}
			}
		} else {
			innerDiff, err := Diff(aValue.Index(i).Interface(), bValue.Index(i).Interface())
			if err != nil && innerDiff == nil {
				diff[i] = err
			}
			if innerDiff != nil {
				diff[i] = innerDiff
			}
		}
	}
	if len(diff) > 0 {
		return diff, nil
	}
	return nil, nil
}

func diffMap(a, b interface{}) (interface{}, error) {
	aValue := reflect.Indirect(reflect.ValueOf(a))
	bValue := reflect.Indirect(reflect.ValueOf(b))
	diff := make(map[interface{}]interface{})

	if aValue.Len() > 0 {
		// Add all a's keys and values to the diff map
		m := aValue.MapRange()
		n := true
		for ; n; n = m.Next() {
			diff[m.Key().Interface()] = []interface{}{m.Value().Interface(), nil}
		}
	}

	if bValue.Len() > 0 {
		// Check all of b's values
		m := bValue.MapRange()
		n := true
		for ; n; n = m.Next() {
			if v, e := diff[m.Key().Interface()]; e {
				// If b's key exists in the diff map, compare values
				innerDiff, err := Diff(v, m.Value().Interface())
				if err != nil {
					// If comparison fails, store error
					diff[m.Key().Interface()] = err
				} else if innerDiff != nil{
					// if values differ, store diff
					diff[m.Key().Interface()] = innerDiff
				} else {
					// if values are equal, remove the diff
					delete(diff, m.Key().Interface())
				}
			} else {
				// if key does not exist, store b's value
				diff[m.Key().Interface()] = []interface{}{nil, m.Value().Interface()}
			}
		}
	}
	if len(diff) > 0 {
		return diff, nil
	}
	return nil, nil
}

func diffStruct(a, b interface{}) (interface{}, error) {
	// Types must be structs
	aValue := reflect.Indirect(reflect.ValueOf(a))
	if aValue.Kind() != reflect.Struct {
		return nil, errors.New(reflect.TypeOf(a).Name() + "is not a struct")
	}

	// Get type of underlying struct, in case a and b are Ptr
	structType := aValue.Type()

	diff := make(map[string]interface{})
	bValue := reflect.Indirect(reflect.ValueOf(b))
	fieldCount := aValue.NumField()
	for i := 0; i < fieldCount; i++ {
		structField := structType.Field(i)
		if structField.PkgPath != "" {
			continue
		}
		aField := reflect.Indirect(aValue.Field(i))
		bField := reflect.Indirect(bValue.Field(i))
		innerDiff, err := Diff(aField.Interface(), bField.Interface())
		if err != nil {
			diff[structField.Name] = err
		}
		if innerDiff != nil {
			diff[structField.Name] = innerDiff
		}
	}
	if len(diff) > 0 {
		return diff, nil
	}
	return nil, nil
}

// MARK: Private Variables

var directlyComparable = [...]string{
	"time.Time",
	"uuid.UUID",
}
