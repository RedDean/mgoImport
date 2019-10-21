package id

import (
	"gopkg.in/mgo.v2/bson"
	"mgoImport"
	"mgoImport/testUtil"
	"testing"
)

const test_load_collection = "test_loader_2"

func TestSelfModeLoad(t *testing.T) {

	url := "mongodb://localhost"
	err := mgoImport.InitMgoCli(url, "djh")
	testUtil.AssertNoError(t, err)
	defer mgoImport.Close()

	t.Run("load data from db ", func(t *testing.T) {
		// arrange
		ids := []string{"", "", "2", "2", "3", "3", "3", "3", "4", "4", "4", "4"}
		loadHelper(ids)

		want := []interface{}{"4", "3", "2"}

		// act
		loader := IDloader{
			targetCollection: test_load_collection,
			targetColumn:     "originalId",
			data:             make([]interface{}, 0),
		}
		loader.Load()

		// assert
		got := loader.GetData()
		testUtil.AssertTwoObjEqual(t, got, want)
	})

	t.Run("separate data slice into small slices", func(t *testing.T) {

		// arrange
		mock := []interface{}{1, 2, 3, 4, 5, 6, 7, 8, 9, 2, 4, 5, 5, 6}
		// NOTICE: DivideIntoSmallChunks chunkSize will calculate by numCPU in real time
		// This mock data might need reset by numCPU so that this un could work.
		want := [][]interface{}{
			{1, 2, 3, 4},
			{5, 6, 7, 8},
			{9, 2, 4, 5},
			{5, 6},
		}

		// act
		got := DivideIntoSmallChunks(mock)

		// assert
		testUtil.AssertTwoObjEqual(t, got, want)
	})

}

func loadHelper(mockData []string) {

	db := mgoImport.GetDb()
	defer db.Close()

	for _, value := range mockData {
		db.DB("djh").C(test_load_collection).Insert(bson.M{"originalId": value})
	}
}
