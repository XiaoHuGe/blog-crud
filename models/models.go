package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/labstack/gommon/log"
	"time"
	"xhblog/utils/setting"
)

var (
	db *gorm.DB
)

type Model struct {
	ID         int `gorm:"primary_key" json:"id"`
	CreatedOn  int `json:"created_on"`
	ModifiedOn int `json:"modified_on"`
}

func Setup() {
	// 加载数据库配置
	//sec, err := setting.Cfg.GetSection("database")
	//if err != nil {
	//	log.Fatalf("配置文件获取'database'失败: %v", err)
	//}
	var err error
	dbType := setting.DatabaseSetting.Type
	user := setting.DatabaseSetting.User
	password := setting.DatabaseSetting.Password
	host := setting.DatabaseSetting.Host
	dbName := setting.DatabaseSetting.Name
	tablePrefix := setting.DatabaseSetting.TablePrefix

	// 连接数据库
	db, err = gorm.Open(dbType, fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local", user, password, host, dbName))

	//dbInfo := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local", user, password, host, dbName)
	//db, err = gorm.Open(dbType, dbInfo)
	if err != nil {
		log.Printf("gorm.Open fail: %v", err)
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
	db.AutoMigrate(&Article{})
	db.AutoMigrate(&Auth{})

	db.Callback().Create().Replace("gorm:update_time_stamp", updateTimeStampForCreateCallback)
	db.Callback().Update().Replace("gorm:update_time_stamp", updateTimeStampForUpdateCallback)
}

// 【？何时需要关闭】
func CloseDB() {
	defer db.Close()
}

func updateTimeStampForCreateCallback(scope *gorm.Scope) {
	if !scope.HasError() {
		nowTime := time.Now().Unix()
		if createTimeField, ok := scope.FieldByName("CreatedOn"); ok {
			if createTimeField.IsBlank {
				createTimeField.Set(nowTime)
			}
		}

		if modifyTimeField, ok := scope.FieldByName("ModifiedOn"); ok {
			if modifyTimeField.IsBlank {
				modifyTimeField.Set(nowTime)
			}
		}
	}
}

func updateTimeStampForUpdateCallback(scope *gorm.Scope) {
	if _, ok := scope.Get("gorm:update_column"); ok {
		scope.SetColumn("ModifiedOn", time.Now().Unix())
	}
}
