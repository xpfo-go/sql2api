package db

import (
	"fmt"
	"github.com/xpfo-go/sql2api/database"
	"github.com/xpfo-go/sql2api/util"
	"net/http"
)

type CreateDBConnectReq struct {
	DatabaseType database.DBType `json:"database_type"`
	UniqueDBName string          `json:"unique_db_name"`

	User     string `json:"user"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Database string `json:"database"`

	// have default value
	MaxOpenConn           int `json:"maxOpenConn"`
	MaxIdleConn           int `json:"maxIdleConn"`
	ConnMaxLifetimeSecond int `json:"connMaxLifetimeSecond"`
}

func CreateDBConnect(w http.ResponseWriter, r *http.Request) {
	var params CreateDBConnectReq
	if err := util.BindJson(r, &params); err != nil {
		util.ResponseJson(&w, http.StatusBadRequest, []byte(err.Error()))
		return
	}

	switch params.DatabaseType {
	case database.TypeOfMysql:
		if _, ok := database.MysqlManage.IsExist(params.UniqueDBName); ok {
			util.ResponseJson(&w, http.StatusBadRequest, []byte(fmt.Sprintf("%s is exist.", params.UniqueDBName)))
			return
		}

		client := database.NewMysqlClient(&database.MysqlConfig{
			User:                  params.User,
			Password:              params.Password,
			Host:                  params.Host,
			Port:                  params.Port,
			Database:              params.Database,
			MaxOpenConn:           params.MaxOpenConn,
			MaxIdleConn:           params.MaxIdleConn,
			ConnMaxLifetimeSecond: params.ConnMaxLifetimeSecond,
		})

		if err := client.Connect(); err != nil {
			util.ResponseJson(&w, http.StatusInternalServerError, []byte(err.Error()))
			return
		}

		if err := database.MysqlManage.AddClient(params.UniqueDBName, client); err != nil {
			util.ResponseJson(&w, http.StatusInternalServerError, []byte(err.Error()))
			return
		}
		// todo more type
	default:
		util.ResponseJson(&w, http.StatusBadRequest, []byte("not have database type."))
		return
	}

	util.ResponseJson(&w, http.StatusOK, map[string]interface{}{
		"message": "success",
	})
	return
}
