package persistence

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/xpfo-go/sql2api/persistence/model"
)

type IRouterManage interface {
	GetRouterList(pn, pageSize int) (out []*model.Router, err error)
	CreateRouter(in *model.Router) error
	DeleteRouter(DBName string) error
}

func NewRouterManage() IRouterManage {
	return &routerManage{
		fields:    "id,method,router,db_name,sql_str",
		tableName: routerTableName,
		client:    sqliteClient,
	}
}

type routerManage struct {
	fields    string
	tableName string
	client    *sqlx.DB
}

func (r *routerManage) GetRouterList(pn, pageSize int) (out []*model.Router, err error) {
	querySql := fmt.Sprintf("select %s from %s", r.fields, r.tableName)
	if pn != 0 && pageSize != 0 {
		querySql = fmt.Sprintf("%s limit %d,%d", querySql, (pn-1)*pageSize, pageSize)
	}
	rows, err := r.client.Query(querySql)
	if err != nil {
		return nil, err
	}
	return r.scanRouterRows(rows)
}

func (r *routerManage) CreateRouter(in *model.Router) error {
	if in == nil {
		return nil
	}
	insertSql := fmt.Sprintf("insert into %s (method,router,db_name,sql_str) values (?,?,?,?)", r.tableName)
	_, err := r.client.Exec(insertSql, in.Method, in.Router, in.DBName, in.SqlStr)
	return err
}

func (r *routerManage) DeleteRouter(router string) error {
	delSql := fmt.Sprintf("delete from %s where router = ?", r.tableName)
	_, err := r.client.Exec(delSql, router)
	return err
}

func (r *routerManage) scanRouterRows(rows *sql.Rows) ([]*model.Router, error) {
	res := make([]*model.Router, 0)
	for rows.Next() {
		row := model.Router{}
		if err := rows.Scan(&row.Id, &row.Method, &row.Router, &row.DBName, &row.SqlStr); err != nil {
			return nil, err
		}
		res = append(res, &row)
	}
	return res, nil
}
