// @Datetime  : 2019-06-18 15:03
// @Author    : psyduck
// @Purpose   :
// @TODO      : Pair programming
//
package manager

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"rhizobium/adapter"
	"rhizobium/common"
	"time"
)

const AgentHeartbeatRouteKey = "agent_info"

type AgentHeartbeat struct {
	pid     int
	message struct {
		MemoryUsage    int       `json:"memory_usage"` // rss kb
		CPUUsage       int       `json:"cpu_usage"`
		AgentVersion   string    `json:"agent_version"`
		HostIP         string    `json:"host_ip"`
		HostName       string    `json:"host_name"`
		AgentStartTime time.Time `json:"uptime"` // 启动时间
		AgentHeartbeat time.Time `json:"heartbeat"`
	}
	logger *logrus.Logger
}

func (a *AgentHeartbeat) collectInfo() {
	// 获取agent本身资源占用信息
	a.pid = os.Getpid()
	agentStat, err := common.GetStat(a.pid)
	common.FailOnError(err, "Fetch pid mem cpu info failed, %+v")
	a.message.MemoryUsage = int(agentStat.Memory)
	a.message.CPUUsage = int(agentStat.CPU)
	a.message.AgentStartTime = agentStat.StartTime
	// 获取其他信息
	a.message.AgentHeartbeat = time.Now().Local()
	a.message.HostIP = common.LocalIP
	a.message.HostName = common.GetHostName()
	// django严格按照mysql的推荐范围1000-01-01到9999-12-32来算，所以这边不能空置为0001-01-01
	// 默认设置为linux的最小时间，即1970-01-01，让django可以正常消费
	// a.message.AgentStartTime = time.Unix(0,0)
}

func (a *AgentHeartbeat) pushMessage() {
	// pa.Debug("Now in PutMessage")
	jsonBytes, err := json.Marshal(a.message)
	// pa.Debug("Now after json marshal")
	common.FailOnError(err, "Serialize struct to json bytes failed")
	a.logger.Debug(fmt.Sprintf("Now send %s", jsonBytes))
	adapter.RabbitPublishChan <- adapter.RabbitMessage{RouteKey: AgentHeartbeatRouteKey, BytesMessage: jsonBytes}
	// adapter.PushMessage(AgentHeartbeatRouteKey, jsonBytes)
}

func NewAgentHeartBeat() AgentHeartbeat {
	x := AgentHeartbeat{}
	x.message.AgentVersion = common.Version
	x.logger = common.GetLogger("agent", logrus.Fields{})
	return x
}

func (a *AgentHeartbeat) Run() {
	a.collectInfo()
	a.pushMessage()
}
