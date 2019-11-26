package mgoImport

import (
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"strconv"
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

func (m *Mgr) Run(modeType int) {

	wg := new(sync.WaitGroup)

	fmt.Println("[INFO] 开始处理！")
	fmt.Println(strings.Repeat("-", 20))

	go func() {
		if err := m.parser.readLine(); err != nil {
			panic(err)
		}
	}()

	wg.Add(m.workerSize)

	for i := 0; i < m.workerSize; i++ {
		go m.process(wg, modeType)
	}
	wg.Wait()

	fmt.Println("[INFO] 导数完成")
}

func (m *Mgr) process(wg *sync.WaitGroup, modeType int) {
	switch modeType {
	default:
		m.normalImport(wg)
	case NORMAL:
		m.normalImport(wg)
	case MODIFY:
		m.modifyImport(wg)
	case ITEM:
		m.itemImport(wg)
	case ITEM_HIS:
		m.itemHisImport(wg)
	}
}

func (m *Mgr) normalImport(wg *sync.WaitGroup) {
	defer wg.Done()
	for value := range m.parser.DataCh {
		if model, err := m.repo.BuildModel(value); err != nil {
			fmt.Printf("[ERROR] build model err: %v \n", err)
			continue
		} else {
			if len(model) != 0 {
				insert(model, m.repo.DbName, m.repo.Collection)
			}
		}
	}
}

func (m *Mgr) modifyImport(wg *sync.WaitGroup) {
	defer wg.Done()
	for value := range m.parser.DataCh {
		if value[1] != "" {
			continue
		}
		model := m.repo.BuildModifyModel(value[1])
		modify(value[0], model, m.repo.DbName, m.repo.Collection)
	}
}

func (m *Mgr) itemImport(wg *sync.WaitGroup) {
	defer func() {
		wg.Done()
		if r := recover(); r != nil {
			fmt.Printf("[ERROR] catch a panic in itemImport err: %v \n", r)
		}
		//debug.PrintStack()
	}()

	for value := range m.parser.DataCh {
		if model, err := m.repo.BuildItemModel(value); err != nil {
			fmt.Printf("[ERROR] build model err: %v \n", err)
			continue
		} else {
			if len(model) != 0 {

				if _, ok := model["type"]; !ok {
					fmt.Printf("[ERROR] invaild string : %s \n", value)
					continue
				}

				itemType := model["type"].(string)
				if "APP" != itemType && "IAP" != itemType {
					fmt.Printf("[ERROR] item has wrong type: %s! should be APP or IAP. ", model["type"].(string))
					continue
				}

				channels, ok := model["channels"].(map[string]interface{})
				if ok {
					model["channels"] = resetChannels(channels)
				}

				delete(model, "type")

				insert(G_item_rebuild_func_map[itemType](model),
					m.repo.DbName,
					G_ItemCollectionMap[itemType])
			}
		}
	}
}

func (m *Mgr) itemHisImport(wg *sync.WaitGroup) {
	defer func() {
		wg.Done()
		if r := recover(); r != nil {
			fmt.Printf("[ERROR] catch a panic in itemImport err: %v \n", r)
		}
	}()

	for value := range m.parser.DataCh {
		if model, err := m.repo.BuildItemModel(value); err != nil {
			fmt.Printf("[ERROR] build model err: %v \n", err)
			continue
		} else {
			if len(model) != 0 {
				modelRebuilt, itemType, e := rebuildItemModel(model)
				if e != nil {
					fmt.Printf("[ERROR] itemHisImport err: %v ", e)
					continue
				}

				insert(G_item_rebuild_func_map[itemType](modelRebuilt),
					m.repo.DbName,
					G_ItemCollection_His_Map[itemType])
			}
		}
	}
}

func insert(model map[string]interface{}, dbName, collection string) {
	session := GetDb()
	defer session.Close()
	if err := session.DB(dbName).C(collection).Insert(model); err != nil {
		fmt.Printf("[ERROR] insert err :%v ", err)
	}
}

func rebuildItemModel(model map[string]interface{}) (map[string]interface{}, string, error) {

	itemType := model["type"].(string)
	if "APP" != itemType && "IAP" != itemType {
		//fmt.Printf("[ERROR] item has wrong type: %s! should be APP or IAP. ", model["type"].(string))
		return nil, "", fmt.Errorf("item has wrong type: %s! should be APP or IAP. ", model["type"].(string))
	}

	channels, ok := model["channels"].(map[string]interface{})
	if ok {
		model["channels"] = resetChannels(channels)
	}

	delete(model, "type")

	return model, itemType, nil
}

func modify(id string, model map[string]interface{}, dbName, collection string) {
	session := GetDb()
	defer session.Close()
	val, err := strconv.Atoi(id)
	if err != nil {
		panic(err)
	}
	if err := session.DB(dbName).C(collection).Update(bson.M{"originalId": val}, model); err != nil {
		fmt.Printf("[ERROR] update err :%v ", err)
	}
}
