package util

import (
	"errors"
	"reflect"
)

func Call(m map[string]interface{}, funcName string, params ...interface{}) (result []reflect.Value,
	err error)  {
	f := reflect.ValueOf(m[funcName])

	// compare the number of func arguments
	if len(params) != f.Type().NumIn() {
		err = errors.New("The number of params is not adapted.")
		return
	}

	// make args for func from params interface
	in := make([]reflect.Value,len(params))
	for k, p := range params {
		in[k] =  reflect.ValueOf(p)
	}

	result = f.Call(in)
	return
}