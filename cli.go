package mgoImport

import (
	"fmt"
	"gopkg.in/mgo.v2"
)

var G_session *mgo.Session
var G_DBname string

func InitMgoCli(url string, dbName string) error {

	//mongdbUrl = fmt.Sprintf("mongodb://%s:%s@%s", username, password, Url)
	if cli, err := mgo.Dial(url); err != nil {
		return err
	} else {
		G_session = cli
		G_DBname = dbName
	}
	return nil
}

func Close() {
	G_session.Close()
}

func GetDb() *mgo.Session {
	return G_session.Copy()
}

func TestDb(s *mgo.Session) bool {
	if err := s.Ping(); err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

var G_ItemCollectionMap = map[string]string{
	"APP": "udp_item_app",
	"IAP": "udp_item_iap",
}
