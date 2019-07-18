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

func NewMgr(p *DataParser,repository *Repository, workerSize int) *Mgr {
	return &Mgr{
		parser:p,
		repo:repository,
		workerSize:workerSize,
	}
}

func (m *Mgr) Run()  {

	wg := new(sync.WaitGroup)

	go func() {
		if err := m.parser.readLine(); err != nil {
			panic(err)
		}
	}()

	wg.Add(m.workerSize)
	for i := 0; i < m.workerSize; i++ {
		go m.process(wg)
	}
	wg.Wait()

}

func (m *Mgr) process(wg *sync.WaitGroup, deli string)  {
	defer wg.Done()
	for value := range m.parser.DataCh {
		if strArr,err := splitByDelimiter(value,deli);err != nil {
			fmt.Printf("split err: %v", err)
			continue
		}else {
			m.importData(strArr)
		}
	}
}

func (m *Mgr) importData(dataArr []string)  {
	if model, err := m.repo.BuildModel(dataArr); err != nil {
		fmt.Printf("build model err %v \n", err)
	}else {
		// todo
		model = model
	}
}