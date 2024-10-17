package api

import (
	"github.com/xpfo-go/sql2api/server"
	"github.com/xpfo-go/sql2api/util"
	"net/http"
)

type DeleteApiReq struct {
	Url string `json:"url"`
}

func DeleteApi(w http.ResponseWriter, r *http.Request) {
	var params DeleteApiReq
	if err := util.BindJson(r, &params); err != nil {
		util.ResponseJson(&w, http.StatusInternalServerError, []byte(err.Error()))
		return
	}

	// 注册路由
	server.GetRouter().DeleteRouter(params.Url)

	util.ResponseJson(&w, http.StatusOK, map[string]interface{}{
		"message": "success",
	})
	return
}
