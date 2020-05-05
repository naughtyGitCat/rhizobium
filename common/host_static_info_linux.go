// @Datetime  : 2019-06-18 20:44
// @Author    : psyduck
// @Purpose   :
// @TODO      : Pair programming
//
package common

import (
	"net"
	"os"
	"strconv"
	"syscall"
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
	ret := StaticMemoryInfo{}
	x := syscall.Sysinfo_t{}
	if err := syscall.Sysinfo(&x); err != nil {
		p.Warn(err)
	}
	ret.TotalMemorySizeByte = x.Totalram
	ret.TotalMemorySizeKB = x.Totalram / 1024
	ret.TotalMemorySizeMB = x.Totalram / 1024 * 1024
	ret.TotalMemorySizeGB = x.Totalram / 1024 * 1024 * 1024
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
	x := syscall.Sysinfo_t{}
	if err := syscall.Sysinfo(&x); err != nil {
		p.Warn(err)
	}
	return time.Unix(time.Now().Unix()-x.Uptime, 0)
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
