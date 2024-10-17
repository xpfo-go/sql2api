package api

import (
	"fmt"
	"github.com/xpfo-go/sql2api/inject"
	"github.com/xpfo-go/sql2api/server"
	"github.com/xpfo-go/sql2api/util"
	"net/http"
)

type CreateApiReq struct {
	Method string `json:"method"`
	Url    string `json:"url"`
	Sql    string `json:"sql"`
}

func CreateApi(w http.ResponseWriter, r *http.Request) {
	var params CreateApiReq
	if err := util.BindJson(r, &params); err != nil {
		util.ResponseJson(&w, http.StatusInternalServerError, []byte(err.Error()))
		return
	}

	if params.Method == "" {
		params.Method = http.MethodGet
	}

	if params.Url == "" {
		params.Url = fmt.Sprintf("/%s", util.GetRandomStr())
	}

	// TODO: 验证sql

	// 注册路由
	server.GetRouter().RegisterFunc(params.Method, params.Url, inject.CreateHandler(params.Sql))

	util.ResponseJson(&w, http.StatusOK, map[string]interface{}{
		"message": "success",
	})
	return
}
