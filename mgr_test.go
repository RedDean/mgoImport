package mgoImport

import (
	"strings"
	"testing"
)

func TestMgr(t *testing.T) {

	input := strings.NewReader("harden|{\"team\":\"rocket\"}\nrussell|{\"team\":\"rocket\"}\n")
	database := &Repository{}
	_ = database.buildProperties([]string{"name", "payload"}, []string{"string", "json"})

	t.Logf("properties %v", database.Properties)

	mgr := &Mgr{
		NewDataParser(input, 2, "|"),
		database,
		2,
	}

	mgr.Run(1)

	got := database.DbName
	want := []map[string]interface{}{
		{
			"name": "harden",
			"team": "rocket",
		},
		{
			"name": "russell",
			"team": "rocket",
		},
	}

	assertTwoObjEqual(t, got, want)

}
