package main

import (
	"flag"
	"github.com/PandoCloud/pando-cloud/pkg/mysql"
	"github.com/jinzhu/gorm"
)

const (
	flagDBHost = "dbhost"
	flagDBPort = "dbport"
	flagDBName = "dbname"
	flagDBUser = "dbuser"
	flagDBPass = "dbpass"

	defaultDBHost = "localhost"
	defaultDBPort = "3306"
	defaultDBName = "PandoCloud"
	defaultDBUser = "root"
)

var (
	confDBHost = flag.String(flagDBHost, defaultDBHost, "database host address.")
	confDBPort = flag.String(flagDBPort, defaultDBPort, "database host port.")
	confDBName = flag.String(flagDBName, defaultDBName, "database name.")
	confDBUser = flag.String(flagDBUser, defaultDBUser, "database user.")
	confDBPass = flag.String(flagDBPass, "", "databse password.")
)

var DB *gorm.DB

func getDB() (*gorm.DB, error) {
	db, err := mysql.GetClient(*confDBHost, *confDBPort, *confDBName, *confDBUser, *confDBPass)
	if err != nil {
		return nil, err
	}
	gormdb, err := gorm.Open("mysql", db)
	if err != nil {
		return nil, err
	}
	gormdb.SingularTable(true)
	gormdb.LogMode(true)
	return &gormdb, nil
}
