package id

import (
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"mgoImport"
	"sync"
	"time"
)

const RETRY_TIMES = 3

type ForeignKeyWorker struct {
	collections   []string
	relatedColumn string
	foreignColumn string
}

func NewForeignKeyWorker(collections []string, relatedColumn, foreignColumn string) *ForeignKeyWorker {
	if len(collections) == 0 {
		panic(fmt.Errorf("collection slice is empty, this is required"))
	}
	return &ForeignKeyWorker{
		collections:   collections,
		relatedColumn: relatedColumn,
		foreignColumn: foreignColumn,
	}
}

func (f ForeignKeyWorker) Do(dataCh <-chan interface{}, swg *sync.WaitGroup) {
	defer swg.Done()
	for {
		data, ok := <-dataCh
		if !ok {
			return
		}

		for i := 1; i <= RETRY_TIMES; i++ {
			if err := f.updateID(data.(ForeignKeyIdObj)); err == nil {
				break
			} else {
				time.Sleep(time.Second * time.Duration(i))
				fmt.Printf("[ERROR] error : %v ouccred when update id: %s. Retry times: %d  \n", err, data.(ForeignKeyIdObj).OriginalID, i)
			}
		}
	}
}

func (f ForeignKeyWorker) updateID(obj ForeignKeyIdObj) error {
	session := mgoImport.GetDb()
	defer session.Close()

	for _, foreign_col := range f.collections {

		where := bson.M{
			f.foreignColumn: obj.OriginalID,
		}

		set := bson.M{
			"$set": bson.M{
				f.foreignColumn: obj.ID,
			},
		}

		_, err := session.DB(mgoImport.G_DBname).C(foreign_col).UpdateAll(where, set)
		if err != nil {
			fmt.Printf("[ERROR] err:%v, can't update collection: %s, originalId: %s", foreign_col, obj.OriginalID)
			return err
		}
	}

	return nil
}
