package persistence

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/xpfo-go/sql2api/persistence/model"
)

type IDBConnManage interface {
	GetDBConnList(pn, pageSize int) (out []*model.DBConn, err error)
	CreateDBConn(in *model.DBConn) error
	DeleteDBConn(DBName string) error
}

func NewDBConnManage() IDBConnManage {
	return &dbConnManage{
		fields:    "id,db_name,db_type,config_json",
		tableName: dbConnTableName,
		client:    sqliteClient,
	}
}

type dbConnManage struct {
	fields    string
	tableName string
	client    *sqlx.DB
}

func (d *dbConnManage) GetDBConnList(pn, pageSize int) (out []*model.DBConn, err error) {
	querySql := fmt.Sprintf("select %s from %s", d.fields, d.tableName)
	if pn != 0 && pageSize != 0 {
		querySql = fmt.Sprintf("%s limit %d,%d", querySql, (pn-1)*pageSize, pageSize)
	}
	rows, err := d.client.Query(querySql)
	if err != nil {
		return nil, err
	}
	return d.scanDBConnRows(rows)
}

func (d *dbConnManage) CreateDBConn(in *model.DBConn) error {
	if in == nil {
		return nil
	}
	insertSql := fmt.Sprintf("insert into %s (db_name,db_type,config_json) values (?,?,?)", d.tableName)
	_, err := d.client.Exec(insertSql, in.DBName, in.DBType, in.ConfigJson)
	return err
}

func (d *dbConnManage) DeleteDBConn(DBName string) error {
	delSql := fmt.Sprintf("delete from %s where db_name = ?", d.tableName)
	_, err := d.client.Exec(delSql, DBName)
	return err
}

func (d *dbConnManage) scanDBConnRows(rows *sql.Rows) ([]*model.DBConn, error) {
	res := make([]*model.DBConn, 0)
	for rows.Next() {
		row := model.DBConn{}
		if err := rows.Scan(&row.Id, &row.DBName, &row.DBType, &row.ConfigJson); err != nil {
			return nil, err
		}
		res = append(res, &row)
	}
	return res, nil
}
