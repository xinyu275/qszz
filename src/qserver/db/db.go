package db

//mysql
import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/liangdas/mqant/conf"
)

var DB *sql.DB

func InitDB(config conf.Config) {
	if dbDriver, ok := config.Settings["db"]; ok {
		dbDriver := dbDriver.(string)
		if dbDriver == "mysql" {
			InitMysqlDB(config)
		}
	}
}

func InitMysqlDB(config conf.Config) {
	host := config.Settings["host"].(string)
	port := int(config.Settings["port"].(float64))
	user := config.Settings["user"].(string)
	pass := config.Settings["pass"].(string)
	dbName := config.Settings["dbName"].(string)
	charset := config.Settings["charset"].(string)
	maxOpenConns := int(config.Settings["maxOpenConns"].(float64))
	maxIdleConns := int(config.Settings["maxIdleConns"].(float64))

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s",
		user, pass, host, port, dbName, charset)
	var err error
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		panic("InitMysqlDB error" + err.Error())
	}
	DB.SetMaxOpenConns(maxOpenConns)
	DB.SetMaxIdleConns(maxIdleConns)
}
