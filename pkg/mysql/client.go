package mysql

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

var mapClients map[string]*sql.DB

func GetClient(dbhost, dbport, dbname, dbuser, dbpass string) (*sql.DB, error) {

	pattern := dbuser + ":" + dbpass + "@tcp(" + dbhost + ":" + dbport + ")/" + dbname
	_, exist := mapClients[pattern]
	if !exist {
		var err error
		mapClients[pattern], err = sql.Open("mysql", pattern+"?charset=utf8&parseTime=True")
		if err != nil {
			return nil, err
		}
		err = mapClients[pattern].Ping()
		if err != nil {
			return nil, err
		}
	}

	return mapClients[pattern], nil
}

func init() {
	mapClients = make(map[string]*sql.DB)

	timer := time.NewTicker(30 * time.Second)
	go func() {
		for {
			<-timer.C
			for _, db := range mapClients {
				db.Ping()
			}
		}
	}()
}
