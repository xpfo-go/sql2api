package api

import (
	"github.com/xpfo-go/sql2api/api/db"
	"github.com/xpfo-go/sql2api/api/inject"
	"github.com/xpfo-go/sql2api/server"
	"net/http"
)

func RegisterRouter() {
	// db router register
	{
		server.GetRouter().RegisterFunc(http.MethodGet, "/db/list", db.ListDB)
		server.GetRouter().RegisterFunc(http.MethodPost, "/db/create", db.CreateDBConnect)
	}

	// api router register
	{
		server.GetRouter().RegisterFunc(http.MethodGet, "/api/list", inject.ListApi)
		server.GetRouter().RegisterFunc(http.MethodPost, "/api/create", inject.CreateApi)
		server.GetRouter().RegisterFunc(http.MethodPost, "/api/delete", inject.DeleteApi)
	}

}
