package mgoImport

import (
	"fmt"
	"sync"
)

type Mgr struct {
	parser *DataParser
	repo   *Repository

	workerSize int
}

func NewMgr(p *DataParser,repository *Repository, workerSize int){

}

func (m *Mgr) Run()  {

	wg := new(sync.WaitGroup)
	ch := make(chan map[string]interface{},  10)

	go func() {
		if err := m.parser.readLine(); err != nil {
			panic(err)
		}
	}()

	wg.Add(m.workerSize)
	for i := 0; i < m.workerSize; i++ {
		go m.process(wg,ch)
	}
	wg.Wait()
	close(ch)

	for v := range ch {
		fmt.Println(v)

		// save operation
		m.repo.Db = append(m.repo.Db, v)
	}

}

func (m *Mgr) process(wg *sync.WaitGroup, dataCh chan map[string]interface{})  {
	defer wg.Done()
	for value := range m.parser.DataCh {
		if strArr,err := splitByDelimiter(value,"|");err != nil {
			fmt.Printf("split err: %v", err)
			continue
		}else {
			m.importData(strArr,dataCh)
		}
	}
}

func (m *Mgr) importData(dataArr []string, ch chan map[string]interface{})  {
	if model, err := m.repo.BuildModel(dataArr); err != nil {
		fmt.Printf("build model err %v \n", err)
	}else {
		ch <- model
	}
}