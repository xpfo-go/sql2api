package model

type Router struct {
	Id     int    `orm:"id"`
	Method string `orm:"method"`
	Router string `orm:"router"`
	DBName string `orm:"db_name"`
	SqlStr string `orm:"sql_str"`
}
