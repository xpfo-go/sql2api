package inject

import (
	"github.com/xpfo-go/sql2api/server"
	"github.com/xpfo-go/sql2api/util"
	"net/http"
)

func ListApi(w http.ResponseWriter, r *http.Request) {
	util.ResponseJson(&w, http.StatusOK, map[string]interface{}{
		"message": "success",
		"data":    server.GetRouter().GetApiList(),
	})
	return
}
