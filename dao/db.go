package dao

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
	"github.com/wsloong/go-crawler/global"
)

var masterDB *sql.DB

func init() {
	// 确保viper配置文件已经处理
	global.Init()

	var err error

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
		viper.GetString("storage.user"),
		viper.GetString("storage.password"),
		viper.GetString("storage.host"),
		viper.GetString("storage.port"),
		viper.GetString("storage.dbname"),
		viper.GetString("storage.charset"))

	masterDB, err = sql.Open(viper.GetString("storage.driver"), dsn)
	if err != nil {
		panic(err)
	}

	// 测试数据库链接是否ok
	if err = masterDB.Ping(); err != nil {
		panic(err)
	}

}
