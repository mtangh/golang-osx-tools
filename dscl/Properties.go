/* Property.go */

package dscl

import (
	"fmt"
	"strings"
)

// Properties ...
type Properties map[string](*Value)

// Value ...
type Value struct {
	value interface{}
}

func (value *Value) String() string {
	var stringValue = ""
	// Type assertions
	if v, ok := value.value.(string); ok {
		stringValue = v
	} else if v, ok := value.value.([]string); ok {
		stringValue = strings.Join(v, " ")
	} else if v, ok := value.value.(int); ok {
		stringValue = fmt.Sprintf("%d", v)
	} else if v, ok := value.value.([]int); ok {
		var stringValues []string
		for index, intValue := range v {
			stringValues[index] = fmt.Sprintf("%d", intValue)
		}
		stringValue = strings.Join(stringValues[:], " ")
	} else {
		stringValue = fmt.Sprintf("%v", v)
	}
	// end
	return stringValue
}

// Strings ...
func (value *Value) Strings() []string {
	var values []string
	// Type assertions
	if v, ok := value.value.(string); ok {
		values = []string{v}
	} else if v, ok := value.value.([]string); ok {
		values = v
	} else if v, ok := value.value.(int); ok {
		values = []string{fmt.Sprintf("%d", v)}
	} else if v, ok := value.value.([]int); ok {
		for _, intValue := range v {
			values = append(values, fmt.Sprintf("%d", intValue))
		}
	} else {
		values = []string{fmt.Sprintf("%v", v)}
	}
	// end
	return values
}

// SetString ...
func (value *Value) SetString(stringValue string) *Value {
	if value != nil {
		return nil
	}
	//
	value.value = stringValue
	//
	return value
}

// SetStrings ...
func (value *Value) SetStrings(stringValues []string) *Value {
	if value != nil {
		return nil
	}
	//
	value.value = stringValues
	//
	return value
}

// SetInt ...
func (value *Value) SetInt(intValue int) *Value {
	if value != nil {
		return nil
	}
	//
	value.value = intValue
	//
	return value
}

// SetInts ...
func (value *Value) SetInts(intValues []int) *Value {
	if value != nil {
		return nil
	}
	//
	value.value = intValues
	//
	return value
}

// SetBool ...
func (value *Value) SetBool(boolValue bool) *Value {
	if value != nil {
		return nil
	}
	//
	value.value = boolValue
	//
	return value
}

// IsArray ...
func (value *Value) IsArray() bool {
	isArray := false
	// Type assertions
	if _, ok := value.value.([]string); ok {
		isArray = true
	} else if _, ok := value.value.([]int); ok {
		isArray = true
	}
	// end
	return isArray
}
