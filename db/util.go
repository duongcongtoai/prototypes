package db

import (
	"errors"
	"path"
	"runtime"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/spf13/viper"
)

var (
	username string
	password string
	dbName   string
	port     string
	host     string
	initErr  error
)

func init() {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		initErr = errors.New("No caller information")
		return
	}

	viper.AddConfigPath(path.Join(path.Dir(filename)) + "/config")

	viper.SetConfigName("database")

	err := viper.ReadInConfig()
	if err != nil {
		initErr = err
		return
	}
	username = viper.GetString("username")
	password = viper.GetString("password")
	dbName = viper.GetString("db_name")
	host = viper.GetString("host")
	port = viper.GetString("port")
}

// OpenDb used to open new connection, used only in main function, implementation needs to close the connection
func OpenDb() (*gorm.DB, func(), error) {
	if initErr != nil {
		return nil, func() {}, initErr
	}

	connString := username + ":" + password +
		"@tcp(" + host + ":" + port + ")/" +
		dbName + "?charset=utf8&parseTime=True&loc=Local"
	db, err := gorm.Open("mysql", connString)
	if err != nil {
		return nil, func() {}, err
	}
	return db, func() {
		db.Close()
	}, nil
}
