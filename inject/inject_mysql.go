package inject

import (
	"github.com/xpfo-go/sql2api/database"
	"github.com/xpfo-go/sql2api/util"
	"net/http"
)

func CreateHandler(db *database.MysqlClient, sql string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := db.TestConnection(); err != nil {
			util.ResponseJson(&w, http.StatusInternalServerError, "Connect db failed.")
			return
		}

		rows, err := db.DB.Query(sql)
		if err != nil {
			util.ResponseJson(&w, http.StatusInternalServerError, "Query error.")
			return
		}

		columns, err := rows.Columns()
		if err != nil {
			util.ResponseJson(&w, http.StatusInternalServerError, "Failed to get columns.")
			return
		}

		values := make([]interface{}, len(columns))
		valuesPtr := make([]interface{}, len(columns))
		for i := range columns {
			valuesPtr[i] = &values[i]
		}

		data := make([]map[string]interface{}, 0)
		for rows.Next() {

			if err := rows.Scan(valuesPtr...); err != nil {
				util.ResponseJson(&w, http.StatusInternalServerError, "Failed to scan row.")
				return
			}

			rowMap := make(map[string]interface{})
			for i, col := range columns {
				var v interface{}

				// 判断是否是byte
				if b, ok := values[i].([]byte); ok {
					v = string(b)
				} else {
					v = values[i]
				}

				rowMap[col] = v
			}
			data = append(data, rowMap)
		}

		util.ResponseJson(&w, http.StatusOK, data)
	}
}
