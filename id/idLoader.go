package id

import (
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"mgoImport"
	"runtime"
)

type IDloader struct {
	targetCollection string
	targetColumn     string

	data       []interface{}
	dataLength int
}

func NewIDLoader(collection string, column string) *IDloader {
	return &IDloader{
		targetColumn:     column,
		targetCollection: collection,
		data:             make([]interface{}, 0),
	}
}

func (l *IDloader) Load() {

	db := mgoImport.GetDb()
	defer db.Close()

	query := []bson.M{
		{"$match": bson.M{
			l.targetColumn: bson.M{"$ne": ""},
		},
		},
		{
			"$group": bson.M{
				"_id":        "$" + l.targetColumn,
				"groupCount": bson.M{"$sum": 1},
			},
		},
		{
			"$project": bson.M{
				"_id": 1,
			},
		},
	}

	pipe := db.DB(mgoImport.G_DBname).C(l.targetCollection).Pipe(query)

	var results []struct {
		ID string `bson:"_id"`
	}

	if err := pipe.All(&results); err != nil {
		panic(fmt.Errorf("load err failed! %w", err))
	}

	for key := range results {
		l.data = append(l.data, results[key].ID)
	}

	l.dataLength = len(l.data)
}

func (l IDloader) GetData() []interface{} {
	return l.data
}

func DivideIntoSmallChunks(data []interface{}) [][]interface{} {
	dataLength := len(data)
	if dataLength == 0 {
		fmt.Println("[ERROR] No data to divide.")
		return nil
	}
	//
	var (
		divided [][]interface{}
		numCPU  int
	)

	numCPU = runtime.NumCPU()
	chunkSize := (dataLength + numCPU - 1) / numCPU
	// for test mock
	// chunkSize := 4

	for i := 0; i < dataLength; i += chunkSize {
		end := i + chunkSize

		if end > dataLength {
			end = dataLength
		}

		divided = append(divided, data[i:end])
	}

	return divided
}
