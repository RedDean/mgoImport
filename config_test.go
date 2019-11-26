package mgoImport

import (
	"mgoImport/testUtil"
	"os"
	"reflect"
	"testing"
)

func TestConfigFile(t *testing.T) {
	t.Run("read configure file by given directory", func(t *testing.T) {
		cf := &ConfigFile{}
		dir := "./cmd/app/config.json"

		file, err := os.OpenFile(dir, os.O_RDONLY, 0666)
		testUtil.AssertNoError(t, err)

		if err := cf.load(file); err != nil {
			t.Fatalf("load json err: %v", err)
		}

		got := cf.JsonField
		want := map[string]string{
			"gameCount":       "int",
			"activeGameCount": "int",
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v , want %v", got, want)
		}
	})
}
