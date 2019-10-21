package id

import (
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"mgoImport"
	"sync"
)

type IDWorker struct {
	collection   string
	targetColumn string
}

func NewIDWorker(collection, column string) *IDWorker {
	return &IDWorker{
		collection:   collection,
		targetColumn: column,
	}
}

func (w IDWorker) Do(dataCh <-chan interface{}, swg *sync.WaitGroup) {
	defer swg.Done()
	for {
		data, ok := <-dataCh
		if !ok {
			return
		}

		if err := w.updateID(data.(string)); err != nil {
			fmt.Printf("[ERROR] error : %v ouccred when update id: %s", err, data.(string))
			continue
		}
	}
}

func (w IDWorker) updateID(originalId string) error {
	session := mgoImport.GetDb()
	defer session.Close()

	where := bson.M{
		w.targetColumn: originalId,
	}

	set := bson.M{
		"$set": bson.M{
			w.targetColumn: bson.NewObjectId(),
		},
	}

	_, err := session.DB(mgoImport.G_DBname).C(w.collection).UpdateAll(where, set)
	if err != nil {
		return err
	}

	return nil
}
