package mysql

import (
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql" //匿名导入 默认会自动执行该包中的init()方法
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var db *sqlx.DB

//Init 初始化方法
func Init() (err error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True",
		viper.GetString("mysql.user"),
		viper.GetString("mysql.passwd"),
		viper.GetString("mysql.host"),
		viper.GetInt("mysql.port"),
		viper.GetString("mysql.dbname"),
	)

	db, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		zap.L().Error("connect DB failed", zap.Error(err))
		return
	}
	db.SetMaxOpenConns(viper.GetInt("mysql.max_open_conns"))
	db.SetMaxIdleConns(viper.GetInt("mysql.max_idle_conns"))
	db.SetConnMaxLifetime(time.Duration(viper.GetInt("mysql.max_left_time")) * time.Second)
	return
}

//Close 关闭MySQL实例
func Close() {
	db.Close()
}
