package mgoImport

import (
	"fmt"
	"gopkg.in/mgo.v2"
)

var G_session *mgo.Session
var G_DBname string

func InitMgoCli(db DbConfig) error {

	//mongdbUrl = fmt.Sprintf("mongodb://%s:%s@%s", username, password, Url)
	if cli, err := mgo.Dial(db.Url); err != nil {
		return err
	} else {
		G_session = cli
		G_DBname = db.Name

		G_session.SetPoolLimit(50)
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

var G_ItemCollection_His_Map = map[string]string{
	"APP": "udp_item_app_history",
	"IAP": "udp_item_iap_history",
}
