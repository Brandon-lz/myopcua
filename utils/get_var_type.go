package utils

import (
	"reflect"
)
func Typeof(v interface{}) string {
    return reflect.TypeOf(v).String()    
}