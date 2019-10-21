package id

import (
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"mgoImport"
)

type ForeignKeyLoader struct {
	mainCollection string
	relatedColumn  string
	idColumn       string

	data       []interface{}
	dataLength int
}

type ForeignKeyIdObj struct {
	ID         bson.ObjectId `bson:"_id"`
	OriginalID string        `bson:"originalId"`
}

func NewForeignKeyLoader(collection string, column string, idColumn string) *ForeignKeyLoader {
	return &ForeignKeyLoader{
		idColumn:       idColumn,
		relatedColumn:  column,
		mainCollection: collection,
		data:           make([]interface{}, 0),
	}
}

func (f ForeignKeyLoader) GetData() []interface{} {
	return f.data
}

func (f *ForeignKeyLoader) Load() {
	db := mgoImport.GetDb()
	defer db.Close()

	query := []bson.M{
		{"$match": bson.M{
			f.idColumn: bson.M{"$ne": ""},
		},
		},
		{
			"$group": bson.M{
				"_id":        "$" + f.idColumn,
				"originalId": bson.M{"$first": "$" + f.relatedColumn},
			},
		},
	}

	pipe := db.DB(mgoImport.G_DBname).C(f.mainCollection).Pipe(query)

	var results []ForeignKeyIdObj

	if err := pipe.All(&results); err != nil {
		panic(fmt.Errorf("load err failed! %w", err))
	}

	for key := range results {
		f.data = append(f.data, results[key])
	}

	f.dataLength = len(f.data)
}
