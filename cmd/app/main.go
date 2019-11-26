package main

import (
	"fmt"
	"gopkg.in/alecthomas/kingpin.v2"
	"mgoImport"
	"runtime"
)

/*
we're gonna build a util that is used for import data from csv file to mongoDB.
In this case, I will attempt to build this app by using a classic programming thought
called TDD for the first time.

we're gonna complete these features:

    1. accept command line params.

    2. read configurable file done

    3. set working pool limitation while importing data.

    4. import different model flexibly by using configurable file
*/

var (
	configDir = kingpin.Flag("config", "config file directory").Required().String()
	dbDir     = kingpin.Flag("dbConf", "db config file directory").Default("./conf/db-config.json").String()
	fileDir   = kingpin.Flag("file", "csv file directory").String()

	// may not need this variable.
	isModifyFieldsModel = kingpin.Flag("modify", "modify some filed in mongo").Bool()

	isItem    = kingpin.Flag("item", "import item data specially").Bool()
	isItemHis = kingpin.Flag("itemHis", "import item history data specially").Bool()

	IDMode           = kingpin.Flag("id", "change program mode into id mode.").Bool()
	IDSelfCollection = kingpin.Flag("id-self", " use new ObjectID replace filed id in one collection").Bool()
	IDChPacking      = kingpin.Flag("ch-pack-id", "replace item.Channel[$store_name].PackingId").Bool()

	enum = kingpin.Flag("enum", "replace enum value").Bool()

	limitation = kingpin.Flag("limit", "buffer size limitation while parsing file").Default("30").Int()
	size       = kingpin.Flag("size", "number of processing data workers in import mode").Default("3").Int()
	readerSize = kingpin.Flag("readerSize", "reader buffer size").Default("4096").Int()

	idWorkers = kingpin.Flag("idWorkers", " id substitution process goroutine number").Default("1").Int()
)

var Run func(file *mgoImport.ConfigFile)

var G_WORKERNUM int

func init() {

	if runtime.NumCPU() == 1 {
		runtime.GOMAXPROCS(4)
	}

	kingpin.Version("0.0.1")
	kingpin.Parse()

	if *IDMode && *enum {
		panic(fmt.Errorf("Can't run id mode and enum mode at one time"))
	}

	if *IDMode {
		Run = changeIDMode
		fmt.Println("[debug] idWorkers", *idWorkers)
		G_WORKERNUM = *idWorkers
		return
	}

	if *enum {
		Run = enumMode
		return
	}

	Run = importMode
}

func main() {

	c := mgoImport.InitConfig(*configDir, *dbDir)

	if err := mgoImport.InitMgoCli(c.Db); err != nil {
		panic(err)
	}
	defer mgoImport.Close()

	Run(c)
}
