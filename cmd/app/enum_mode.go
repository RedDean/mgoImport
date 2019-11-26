package main

import (
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"mgoImport"
	"time"
)

func enumMode(config *mgoImport.ConfigFile) {
	fmt.Println("[INFO] mode: replace enum mode")
	startTime := time.Now()

	enums := config.GetEnumNodeArray()
	if len(enums) == 0 {
		panic(fmt.Errorf("enums slice is empty"))
	}

	for key := range enums {
		if err := updateEnum(enums[key]); err != nil {
			fmt.Printf("[ERROR] catch an error: %v while update enum, column:%s \n", err, enums[key].EnumColumn)
			continue
		}
	}

	fmt.Println("[INFO] Elapsed time :", time.Since(startTime))
}

func updateEnum(node mgoImport.EnumNode) (err error) {
	session := mgoImport.GetDb()
	defer session.Close()

	selector := bson.M{
		node.EnumColumn: node.OldValue,
	}

	updated := bson.M{
		"$set": bson.M{
			node.EnumColumn: node.NewValue,
		},
	}

	info, err := session.DB(mgoImport.G_DBname).C(node.CollectionName).UpdateAll(selector, updated)
	if err != nil {
		return
	}

	fmt.Printf("[INFO]  matched: %d, updated %d \n", info.Matched, info.Updated)
	return
}
