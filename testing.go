package mgoImport

import (
	"io"
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
