package model

import "github.com/xpfo-go/sql2api/database"

type DBConn struct {
	Id         int             `orm:"id"`
	DBName     string          `orm:"db_name"`
	DBType     database.DBType `orm:"db_type"` // 1: mysql 2: pgsql 3: clickhouse 4: redis
	ConfigJson string          `orm:"config_json"`
}
