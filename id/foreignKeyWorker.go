package id

import (
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"sync"
)

const RETRY_TIMES = 3

type ForeignKeyWorker struct {
	collections   []string
	foreignColumn string
}

func NewForeignKeyWorker(collections []string, foreignColumn string) *ForeignKeyWorker {
	if len(collections) == 0 {
		panic(fmt.Errorf("collection slice is empty, this is required"))
	}
	return &ForeignKeyWorker{
		collections:   collections,
		foreignColumn: foreignColumn,
	}
}

func (f ForeignKeyWorker) Do(dataCh <-chan interface{}, swg *sync.WaitGroup) {
	poolMap := make(map[string]*UpdateBatchPool, len(f.collections))
	for k := range f.collections {
		poolMap[f.collections[k]] = NewUpdateBatchPool(DefaultPoolSize, f.collections[k])
	}

	defer func() {
		for _, p := range poolMap {
			p.Clean()
		}
		swg.Done()
	}()

	for {
		data, ok := <-dataCh
		if !ok {
			return
		}

		for _, p := range poolMap {
			p.Add(f.buildUpdateOpsObj(data.(ForeignKeyIdObj)))
		}
	}
}

func (f ForeignKeyWorker) buildUpdateOpsObj(obj ForeignKeyIdObj) *updateObj {
	return &updateObj{
		selector: bson.M{
			f.foreignColumn: obj.OriginalID,
		},
		setter: bson.M{
			"$set": bson.M{
				f.foreignColumn: obj.ID,
			},
		},
	}
}
