package db

import (
	"github.com/xpfo-go/sql2api/database"
	"github.com/xpfo-go/sql2api/util"
	"net/http"
)

func ListDB(w http.ResponseWriter, r *http.Request) {
	util.ResponseJson(&w, http.StatusOK, map[string]interface{}{
		"message": "success",
		"data":    database.MysqlManage.GetClientList(),
	})
	return
}
