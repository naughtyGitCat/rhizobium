// @Datetime  : 2019-06-18 20:44
// @Author    : psyduck
// @Purpose   :
// @TODO      : Pair programming
//
package common

import (
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"net"
	"os"
	"strconv"
	"time"
)

type CPUInfo struct {
}

type StaticMemoryInfo struct {
	TotalMemorySizeByte uint64
	TotalMemorySizeKB   uint64
	TotalMemorySizeMB   uint64
	TotalMemorySizeGB   uint64
}

type OSInfo struct {
}

type OEMInfo struct {
}

type HostStaticParameters struct {
	IP         string
	HostName   string
	MemoryInfo StaticMemoryInfo
	CPUInfo    CPUInfo
	OSInfo     OSInfo
	OEMInfo    OEMInfo
	UpSince    time.Time
}

/*
获取内存信息(固定不变的只有总大小)
*/
func GetStaticMemoryInfo() StaticMemoryInfo {
	x, err := mem.VirtualMemory()
	FailOnError(err, "Fetch system memory info via psutil failed")
	ret := StaticMemoryInfo{}
	ret.TotalMemorySizeByte = x.Total
	ret.TotalMemorySizeKB = x.Total / 1024
	ret.TotalMemorySizeMB = x.Total / 1024 * 1024
	ret.TotalMemorySizeGB = x.Total / 1024 * 1024 * 1024
	return ret
}

func GetCPUInfo() {

}

func GetOSInfo() {

}

func GetOEMInfo() {

}

/*
获取系统启动时间
*/
func GetUpSince() time.Time {
	seconds, err := host.BootTime()
	FailOnError(err, "Fetch sys uptime with psutil failed")
	return time.Unix(int64(seconds), 0)
}

/*
获取本地主机名
*/
func GetHostName() string {
	name, _ := os.Hostname()
	return name
}

/*
获取本地IP地址
*/
func GetLocalIP() net.IP {
	judgeServer := CONF.RabbitMQ.Host + ":" + strconv.Itoa(CONF.RabbitMQ.Port)
	// judgeServer := "4.4.4.4:53"
	// p.Debug(judgeServer)
	conn, err := net.Dial("udp", judgeServer)
	if err != nil {
		p.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	p.Debugf("%+v", localAddr.String())
	return localAddr.IP
}
