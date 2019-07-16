package mgoImport

import (
	"testing"
)

func TestRepository(t *testing.T) {

	//want := map[string]interface{}{
	//    "name":"harden",
	//	"age": 23,
	//	"tel": 10086,
	//}

	t.Run("build properties map", func(t *testing.T) {

		got := &Repository{}
		err := got.BuildProperties([]string{"age", "name"}, []string{"int", "string"})
		assertNoError(t, err)

		want := &Repository{
			Properties: []Model{
				{
					FieldType: "int",
					FieldName: "age",
				}, {
					FieldName: "name",
					FieldType: "string",
				},
			},
		}

		assertTwoObjEqual(t, got, want)

	})

	t.Run("build mongo model", func(t *testing.T) {

		repo := &Repository{}
		err1 := repo.BuildProperties([]string{"name", "number", "payload"}, []string{"string", "int", "json"})
		assertNoError(t, err1)

		input := []string{"harden", "13", "{\"team\":\"H-town\",\"age\":29}"}

		got, err2 := repo.BuildModel(input)
		assertNoError(t, err2)
		want := map[string]interface{}{
			"name":   "harden",
			"number": 13,
			"team":   "H-town",
			"age":    29,
		}

		assertTwoObjEqual(t, got, want)
	})

}
