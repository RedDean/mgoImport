package mgoImport

import (
	"fmt"
	"strings"
	"sync"
)

type Mgr struct {
	parser *DataParser
	repo   *Repository

	workerSize int
}

func NewMgr(p *DataParser, repository *Repository, workerSize int) *Mgr {
	return &Mgr{
		parser:     p,
		repo:       repository,
		workerSize: workerSize,
	}
}

func (m *Mgr) Run() {

	wg := new(sync.WaitGroup)

	fmt.Println("开始处理！")
	fmt.Println(strings.Repeat("-", 20))
	go func() {
		if err := m.parser.readLine(); err != nil {
			panic(err)
		}
	}()

	wg.Add(m.workerSize)
	for i := 0; i < m.workerSize; i++ {
		go m.process(wg, m.parser.deli)
	}
	wg.Wait()

	fmt.Println("导数完成")
}

func (m *Mgr) process(wg *sync.WaitGroup, deli string) {
	defer wg.Done()
	for value := range m.parser.DataCh {
		if model, err := m.repo.BuildModel(value); err != nil {
			fmt.Printf("build model err %v \n", err)
		} else {
			insert(model, m.repo.DbName, m.repo.Collection)
		}
	}
}

func insert(model map[string]interface{}, dbName, collection string) {
	session := getDb()
	defer session.Close()
	if err := session.DB(dbName).C(collection).Insert(model); err != nil {
		fmt.Printf("insert err :%v ", err)
	}
}
