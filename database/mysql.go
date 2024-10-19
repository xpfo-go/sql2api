package database

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/xpfo-go/logs"
	"net/url"
	"time"
)

const (
	defaultMaxOpenConn     = 100
	defaultMaxIdleConn     = 25
	defaultConnMaxLifetime = 10 * time.Minute
)

//var (
//	DefaultDBClient *MysqlClient
//)

// NewMysqlClient :
func NewMysqlClient(cfg *MysqlConfig) *MysqlClient {
	dataSource := fmt.Sprintf("%s:%s@(%s:%d)/%s?charset=%s&parseTime=True&interpolateParams=true&loc=%s&time_zone=%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
		"utf8",
		"UTC",
		url.QueryEscape("'+00:00'"),
	)

	maxOpenConn := defaultMaxOpenConn
	if cfg.MaxOpenConn > 0 {
		maxOpenConn = cfg.MaxOpenConn
	}

	maxIdleConn := defaultMaxIdleConn
	if cfg.MaxIdleConn > 0 {
		maxIdleConn = cfg.MaxIdleConn
	}

	if maxOpenConn < maxIdleConn {
		logs.Errorf("error config for database %s, maxOpenConn should greater or equals to maxIdleConn, will"+
			"use the default [defaultMaxOpenConn=%d, defaultMaxIdleConn=%d]",
			cfg.Database, defaultMaxOpenConn, defaultMaxIdleConn)
		maxOpenConn = defaultMaxOpenConn
		maxIdleConn = defaultMaxIdleConn
	}

	connMaxLifetime := defaultConnMaxLifetime
	if cfg.ConnMaxLifetimeSecond > 0 {
		if cfg.ConnMaxLifetimeSecond >= 60 {
			connMaxLifetime = time.Duration(cfg.ConnMaxLifetimeSecond) * time.Second
		} else {
			logs.Errorf("error config for database %s, connMaxLifetimeSeconds should be greater than 60 seconds"+
				"use the default [defaultConnMaxLifetime=%s]",
				cfg.Database, defaultConnMaxLifetime)
		}
	}

	return &MysqlClient{
		name:            cfg.Database,
		dataSource:      dataSource,
		maxOpenConn:     maxOpenConn,
		maxIdleConn:     maxIdleConn,
		connMaxLifetime: connMaxLifetime,
	}
}

// MysqlClient MySQL DB Instance
type MysqlClient struct {
	name            string
	DB              *sqlx.DB
	dataSource      string
	maxOpenConn     int
	maxIdleConn     int
	connMaxLifetime time.Duration
}

type MysqlConfig struct {
	User     string
	Password string
	Host     string
	Port     int
	Database string

	// have default value
	MaxOpenConn           int
	MaxIdleConn           int
	ConnMaxLifetimeSecond int
}

// TestConnection ...
func (db *MysqlClient) TestConnection() (err error) {
	// TODO：超时控制
	conn, err := sqlx.Connect("mysql", db.dataSource)
	if err != nil {
		return
	}

	_ = conn.Close()
	return nil
}

// Connect to db, and update some settings
func (db *MysqlClient) Connect() error {
	// TODO：超时控制
	var err error
	db.DB, err = sqlx.Connect("mysql", db.dataSource)
	if err != nil {
		return err
	}

	db.DB.SetMaxOpenConns(db.maxOpenConn)
	db.DB.SetMaxIdleConns(db.maxIdleConn)
	db.DB.SetConnMaxLifetime(db.connMaxLifetime)

	logs.Infof("connect to database: %s[maxOpenConn=%d, maxIdleConn=%d, connMaxLifetime=%s]",
		db.name, db.maxOpenConn, db.maxIdleConn, db.connMaxLifetime)

	return nil
}

func (db *MysqlClient) Query(sql string) (*sql.Rows, error) {
	return db.DB.Query(sql)
}

// Close db connection
func (db *MysqlClient) Close() {
	if db.DB != nil {
		_ = db.DB.Close()
	}
}
