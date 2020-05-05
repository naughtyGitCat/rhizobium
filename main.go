// @Datetime  : 2020/4/23 8:36 下午
// @Author    : psyduck
// @Purpose   :
// @TODO      : Pair programming
//
package main

import (
	"rhizobium/common"
	"rhizobium/manager"
	"time"
)

func main() {
	var p = common.GetLoggerNoFields("main")
	// 打印程序版本
	p.Infof("program version: %s", common.Version)

	////////////// 开始执行Rhizobium程序服务 ///////////
	p.Debug("Now do real jobs........................")
	// 程序本身的心跳信息
	go agentHeartbeat()
	for true {
		time.Sleep(1 * time.Second)
		// 这里需要做一下检测ctrl + c，然后立即保存lastPos，不过默认从倒数几行读的话，也无所谓了
	}
}

/*
	Agent自身模块
*/
func agentHeartbeat() {
	x := manager.NewAgentHeartBeat()
	for true {
		x.Run()
		time.Sleep(18 * time.Second)
	}
}
