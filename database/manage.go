package database

type IDatabaseClient interface {
	TestConnection() (err error)
	Connect() error
	Close()
}

type DBType int

const (
	TypeOfMysql DBType = iota + 1
	TypeOfPgsql
	TypeOfClickhouse
	TypeOfRedis
)

type IDatabaseManage struct {
}

func init() {

}
