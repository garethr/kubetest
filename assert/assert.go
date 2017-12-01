package assert

// The following assertions are ported directly from https://github.com/stretchr/testify
// with minor modifications to allow use outside the Go test framework

// Copyright (c) 2012 - 2013 Mat Ryer and Tyler Bunnell
//
// Please consider promoting this project if you find it useful.
//
// Permission is hereby granted, free of charge, to any person
// obtaining a copy of this software and associated documentation
// files (the "Software"), to deal in the Software without restriction,
// including without limitation the rights to use, copy, modify, merge,
// publish, distribute, sublicense, and/or sell copies of the Software,
// and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included
// in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
// OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
// IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM,
// DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT
// OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE
// OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
//

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"
)

const (
	AssertionError   = iota
	AssertionFailure = iota
	AssertionSuccess = iota
)

var Results []assertionResult

type assertionResult struct {
	Message string
	Kind    int
}

func NotEqual(actual, expected interface{}, msg string) bool {
	if err := validateEqualArgs(expected, actual); err != nil {
		result := assertionResult{Message: fmt.Sprintf("Invalid operation: %#v == %#v (%s)",
			expected, actual, err), Kind: AssertionError}
		Results = append(Results, result)
		return false
	}

	if objectsAreEqual(expected, actual) {
		result := assertionResult{Message: fmt.Sprintf("%s but does. actual: %s", msg, actual),
			Kind: AssertionFailure}
		Results = append(Results, result)
		return false
	}

	result := assertionResult{Message: msg, Kind: AssertionSuccess}
	Results = append(Results, result)
	return true
}

func Equal(actual, expected interface{}, msg string) bool {
	if err := validateEqualArgs(expected, actual); err != nil {
		result := assertionResult{Message: fmt.Sprintf("Invalid operation: %#v == %#v (%s)",
			expected, actual, err), Kind: AssertionError}
		Results = append(Results, result)
		return false
	}

	if !objectsAreEqual(expected, actual) {
		expected, actual = formatUnequalValues(expected, actual)
		result := assertionResult{Message: fmt.Sprintf("%s but doesn't. expected: %s actual: %s", msg, expected, actual),
			Kind: AssertionFailure}
		Results = append(Results, result)
		return false
	} else {
		result := assertionResult{Message: msg, Kind: AssertionSuccess}
		Results = append(Results, result)
		return true
	}
}

func validateEqualArgs(expected, actual interface{}) error {
	if isFunction(expected) || isFunction(actual) {
		return errors.New("cannot take func type as argument")
	}
	return nil
}

func isFunction(arg interface{}) bool {
	if arg == nil {
		return false
	}
	return reflect.TypeOf(arg).Kind() == reflect.Func
}

func objectsAreEqual(expected, actual interface{}) bool {
	if expected == nil || actual == nil {
		return expected == actual
	}
	if exp, ok := expected.([]byte); ok {
		act, ok := actual.([]byte)
		if !ok {
			return false
		} else if exp == nil || act == nil {
			return exp == nil && act == nil
		}
		return bytes.Equal(exp, act)
	}
	return reflect.DeepEqual(expected, actual)
}

func formatUnequalValues(expected, actual interface{}) (e string, a string) {
	if reflect.TypeOf(expected) != reflect.TypeOf(actual) {
		return fmt.Sprintf("%T(%#v)", expected, expected),
			fmt.Sprintf("%T(%#v)", actual, actual)
	}

	return fmt.Sprintf("%#v", expected),
		fmt.Sprintf("%#v", actual)
}

func typeAndKind(v interface{}) (reflect.Type, reflect.Kind) {
	t := reflect.TypeOf(v)
	k := t.Kind()

	if k == reflect.Ptr {
		t = t.Elem()
		k = t.Kind()
	}
	return t, k
}

func includeElement(list interface{}, element interface{}) (ok, found bool) {

	listValue := reflect.ValueOf(list)
	elementValue := reflect.ValueOf(element)
	defer func() {
		if e := recover(); e != nil {
			ok = false
			found = false
		}
	}()

	if reflect.TypeOf(list).Kind() == reflect.String {
		return true, strings.Contains(listValue.String(), elementValue.String())
	}

	if reflect.TypeOf(list).Kind() == reflect.Map {
		mapKeys := listValue.MapKeys()
		for i := 0; i < len(mapKeys); i++ {
			if objectsAreEqual(mapKeys[i].Interface(), element) {
				return true, true
			}
		}
		return true, false
	}

	for i := 0; i < listValue.Len(); i++ {
		if objectsAreEqual(listValue.Index(i).Interface(), element) {
			return true, true
		}
	}
	return true, false

}

func Fail(msg string) bool {
	result := assertionResult{Message: msg, Kind: AssertionFailure}
	Results = append(Results, result)
	return false
}

func FailNow(msg string) bool {
	result := assertionResult{Message: msg, Kind: AssertionError}
	Results = append(Results, result)
	return false
}

func NotNil(object interface{}, msg string) bool {
	if !isNil(object) {
		result := assertionResult{Message: msg, Kind: AssertionSuccess}
		Results = append(Results, result)
		return true
	}
	result := assertionResult{Message: fmt.Sprintf("%s Expected value not to be nil", msg),
		Kind: AssertionFailure}
	Results = append(Results, result)
	return false
}

