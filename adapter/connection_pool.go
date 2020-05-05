// @Datetime  : 2020/4/4 1:22 下午
// @Author    : psyduck
// @Purpose   :
// @TODO      : Pair programming
//
package adapter

type DBConnectionPool interface {
	Init()
	CheckAlive()
	Dispose()
}

type DBConnectionPoolMap map[int]DBConnectionPool

type ConnectionPoolKeeper struct {
}

func (m *ConnectionPoolKeeper) Run() {
	panic("Not implemented")
}
