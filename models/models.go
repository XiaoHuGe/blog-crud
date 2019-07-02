package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/labstack/gommon/log"
	"xh-blog/utils/setting"
)

var (
	db *gorm.DB
)

type Model struct {
	ID int `gorm:"primary_key" json:"id"`
	CreatedOn int `json:"created_on"`
	ModifiedOn int `json:"modified_on"`
}

func init()  {
	// 加载数据库配置
	sec, err := setting.Cfg.GetSection("database")
	if err != nil {
		log.Fatalf("配置文件获取'database'失败: %v", err)
	}
	dbType := sec.Key("TYPE").String()
	user := sec.Key("USER").String()
	password := sec.Key("PASSWORD").String()
	host := sec.Key("HOST").String()
	dbName := sec.Key("NAME").String()
	tablePrefix := sec.Key("TABLE_PREFIX").String()

	// 连接数据库
	db, err = gorm.Open(dbType, fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local", user, password, host, dbName))

	//dbInfo := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local", user, password, host, dbName)
	//db, err = gorm.Open(dbType, dbInfo)
	if err != nil {
		log.Fatal("gorm.Open fail:", err)
	}

	// 修改默认表名
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return tablePrefix + defaultTableName
	}

	// 全局禁用表名复数 [不太理解]
	db.SingularTable(true)

	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)

	db.AutoMigrate(&Tag{})
}

// 【？何时需要关闭】
func CloseDB()  {
	defer db.Close()
}