// isNil checks if a specified object is nil or not, without Failing.
func isNil(object interface{}) bool {
	if object == nil {
		return true
	}

	value := reflect.ValueOf(object)
	kind := value.Kind()
	if kind >= reflect.Chan && kind <= reflect.Slice && value.IsNil() {
		return true
	}

	return false
}

// Nil asserts that the specified object is nil.
//
//    assert.Nil(t, err)
//
// Returns whether the assertion was successful (true) or not (false).
func Nil(object interface{}, msg string) bool {
	if isNil(object) {
		result := assertionResult{Message: msg, Kind: AssertionSuccess}
		Results = append(Results, result)
		return true
	}
	result := assertionResult{Message: fmt.Sprintf("%s Expected nil, but got: %#v", msg, object),
		Kind: AssertionFailure}
	Results = append(Results, result)
	return false
}

var numericZeros = []interface{}{
	int(0),
	int8(0),
	int16(0),
	int32(0),
	int64(0),
	uint(0),
	uint8(0),
	uint16(0),
	uint32(0),
	uint64(0),
	float32(0),
	float64(0),
}

// isEmpty gets whether the specified object is considered empty or not.
func isEmpty(object interface{}) bool {

	if object == nil {
		return true
	} else if object == "" {
		return true
	} else if object == false {
		return true
	}

	for _, v := range numericZeros {
		if object == v {
			return true
		}
	}

	objValue := reflect.ValueOf(object)

	switch objValue.Kind() {
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice, reflect.String:
		{
			return (objValue.Len() == 0)
		}
	case reflect.Struct:
		switch object.(type) {
		case time.Time:
			return object.(time.Time).IsZero()
		}
	case reflect.Ptr:
		{
			if objValue.IsNil() {
				return true
			}
			switch object.(type) {
			case *time.Time:
				return object.(*time.Time).IsZero()
			default:
				return false
			}
		}
	}
	return false
}

// Empty asserts that the specified object is empty.  I.e. nil, "", false, 0 or either
// a slice or a channel with len == 0.
//
//  assert.Empty(t, obj)
//
// Returns whether the assertion was successful (true) or not (false).
func Empty(object interface{}, msg string) bool {

	pass := isEmpty(object)
	var result assertionResult
	if !pass {
		result = assertionResult{Message: fmt.Sprintf("%s Should be empty, but was %v", msg, object),
			Kind: AssertionError}
	} else {
		result = assertionResult{Message: msg, Kind: AssertionSuccess}
	}
	Results = append(Results, result)
	return pass

}

// NotEmpty asserts that the specified object is NOT empty.  I.e. not nil, "", false, 0 or either
// a slice or a channel with len == 0.
//
//  if assert.NotEmpty(t, obj) {
//    assert.Equal(t, "two", obj[1])
//  }
//
// Returns whether the assertion was successful (true) or not (false).
func NotEmpty(object interface{}, msg string) bool {

	pass := !isEmpty(object)
	if !pass {
		result := assertionResult{Message: fmt.Sprintf("%s Should NOT be empty, but was %v", msg, object),
			Kind: AssertionError}
		Results = append(Results, result)
		return false
	}

	result := assertionResult{Message: msg, Kind: AssertionSuccess}
	Results = append(Results, result)
	return pass
}

func True(value bool, msg string) bool {
	if value != true {
		result := assertionResult{Message: msg, Kind: AssertionFailure}
		Results = append(Results, result)
		return false
	}
	result := assertionResult{Message: msg, Kind: AssertionSuccess}
	Results = append(Results, result)
	return true
}

func False(value bool, msg string) bool {
	if value != false {
		result := assertionResult{Message: msg, Kind: AssertionFailure}
		Results = append(Results, result)
		return false
	}
	result := assertionResult{Message: msg, Kind: AssertionSuccess}
	Results = append(Results, result)
	return true
}

func Contains(s, contains interface{}, msg string) bool {
	ok, found := includeElement(s, contains)
	if !ok {
		result := assertionResult{Message: fmt.Sprintf("An error occured with %s", msg),
			Kind: AssertionError}
		Results = append(Results, result)
		return false
	}
	if !found {
		result := assertionResult{Message: fmt.Sprintf("\"%s\" does not contain \"%s\"", s, contains),
			Kind: AssertionFailure}
		Results = append(Results, result)
		return false
	}

	result := assertionResult{Message: msg, Kind: AssertionSuccess}
	Results = append(Results, result)
	return true
}

// NotContains asserts that the specified string, list(array, slice...) or map does NOT contain the
// specified substring or element.
//
//    assert.NotContains(t, "Hello World", "Earth")
//    assert.NotContains(t, ["Hello", "World"], "Earth")
//    assert.NotContains(t, {"Hello": "World"}, "Earth")
//
// Returns whether the assertion was successful (true) or not (false).
func NotContains(s, contains interface{}, msg string) bool {

	ok, found := includeElement(s, contains)
	if !ok {
		result := assertionResult{Message: fmt.Sprintf("An error occured with %s", msg),
			Kind: AssertionError}
		Results = append(Results, result)
		return false
	}
	if found {
		result := assertionResult{Message: fmt.Sprintf("\"%s\" should not contain \"%s\"", s, contains),
			Kind: AssertionFailure}
		Results = append(Results, result)
		return false
	}

	result := assertionResult{Message: msg, Kind: AssertionSuccess}
	Results = append(Results, result)
	return true
}
