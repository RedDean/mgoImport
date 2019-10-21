package id

import (
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"mgoImport"
	"sync"
)

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

		idObj := data.(ForeignKeyIdObj)
		if err := f.updateID(idObj); err != nil {
			fmt.Printf("[ERROR] error : %v ouccred when update id: %s", err, data.(string))
			continue
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

		//fmt.Printf("[DEBUG] foreign column :%s, originald: %s \n", foreign_col, obj.OriginalID )
		_, err := session.DB(mgoImport.G_DBname).C(foreign_col).UpdateAll(where, set)
		if err != nil {
			fmt.Printf("[ERROR] can't update collection: %s, originalId: %s", foreign_col, obj.OriginalID)
			continue
		}
	}

	return nil
}
