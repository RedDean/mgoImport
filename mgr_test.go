package mgoImport

import (
	"strings"
	"testing"
)

func TestMgr(t *testing.T) {

	input :=  strings.NewReader("harden|{\"team\":\"rocket\"}\nrussell|{\"team\":\"rocket\"}\n")
	database := &Repository{}
	_ = database.BuildProperties([]string{"name","payload"},[]string{"string","json"})

	t.Logf("properties %v", database.Properties)

	mgr := &Mgr{
		NewDataParser(input,2),
		database,
	}

	mgr.Run()

	got := database.Db
	want := []map[string]interface{}{
		{
			"name":"harden",
			"team":"rocket",
		},
		{
			"name":"russell",
			"team":"rocket",
		},
	}

	assertTwoObjEqual(t, got , want)

}


