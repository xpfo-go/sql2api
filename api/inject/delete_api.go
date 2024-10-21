package inject

import (
	"github.com/xpfo-go/logs"
	"github.com/xpfo-go/sql2api/persistence"
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
		util.ResponseJson(&w, http.StatusBadRequest, []byte(err.Error()))
		return
	}

	if err := persistence.NewRouterManage().DeleteRouter(params.Url); err != nil {
		logs.Error(err.Error())
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
