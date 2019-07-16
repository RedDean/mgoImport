package mgoImport

import (
	"io"
	"reflect"
	"testing"
)

/*
   common testing helper function
*/

func assertNoError(t *testing.T, err error) {
	t.Helper()

	if err == io.EOF {
		return
	}

	if err != nil {
		t.Fatalf("didn't expect an error but got one, %v", err)
	}
}

func assertTwoObjEqual(t *testing.T, got, want interface{}) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v doesn't equal want %v", got, want)
	}
}
