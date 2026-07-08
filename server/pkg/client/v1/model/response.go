package model

import (
	"encoding/json"
	"reflect"

	"github.com/marmotedu/errors"
)

type ResponseBody[RESP any] struct {
	Code    int             `json:"code"`
	Data    json.RawMessage `json:"data"`
	Message string          `json:"msg"`
}

func (r *ResponseBody[RESP]) Response() (resp RESP, err error) {
	var zero RESP
	if r.Code != 0 {
		err = errors.Errorf("heimdallr response code %d, message %s", r.Code, r.Message)
		return zero, err
	}
	// Handle JSON null value for NilBody type (server returns "null" or "{}" or "" for empty response)
	if isNilBody[RESP]() && (r.Data == nil || string(r.Data) == "null" || string(r.Data) == "{}" || string(r.Data) == "") {
		return zero, nil
	}

	// Check if RESP is NilBody type but response body is not empty
	if isNilBody[RESP]() {
		return zero, errors.New("response declared as NilBody type but actual response body is not empty")
	}

	// Get the type of RESP
	respType := reflect.TypeOf((*RESP)(nil)).Elem()

	// If RESP is a pointer type, we can unmarshal directly to a new instance of it
	// If RESP is a non-pointer type, we need to unmarshal to the address of a new instance
	var result interface{}

	if respType.Kind() == reflect.Ptr {
		// RESP is a pointer type like *T
		// Create a new pointer to the underlying type
		ptrValue := reflect.New(respType.Elem())
		err = json.Unmarshal(r.Data, ptrValue.Interface())
		if err != nil {
			return zero, err
		}
		// Get the pointer value that contains our data
		result = ptrValue.Interface()
	} else {
		// RESP is a non-pointer type like T
		// Create a new instance and unmarshal to its address
		newValue := reflect.New(respType)
		err = json.Unmarshal(r.Data, newValue.Interface())
		if err != nil {
			return zero, err
		}
		// Dereference to get the actual value
		result = newValue.Elem().Interface()
	}

	resp, ok := result.(RESP)
	if !ok {
		return zero, errors.Errorf("failed to convert interface to %s", respType.String())
	}
	return resp, err
}
