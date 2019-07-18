package main

import (
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
	configDir = kingpin.Flag("config","config file directory").Default("./config.json").String()
	// set a limitation
	parseLimitation = kingpin.Flag("limit","limitation while parsing file").Default("500").Int()
)

func main() {

	kingpin.Version("0.0.1")
	kingpin.Parse()

	// config.json
	config := mgoImport.InitConfig(*configDir)
	parser := mgoImport.InitParser("",0,config.Delimiter)
	repo := &mgoImport.Repository{}

	mgr := mgoImport.NewMgr(parser,repo, 0)
	mgr.Run()
}
