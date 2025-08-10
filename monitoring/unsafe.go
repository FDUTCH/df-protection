package monitoring

import (
	"reflect"
	"unsafe"
)

// Credit: https://github.com/bedrock-gophers/unsafe/blob/master/unsafe/internal.go

func fetchPrivateField(s any, name string) reflect.Value {
	reflectedValue := reflect.ValueOf(s).Elem()
	privateFieldValue := reflectedValue.FieldByName(name)
	return reflect.NewAt(privateFieldValue.Type(), unsafe.Pointer(privateFieldValue.UnsafeAddr())).Elem()
}

func updatePrivateField(v any, name string, value any) {
	reflectedValue := reflect.ValueOf(v).Elem()
	privateFieldValue := reflectedValue.FieldByName(name)
	privateFieldValue = reflect.NewAt(privateFieldValue.Type(), unsafe.Pointer(privateFieldValue.UnsafeAddr())).Elem()
	privateFieldValue.Set(reflect.ValueOf(value))
}
