package dbmanager

import (
	"fmt"
	"lbaas/basic/config"
	"sync"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type DBManager struct {
	ormdb *gorm.DB
}

var dbmanager *DBManager = nil
var dbmgronce sync.Once
var dbmclose sync.Once

func GetDBMgrInst() *DBManager {
	dbmgronce.Do(func() {
		dbname := config.GetCommonVipper().GetString("dbconfig.dbname")
		dbuser := config.GetCommonVipper().GetString("dbconfig.dbuser")
		dbpswd := config.GetCommonVipper().GetString("dbconfig.dbpswd")
		dbport := config.GetCommonVipper().GetInt("dbconfig.dbport")

		title := fmt.Sprintf("%s:%s@tcp(127.0.0.1:%d)/%s?charset=utf8&parseTime=true&loc=Local", dbuser,
			dbpswd, dbport, dbname)
		fmt.Println(title)
		db, err := gorm.Open("mysql", title)
		if err != nil {
			fmt.Println("dbmgr init failed, mysql open failed")
			fmt.Println("err is : ", err.Error())
			return
		}

		db.DB().SetMaxIdleConns(10)
		db.DB().SetMaxOpenConns(100)
		db.SingularTable(true)
		dbmanager = &DBManager{ormdb: db}
	})
	return dbmanager
}

func (dbm *DBManager) ReleaseDB() {
	dbmclose.Do(
		func() {
			defer dbm.ormdb.Close()
		})
}

func (dbm *DBManager) GetOrmDB() *gorm.DB {
	return dbm.ormdb
}
