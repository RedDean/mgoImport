package id

import (
	"gopkg.in/mgo.v2/bson"
	"mgoImport"
	"sync"
)

const DefaultPoolSize = 30

type updateObj struct {
	setter   bson.M
	selector bson.M
}

type UpdateBatchPool struct {
	collectionName string
	size           int
	buf            []*updateObj
	NextIndex      int
	mux            *sync.Mutex
	closed         bool
}

func NewUpdateBatchPool(size int, collection string) *UpdateBatchPool {
	return &UpdateBatchPool{
		size:           size,
		buf:            make([]*updateObj, size),
		mux:            new(sync.Mutex),
		collectionName: collection,
	}
}

func (up *UpdateBatchPool) Add(obj *updateObj) {
	if up.NextIndex < up.size {
		up.buf[up.NextIndex] = obj
		up.NextIndex++
	}

	// Pool is full.
	if up.NextIndex == up.size {
		//fmt.Printf("[DEBUG] pool is full with %d records, update a batch of operation \n", up.size)
		up.batchUpdate()
		up.empty()
	}
}

func (up *UpdateBatchPool) AddWithMux(obj *updateObj) {
	up.mux.Lock()
	defer up.mux.Unlock()
	if up.NextIndex < up.size {
		up.buf[up.NextIndex] = obj
		up.NextIndex++
	}

	// Pool is full.
	if up.NextIndex == up.size {
		//fmt.Printf("[DEBUG] pool is full with %d records, update a batch of operation \n", up.size)
		up.batchUpdate()
		up.empty()
	}
}

func (up *UpdateBatchPool) Clean() {
	if up.closed {
		return
	}
	up.mux.Lock()
	defer up.mux.Unlock()
	up.batchUpdate()
	up.closed = true
}

func (up UpdateBatchPool) batchUpdate() {
	session := mgoImport.GetDb()
	defer session.Close()

	bulk := session.DB(mgoImport.G_DBname).C(up.collectionName).Bulk()

	for key := range up.buf {
		ops := up.buf[key]
		if ops != nil {
			bulk.UpdateAll(ops.selector, ops.setter)
		}
	}

	if _, err := bulk.Run(); err != nil {
		panic(err)
	}
}

func (up *UpdateBatchPool) empty() {
	up.NextIndex = 0
	up.buf = make([]*updateObj, up.size)
}
