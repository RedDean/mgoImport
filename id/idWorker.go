package id

import (
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"mgoImport"
	"sync"
	"time"
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

		// There will be some god damn socket errors in my local environment.
		// Maybe it related to docker configure or concurrency.
		// Retry 3 times when capture error.
		for i := 1; i <= 3; i++ {
			if err := w.updateID(data.(string)); err == nil {
				break
			} else {
				time.Sleep(time.Millisecond * time.Duration(i*100))
				fmt.Printf("[ERROR] error : %v ouccred when update id: %s. Retry at %d times \n", err, data.(string), i)
			}
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
