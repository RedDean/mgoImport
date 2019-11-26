package id

import (
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"mgoImport"
)

type ChannelForeignKeyLoader struct {
	mainCollection string
	relatedColumn  string
	idColumn       string

	data       []interface{}
	dataLength int
}

type ChannelForeignKeyIdObj struct {
	ID          bson.ObjectId `bson:"_id"`
	OriginalID  string        `bson:"originalId"`
	ChannelName string        `bson:"channelName"`
}

func NewChannelForeignKeyLoader(collection string, column string, idColumn string) *ChannelForeignKeyLoader {
	return &ChannelForeignKeyLoader{
		idColumn:       idColumn,
		relatedColumn:  column,
		mainCollection: collection,
		data:           make([]interface{}, 0),
	}
}

func (f ChannelForeignKeyLoader) GetData() []interface{} {
	return f.data
}

func (f *ChannelForeignKeyLoader) Load() {
	db := mgoImport.GetDb()
	defer db.Close()

	query := []bson.M{
		{"$match": bson.M{
			f.relatedColumn: bson.M{"$ne": ""},
			"channelName":   bson.M{"$nin": []string{"GSTORE", "AMAZON", "GOOGLE"}},
		},
		},
		{
			"$project": bson.M{
				"_id":         1,
				"originalId":  1,
				"channelName": 1,
			},
		},
	}

	pipe := db.DB(mgoImport.G_DBname).C(f.mainCollection).Pipe(query)

	var results []ChannelForeignKeyIdObj

	if err := pipe.All(&results); err != nil {
		panic(fmt.Errorf("load err failed! %w", err))
	}

	for key := range results {
		f.data = append(f.data, results[key])
	}

	f.dataLength = len(f.data)
}
