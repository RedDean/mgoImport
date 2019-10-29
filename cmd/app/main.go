package main

import (
	"fmt"
	"gopkg.in/alecthomas/kingpin.v2"
	"mgoImport"
)

/*
we're gonna build a util that is used for import data from csv file to mongoDB.
In this case, I will attempt to build this app by using a classic programming thought
called TDD for the first time.

we're gonna complete these features:

TODO:

    1. accept command line params.

    2. read configurable file done

    3. set working pool limitation while importing data.

    4. import different model flexibly by using configurable file
*/

var (
	configDir = kingpin.Flag("config", "config file directory").Required().String()
	fileDir   = kingpin.Flag("file", "csv file directory").String()

	// may not need this variable.
	isModifyFieldsModel = kingpin.Flag("modify", "modify some filed in mongo").Bool()

	isItem    = kingpin.Flag("item", "import item data specially").Bool()
	isItemHis = kingpin.Flag("itemHis", "import item history data specially").Bool()

	IDMode           = kingpin.Flag("id", "change program mode into id mode.").Bool()
	IDSelfCollection = kingpin.Flag("id-self", " use new ObjectID replace filed id in one collection").Bool()

	limitation = kingpin.Flag("limit", "channel size limitation while parsing file").Default("30").Int()
	size       = kingpin.Flag("size", "number of processing data goroutines").Default("3").Int()
	readerSize = kingpin.Flag("readerSize", "reader buffer size").Default("4096").Int()
)

var Run func(file *mgoImport.ConfigFile)

func init() {

	kingpin.Version("0.0.1")
	kingpin.Parse()

	if *IDMode {
		fmt.Println("[INFO] mode: idMode")
		Run = changeIDMode
	} else {
		fmt.Println("[INFO] mode: importMode")
		Run = importMode
	}
}

func main() {

	c := mgoImport.InitConfig(*configDir)

	if err := mgoImport.InitMgoCli(c.Db.Url, c.Db.Name); err != nil {
		panic(err)
	}
	defer mgoImport.Close()

	Run(c)
}
