package inject

import (
	"encoding/json"
	"fmt"
	"github.com/xpfo-go/sql2api/database"
	"github.com/xpfo-go/sql2api/persistence"
	"github.com/xpfo-go/sql2api/server"
)

func ReloadDatabase() error {
	dbConnList, err := persistence.NewDBConnManage().GetDBConnList(0, 0)
	if err != nil {
		return err
	}

	for _, conn := range dbConnList {
		switch conn.DBType {
		case database.TypeOfMysql:
			var cfg database.MysqlConfig
			if err := json.Unmarshal([]byte(conn.ConfigJson), &cfg); err != nil {
				return err
			}
			client := database.NewMysqlClient(&cfg)
			if err := client.Connect(); err != nil {
				return err
			}

			if err := database.MysqlManage.AddClient(conn.DBName, client); err != nil {
				return err
			}
			// todo: more db type
		}
	}

	return nil
}

func ReloadRouter() error {
	routerList, err := persistence.NewRouterManage().GetRouterList(0, 0)
	if err != nil {
		return err
	}

	for _, router := range routerList {
		dbClient, ok := database.MysqlManage.IsExist(router.DBName)
		if !ok {
			return fmt.Errorf("%s not exist", router.DBName)
		}
		server.GetRouter().RegisterFunc(router.Method, router.Router, CreateHandler(dbClient, router.SqlStr))
	}

	return nil
}
