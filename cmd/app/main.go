package main

import (
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

func main() {

	// config.json
	config := mgoImport.InitConfig("")


	parser := mgoImport.InitParser("",0)



	mgr := mgoImport.Mgr{
    		parser:parser,

	}

	mgr.Run()
}
