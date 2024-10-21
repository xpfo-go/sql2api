package inject

import (
	"fmt"
	"github.com/xpfo-go/logs"
	"github.com/xpfo-go/sql2api/database"
	"github.com/xpfo-go/sql2api/inject"
	"github.com/xpfo-go/sql2api/persistence"
	"github.com/xpfo-go/sql2api/persistence/model"
	"github.com/xpfo-go/sql2api/server"
	"github.com/xpfo-go/sql2api/util"
	"net/http"
)

type CreateApiReq struct {
	UniqueDBName string `json:"unique_db_name"`
	Method       string `json:"method"`
	Url          string `json:"url"`
	Sql          string `json:"sql"`
}

func CreateApi(w http.ResponseWriter, r *http.Request) {
	var params CreateApiReq
	if err := util.BindJson(r, &params); err != nil {
		util.ResponseJson(&w, http.StatusBadRequest, []byte(err.Error()))
		return
	}

	dbClient, ok := database.MysqlManage.IsExist(params.UniqueDBName)
	if !ok {
		util.ResponseJson(&w, http.StatusBadRequest, []byte(fmt.Sprintf("%s not exist.", params.UniqueDBName)))
		return
	}

	if params.Method == "" {
		params.Method = http.MethodGet
	}

	if params.Url == "" {
		params.Url = fmt.Sprintf("/%s", util.GetRandomStr())
	}

	// TODO: 验证sql

	if err := persistence.NewRouterManage().CreateRouter(&model.Router{
		Method: params.Method,
		Router: params.Url,
		DBName: params.UniqueDBName,
		SqlStr: params.Sql,
	}); err != nil {
		logs.Error(err.Error())
		util.ResponseJson(&w, http.StatusInternalServerError, []byte(err.Error()))
		return
	}
	// 注册路由
	server.GetRouter().RegisterFunc(params.Method, params.Url, inject.CreateHandler(dbClient, params.Sql))

	util.ResponseJson(&w, http.StatusOK, map[string]interface{}{
		"message": "success",
	})
	return
}
