package main

import (
	"fmt"
	"mgoImport"
	"mgoImport/id"
	"sync"
	"time"
)

// replace id mode
const (
	SELF_REPLACE = iota + 1
	MULTI_REPLACE
)

func changeIDMode(conf *mgoImport.ConfigFile) {
	var (
		data   [][]interface{}
		wg     sync.WaitGroup
		taskCh chan []interface{}
		worker id.Worker
	)

	startTime := time.Now()
	// 1. get distinct id
	fmt.Println("[INFO] start to get ids.")
	data = getLoadData(conf)

	fmt.Printf("[INFO] number of chunks: %d \n", len(data))

	// 2. get worker
	fmt.Println("[INFO] get worker.")
	worker = getWorker(conf)

	// 3. start up producers, distribute task
	wg.Add(len(data))
	taskCh = make(chan []interface{})

	for i := 0; i < len(data); i++ {
		go producer(&wg, taskCh, worker)
	}

	for _, value := range data {
		taskCh <- value
	}

	close(taskCh)
	wg.Wait()

	fmt.Println("[INFO] job done!")
	fmt.Println("[INFO] Elapsed time :", time.Since(startTime))
}

func getLoadData(conf *mgoImport.ConfigFile) [][]interface{} {
	var loader id.Loader

	idConf := conf.GetIDConf()

	if *IDSelfCollection {
		loader = id.NewIDLoader(idConf.Collections[0], idConf.RelatedColumn)
	} else if *IDChPacking {
		loader = id.NewChannelForeignKeyLoader(idConf.Collections[0], idConf.RelatedColumn, idConf.IdColumn)
	} else {
		loader = id.NewForeignKeyLoader(idConf.Collections[0], idConf.RelatedColumn, idConf.IdColumn)
	}

	loader.Load()
	return id.DivideIntoSmallChunks(loader.GetData())
}

func getWorker(conf *mgoImport.ConfigFile) id.Worker {
	var worker id.Worker

	idConf := conf.GetIDConf()
	if *IDSelfCollection {
		worker = id.NewIDWorker(idConf.Collections[0], idConf.RelatedColumn)
	} else if *IDChPacking {
		worker = id.NewChannelPackingIdWorker(idConf.Collections[1], idConf.ForeignColumn)
	} else {
		worker = id.NewForeignKeyWorker(idConf.Collections[1:], idConf.ForeignColumn)
	}

	return worker
}

func producer(wg *sync.WaitGroup, taskCh <-chan []interface{}, worker id.Worker) {
	defer wg.Done()
	for {
		task, ok := <-taskCh
		if !ok {
			return
		}

		swg := sync.WaitGroup{}
		swg.Add(G_WORKERNUM)

		workerCh := make(chan interface{})
		for i := 0; i < G_WORKERNUM; i++ {
			go worker.Do(workerCh, &swg)
		}

		for _, value := range task {
			workerCh <- value
		}

		close(workerCh)
		swg.Wait()
	}
}
