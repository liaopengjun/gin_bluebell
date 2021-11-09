package mysql

import (
	"fmt"
	"gin_bluebell/settings"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

var db *sqlx.DB
func Init(cfg *settings.MySQLConfig)(err error){
	//配置数据库连接参数
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=Local",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DB,
	)
	//连接数据库
	db,err = sqlx.Connect("mysql",dsn)
	if err != nil{
		zap.L().Error("connect DB failed err ",zap.Error(err))
	}
	//设置到数据库的最大打开连接数
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	//设置空闲的最大连接数
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	return
}
// Close 关闭MySQL连接
func Close() {
	_ = db.Close()
}