package check

import (
	"reflect"
	"runtime"
	"testing"
)

func Equal[T any](t *testing.T, expected T, actual T) {
	if !reflect.DeepEqual(actual, expected) {
		_, file, line, _ := runtime.Caller(1)
		t.Errorf("%s:%d expected %v but got %v", file, line, expected, actual)
	}
}

func NotEqual[T any](t *testing.T, expected T, actual T) {
	if reflect.DeepEqual(actual, expected) {
		_, file, line, _ := runtime.Caller(1)
		t.Errorf("#%s:%d should not equal %v", file, line, expected)
	}
}
