package mgoImport

import (
	"reflect"
	"testing"
)

func TestRepository(t *testing.T) {

	//want := map[string]interface{}{
	//    "name":"harden",
	//	"age": 23,
	//	"tel": 10086,
	//}

	t.Run("build model map", func(t *testing.T) {

		got, err := buildModelPropertiesMap([]string{"age", "name"}, []string{"int", "string"})
		assertNoError(t, err)

		want := map[string]string{"age": "int", "name": "string"}

		if !reflect.DeepEqual(got, want) {
			t.Fatalf("got %v doesn't equal want %v", got, want)
		}

	})

	//t.Run("set model value", func(t *testing.T) {
	//	got := setModelValue([]string)
	//})

}
