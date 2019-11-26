package id

import (
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"mgoImport"
	"sync"
	"time"
)

//var _StoreName  = []string{
//	"APTOIDE",
//	"ONESTORE",
//	"APPTUTTI",
//	"HTC",
//	"XIAOMISTORE",
//	"CLOUDMOOLAH",
//	"HUAWEI",
//	"JIO",
//	"SAMSUNG",
//}

type ChannelPackingIdWorker struct {
	collection    string
	foreignColumn string
}

func NewChannelPackingIdWorker(collection string, foreignColumn string) *ChannelPackingIdWorker {
	if collection == "" {
		panic("collection is empty")
	}
	return &ChannelPackingIdWorker{
		collection:    collection,
		foreignColumn: foreignColumn,
	}
}

func (f ChannelPackingIdWorker) Do(dataCh <-chan interface{}, swg *sync.WaitGroup) {
	defer swg.Done()
	for {
		data, ok := <-dataCh
		if !ok {
			return
		}
		for i := 1; i <= RETRY_TIMES; i++ {
			if err := f.updateID(data.(ChannelForeignKeyIdObj)); err == nil {
				break
			} else {
				time.Sleep(time.Second * time.Duration(i))
				fmt.Printf("[ERROR] error : %v ouccred when batchUpdate id: %s. Retry times: %d  \n", err, data.(ForeignKeyIdObj).OriginalID, i)
			}
		}
	}
}

func (f ChannelPackingIdWorker) updateID(obj ChannelForeignKeyIdObj) error {
	session := mgoImport.GetDb()
	defer session.Close()

	where := bson.M{
		"channels." + obj.ChannelName + "." + f.foreignColumn: obj.OriginalID,
	}

	set := bson.M{
		"$set": bson.M{
			"channels." + obj.ChannelName + "." + f.foreignColumn: obj.ID,
		},
	}

	_, err := session.DB(mgoImport.G_DBname).C(f.collection).UpdateAll(where, set)
	if err != nil {
		fmt.Printf("[ERROR] err:%v, can't batchUpdate store: %s, originalId: %s \n", err, obj.ChannelName, obj.OriginalID)
		return err
	}

	return nil
}
