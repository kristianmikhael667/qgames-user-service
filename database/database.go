package database

import (
	"main/package/util"
	"sync"

	"gorm.io/gorm"
)

var (
	dbConn *gorm.DB
	once   sync.Once
)

func CreateConnection() {
	conf := dbConfig{
		Host: util.Getenv("DB_HOST", "localhost"),
		User: util.Getenv("DB_USER", "root"),
		Pass: util.Getenv("DB_PASS", ""),
		Port: util.Getenv("DB_PORT", "3306"),
		Name: util.Getenv("DB_NAME", "qgrowid"),
	}
	mysql := mysqlConfig{dbConfig: conf}
	once.Do(func() {
		mysql.Connect()
	})
}

func GetConnection() *gorm.DB {
	if dbConn == nil {
		CreateConnection()
	}
	return dbConn
}
