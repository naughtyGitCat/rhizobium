// @Datetime  : 2019-05-05 21:00
// @Author    : psyduck
// @Purpose   : 适配MySQL常用操作
// @TODO      : Pair programming
//
package adapter

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
	"rhizobium/common"
)

//func getConn(conf common.Config) (*sql.DB, error) {
//	ip := conf.MySQL.Host
//	port := conf.MySQL.Port
//	user := conf.MySQL.User
//	password := conf.MySQL.Password
//	database := conf.MySQL.DB
//	uri := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8", user, password, ip, port, database)
//	dbConn, err := sql.Open("mysql", uri)
//	return dbConn, err
//}
//
//func ExecSQLSelect(sqlString string, needResponse bool) {
//	conn, _ := getConn(common.CONF)
//	if needResponse {
//		rows, _ := conn.Query("select 1")
//		for rows.Next() {
//			colStrings, err := rows.Columns()
//			if err != nil {
//				p.Error(err)
//			}
//			p.Debug(colStrings)
//		}
//
//	}
//}

type MySQLConnectionPool struct {
	Uri                 string
	ConnectionWarehouse chan *sql.DB
	logger              *logrus.Logger
}

func (c *MySQLConnectionPool) Init() {
	c.logger.Infof("Now init connection by %s", c.Uri)
	for i := 1; i < 6; i++ {
		c.initConnection()
	}

	c.logger = common.GetLogger("mysql", logrus.Fields{"uri": c.Uri})
}

func (c *MySQLConnectionPool) initConnection() {
	dbConn, err := sql.Open("mysql", c.Uri)
	if err != nil {
		c.logger.Error(fmt.Sprintf("Init conn by %s", c.Uri))
	} else {
		c.ConnectionWarehouse <- dbConn
	}
}

func (c *MySQLConnectionPool) Dispose() {
	for conn := range c.ConnectionWarehouse {
		conn.Close()
	}
}

func (c *MySQLConnectionPool) CheckAlive() {
	for conn := range c.ConnectionWarehouse {
		pingErr := conn.Ping()
		if pingErr != nil {
			_ = conn.Close()
			c.initConnection()
		} else {
			c.logger.Debug("Connection refresh success....")
		}
	}
}

func (c *MySQLConnectionPool) PopConnection() *sql.DB {
	return <-c.ConnectionWarehouse
}

func (c *MySQLConnectionPool) PushConnection(conn *sql.DB) {
	c.ConnectionWarehouse <- conn
}
