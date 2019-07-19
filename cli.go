package mgoImport

import (
	"fmt"
	"gopkg.in/mgo.v2"
)

var G_session *mgo.Session

func InitMgoCli(url string) error {

	//mongdbUrl = fmt.Sprintf("mongodb://%s:%s@%s", username, password, Url)
	if cli, err := mgo.Dial(url); err != nil {
		return err
	} else {
		G_session = cli
	}
	return nil
}

func Close() {
	G_session.Close()
}

func getDb() *mgo.Session {
	return G_session.Copy()
}

func TestDb(s *mgo.Session) bool {
	if err := s.Ping(); err != nil {
		fmt.Println(err)
		return false
	}
	return true
}
