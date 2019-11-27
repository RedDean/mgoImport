package id

import (
	"gopkg.in/mgo.v2/bson"
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
	pool := NewUpdateBatchPool(DefaultPoolSize, w.collection)
	defer func() {
		pool.Clean()
		swg.Done()
	}()

	for {
		data, ok := <-dataCh
		if !ok {
			return
		}
		pool.Add(w.buildUpdateOpsObj(data.(string)))
	}
}

func (w IDWorker) buildUpdateOpsObj(originalId string) *updateObj {
	return &updateObj{
		selector: bson.M{
			w.targetColumn: originalId,
		},
		setter: bson.M{
			"$set": bson.M{
				w.targetColumn: bson.NewObjectId(),
			},
		},
	}
}

//func (w IDWorker) updateID(originalId string) error {
//	session := mgoImport.GetDb()
//	defer session.Close()
//
//	where := bson.M{
//		w.targetColumn: originalId,
//	}
//
//	set := bson.M{
//		"$set": bson.M{
//			w.targetColumn: bson.NewObjectId(),
//		},
//	}
//
//	_, err := session.DB(mgoImport.G_DBname).C(w.collection).UpdateAll(where, set)
//	if err != nil {
//		return err
//	}
//
//	return nil
//}
