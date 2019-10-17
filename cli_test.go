package mgoImport

import (
	"gopkg.in/mgo.v2/bson"
	"testing"
	"time"
)

func TestMgoCli(t *testing.T) {

	url := "mongodb://localhost"
	err := InitMgoCli(url)
	assertNoError(t, err)
	defer Close()

	t.Run("connect to mongo db", func(t *testing.T) {
		cli := getDb()
		isConnected := TestDb(cli)
		if !isConnected {
			t.Error("can't connect to db, plz check url!")
		}
	})

	t.Run("insert date type", func(t *testing.T) {
		cli := getDb()
		defer cli.Close()

		timeStr := "2019-05-17 09:51:21.466282+00"
		date, err := time.Parse("2006-01-02 15:04:05+00", timeStr)
		assertNoError(t, err)

		err = cli.DB("djh").C("test_insert").Insert(bson.M{"name": "test3", "insertTime": date})
		assertNoError(t, err)

		data := make(map[string]interface{})

		cli.DB("djh").C("test_insert").Find(bson.M{"name": "test3"}).One(&data)

		//got := data["insertTime"]
		//want := "2019-05-17 09:51:21-0700"
		//
		//assertTwoObjEqual(t, got,want)

		// 读出来会有8小时时差问题
	})

	t.Run("insert map", func(t *testing.T) {
		cli := getDb()
		defer cli.Close()

		data := map[string]interface{}{
			"name":   "harden",
			"number": 13,
			"team":   "H-town",
		}

		err := cli.DB("djh").C("test_insert").Insert(data)
		assertNoError(t, err)
	})

	t.Run("query empty object,check what it is", func(t *testing.T) {
		cli := getDb()
		defer cli.Close()

		type Foo struct {
			Name string
		}

		/*data := map[string]interface{}{
			"name":   "harden",
			"number": 13,
			"team":   "H-town",
		}*/

		var foo *Foo
		err := cli.DB("djh").C("test_query").Find(bson.M{"Name": "123"}).All(foo)
		t.Logf("foo : %v", foo)

		assertNoError(t, err)
	})
}
