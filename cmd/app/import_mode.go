package main

import (
	"fmt"
	"mgoImport"
)

func importMode(c *mgoImport.ConfigFile) {
	fmt.Println("[INFO] mode: import data mode")

	p, err := mgoImport.InitParser(*fileDir, *limitation, c.Delimiter, *readerSize)
	if err != nil {
		panic(err)
	}
	//fmt.Println(c)

	r := mgoImport.InitRepository(c, getImportMode())

	mgr := mgoImport.NewMgr(p, r, *size)

	mgr.Run(getImportMode())
}

func getImportMode() int {
	if *isModifyFieldsModel {
		return mgoImport.MODIFY
	} else if *isItem {
		return mgoImport.ITEM
	} else if *isItemHis {
		return mgoImport.ITEM_HIS
	} else {
		return mgoImport.NORMAL
	}
}
