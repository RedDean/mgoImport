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
	configDir = kingpin.Flag("config", "config file directory").Required().String()
	fileDir   = kingpin.Flag("file", "csv file directory").Required().String()

	limitation = kingpin.Flag("limit", "channel size limitation while parsing file").Default("500").Int()
	size       = kingpin.Flag("size", "number of processing data goroutines").Default("5").Int()
)

func main() {

	kingpin.Version("0.0.1")
	kingpin.Parse()

	c := mgoImport.InitConfig(*configDir)

	p, err := mgoImport.InitParser(*fileDir, *limitation, c.Delimiter)
	if err != nil {
		panic(err)
	}
	//fmt.Println(c)

	r := mgoImport.InitRepository(c)

	if err := mgoImport.InitMgoCli(c.Db.Url); err != nil {
		panic(err)
	}
	defer mgoImport.Close()

	mgr := mgoImport.NewMgr(p, r, *size)

	mgr.Run()
}
