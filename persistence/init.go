package persistence

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/xpfo-go/sql2api/util"
)

const (
	SqlitePath      = "data/RUNTIME.DAT"
	dbConnTableName = "db_conn"
	routerTableName = "router"
)

// dbConnSql
// db_type 1: mysql 2: pgsql 3: clickhouse 4: redis
var (
	dbConnSql = fmt.Sprintf(`
        CREATE TABLE IF NOT EXISTS %s (
                id INTEGER PRIMARY KEY,
                db_name TEXT NOT NULL UNIQUE,
                db_type INTEGER NOT NULL DEFAULT 0,
                config_json TEXT NOT NULL
        );
`, dbConnTableName)

	routerSql = fmt.Sprintf(`
        CREATE TABLE IF NOT EXISTS %s (
                id INTEGER PRIMARY KEY,
                method TEXT NOT NULL,
                router TEXT NOT NULL UNIQUE,
                db_name TEXT NOT NULL,
                sql_str TEXT NOT NULL
        );
`, routerTableName)
)

var (
	sqliteClient *sqlx.DB
)

func InitSqlite(ctx context.Context) {
	if err := util.CreateFileIfNotExist(SqlitePath); err != nil {
		panic(err.Error())
	}

	if err := connectSqlite(); err != nil {
		panic(err.Error())
	}

	if err := initTable(); err != nil {
		panic(err.Error())
	}

	go func() {
		<-ctx.Done()
		closeSqlite()
	}()
}

func closeSqlite() {
	if sqliteClient != nil {
		_ = sqliteClient.Close()
	}
}

func connectSqlite() (err error) {
	sqliteClient, err = sqlx.Open("sqlite3", SqlitePath)
	if err != nil {
		return
	}
	err = sqliteClient.Ping()
	return
}

func initTable() error {
	if _, err := sqliteClient.Exec(dbConnSql); err != nil {
		return err
	}
	if _, err := sqliteClient.Exec(routerSql); err != nil {
		return err
	}
	return nil
}
