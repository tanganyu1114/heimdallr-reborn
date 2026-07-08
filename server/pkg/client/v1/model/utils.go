package model

import (
	"reflect"

	httpclientv1 "github.com/ClessLi/component-base/pkg/client-sdk/http/v1"
)

// isNilBody checks if the type parameter RESP is NilBody
func isNilBody[R any]() bool {
	respType := reflect.TypeOf((*R)(nil)).Elem()
	nilType := reflect.TypeOf(httpclientv1.NilBody{})
	return respType.String() == nilType.String()
}
