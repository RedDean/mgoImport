package mgoImport

import (
	"reflect"
	"strings"
	"testing"
)

func TestDataParser(t *testing.T) {

	t.Run("testing parser data line by line", func(t *testing.T) {
		file := strings.NewReader(`hello
world
tdd`)
		parser := NewDataParser(file,3)
		err := parser.readLine()
		assertNoError(t, err)


		want := []string{"hello", "world", "tdd"}
		got := make([]string,0)
		for v := range parser.DataCh {
			got = append(got, v)
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %s , want %s", got, want)
		}
	})

	t.Run("split csv string into string slice by given delimiter", func(t *testing.T) {
		file := strings.NewReader("123|456|789\n")
		parser := NewDataParser(file,1)

		_ = parser.readLine()

		got, err := splitByDelimiter(<-parser.DataCh, "|")
		assertNoError(t, err)

		want := []string{"123", "456", "789"}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %s , want %s", got, want)
		}
	})

	t.Run("given wrong delimiter", func(t *testing.T) {
		file := strings.NewReader("123|456|789\n")
		parser := NewDataParser(file,1)

		_ = parser.readLine()
		str :=  <-parser.DataCh
		if _, err := splitByDelimiter(str, ""); err == nil {
			t.Error("empty string ,expect an error here but got nothing")
		}

		if _, err := splitByDelimiter(str, "@@"); err == nil {
			t.Error("expect an error here but got nothing")
		}
	})

}
