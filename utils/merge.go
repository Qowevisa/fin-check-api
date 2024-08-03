package utils

import (
	"reflect"
)

func MergeNonZeroFields(src, dest interface{}) {
	srcVal := reflect.ValueOf(src).Elem()
	destVal := reflect.ValueOf(dest).Elem()

	for i := 0; i < srcVal.NumField(); i++ {
		srcField := srcVal.Field(i)
		destField := destVal.Field(i)

		if !srcField.IsZero() {
			destField.Set(srcField)
		}
	}
}
