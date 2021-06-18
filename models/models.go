package models

import (
	"GoGin/config"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var db *gorm.DB

type Model struct {
	ID         int `gorm:"primary_key" json:"id"`
	CreatedOn  int `gorm:"created_on" json:"created_on"`
	ModifiedOn int `gorm:"modified_on" json:"modified_on"`
}

func init() {
	var (
		err                                               error
		dbType, dbName, user, password, host, tablePrefix string
	)
	//sec, err := config.Cfg.Sections("database")
	for _, v := range config.Cfg.Sections() {
		fmt.Println(v.KeyStrings())
	}
	dbType = config.Cfg.Section("database").Key("TYPE").String()
	dbName = config.Cfg.Section("database").Key("NAME").String()
	user = config.Cfg.Section("database").Key("USER").String()
	password = config.Cfg.Section("database").Key("PASSWORD").Value()
	host = config.Cfg.Section("database").Key("HOST").String()
	tablePrefix = config.Cfg.Section("database").Key("TABLE_PREFIX").String()

	//
	//if err != nil {
	//	log.Fatal(2, "Fail to get section 'database' : %v", err)
	//}
	//dbType = sec.Key("TYPE").String()
	//dbName = sec.Key("NAME").String()
	//user = sec.Key("USER").String()
	//password = sec.Key("PASSWORD").String()
	//host = sec.Key("HOST").String()
	//tablePrefix = sec.Key("TABLE_PREFIX").String()

	//content_url := fmt.Sprintf("%%:%%@tcp(%%)/%%?charset=utf8&parseTime=True&loc=Local",
	//	user, password, host, dbName)
	contentUrl := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local", user, password, host, dbName)
	db, err = gorm.Open(dbType, contentUrl)
	if err != nil {
		//log.Println(err)
		panic(err)
	}
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return tablePrefix + defaultTableName
	}
	db.SingularTable(true)
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)

}
func CloseDB() {
	defer db.Close()
}
