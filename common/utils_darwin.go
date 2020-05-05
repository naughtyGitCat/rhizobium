// @Datetime  : 2019-05-10 17:20
// @Author    : psyduck
// @Purpose   :
// @TODO      :
//
package common

import (
	"encoding/xml"
	"fmt"
	"github.com/sirupsen/logrus"
	"golang.org/x/text/encoding/simplifiedchinese"
	"io/ioutil"
	"net"
	"os"
	"os/user"
	"reflect"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"
)

/*
中国时区
*/
func ChinaTimezone() *time.Location {
	ret, _ := time.LoadLocation("Asia/Chongqing")
	return ret
}

func FailOnError(err error, msg string) {
	if err != nil {
		p.Fatalf("%s: %s", msg, err)
	}
}

func FailOnErrorNew(fp *logrus.Logger, msg string, err error) {
	if err != nil {
		fp.Fatalf("%s: %s", msg, err)
	}
}

type DDBXMLConfig struct {
	XMLName xml.Name `xml:"cluster"`
	Text    string   `xml:",chardata"`
	Master  struct {
		Text               string `xml:",chardata"`
		Name               string `xml:"name"`
		Ip                 string `xml:"ip"`
		Port               string `xml:"port"`
		DbaPort            string `xml:"dba_port"`
		DbnHbInterval      string `xml:"dbn_hb_interval"`
		DbnReportInterval  string `xml:"dbn_report_interval"`
		DeadCheckInterval  string `xml:"dead_check_interval"`
		DeadAssureInterval string `xml:"dead_assure_interval"`
		XabCheckInterval   string `xml:"xab_check_interval"`
		XabTimeout         string `xml:"xab_timeout"`
		XabRetryTimes      string `xml:"xab_retry_times"`
		XabCommitInterval  string `xml:"xab_commit_interval"`
		AlarmSwitch        string `xml:"alarm_switch"`
		SocketTimeout      string `xml:"socket_timeout"`
		ConnectTimeout     string `xml:"connect_timeout"`
		MigUnit            string `xml:"mig_unit"`
		SysdbURL           string `xml:"sysdb_url"`
		SysdbUser          string `xml:"sysdb_user"`
		SysdbPassword      string `xml:"sysdb_password"`
		Pid                string `xml:"pid"`
	} `xml:"master"`
	Client struct {
		Text                string `xml:",chardata"`
		BufferSize          string `xml:"buffer_size"`
		MaxBlocksPerFile    string `xml:"max_blocks_per_file"`
		LogFileDir          string `xml:"log_file_dir"`
		LogFileName         string `xml:"log_file_name"`
		MaxLogFiles         string `xml:"max_log_files"`
		ReportInterval      string `xml:"report_interval"`
		UseDaemon           string `xml:"use_daemon"`
		MaxConnsPerPool     string `xml:"max_conns_per_pool"`
		MaxConnsPerXaPool   string `xml:"max_conns_per_xa_pool"`
		WaitConnTimeout     string `xml:"wait_conn_timeout"`
		ConnIdleTimeout     string `xml:"conn_idle_timeout"`
		MaxPstPerConn       string `xml:"max_pst_per_conn"`
		AsyncTimeout        string `xml:"async_timeout"`
		AsyncInterval       string `xml:"async_interval"`
		AsyncXidlistMaxSize string `xml:"async_xidlist_max_size"`
		AsyncThreadInterval string `xml:"async_thread_interval"`
		BufferFlushInterval string `xml:"buffer_flush_interval"`
		ImmediateFlush      string `xml:"immediate_flush"`
		ExecTimeout         string `xml:"exec_timeout"`
	} `xml:"client"`
}

/*
解析DDB配置文件到model.DDBXMLConfig结构体
*/
func ParseDDBXML(filePath string) DDBXMLConfig {
	xmlFile, err := os.Open(filePath)
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}
	defer xmlFile.Close()
	byteValue, _ := ioutil.ReadAll(xmlFile)
	var fuck DDBXMLConfig
	err = xml.Unmarshal(byteValue, &fuck)
	if err != nil {
		fmt.Println(err)
	}
	return fuck
}

/*
 获取当前包名,好像有点呆，返回了当前utils的包名,这个应该放到每个包的init中,就没问题了
*/
func PackageName() string {
	type Fuck struct {
		// 这是空的
	}
	return reflect.TypeOf(Fuck{}).PkgPath()
}

/*
可以唯一地标识一个文件
*/
type FileID struct {
	Inode  uint64 `json:"inode,"`
	Device uint64 `json:"device,"`
}

/*
获取文件的inode和磁盘位置信息(简略)
*/
func GetSimplyFileState(info os.FileInfo) FileID {
	stat := info.Sys().(*syscall.Stat_t)

	// Convert inode and dev to uint64 to be cross platform compatible
	fileState := FileID{
		Inode:  uint64(stat.Ino),
		Device: uint64(stat.Dev),
	}

	return fileState
}

type FileStat struct {
	Inode     uint64
	Device    uint64
	Name      string
	Size      int64 // bytes
	Mode      os.FileMode
	CrtTime   time.Time // 创建时间
	AcsTime   time.Time // 上次读取时间
	ModTime   time.Time // 上次修改时间
	isDir     bool
	user      user.User
	userGroup user.Group
}

/*
获取文件的stat信息，MacOS
*/
func GetFileState(filePath string) FileStat {
	info, _ := os.Stat(filePath)
	stat := info.Sys().(*syscall.Stat_t)
	fileState := FileStat{
		Inode:   uint64(stat.Ino),
		Device:  uint64(stat.Dev),
		Name:    info.Name(),
		Size:    info.Size(),
		Mode:    info.Mode(),
		CrtTime: time.Unix(stat.Ctimespec.Sec, stat.Ctimespec.Nsec),
		ModTime: time.Unix(stat.Mtimespec.Sec, stat.Mtimespec.Nsec),
		AcsTime: time.Unix(stat.Atimespec.Sec, stat.Atimespec.Nsec),
		isDir:   info.IsDir(),
	}
	groupPointer, _ := user.LookupGroupId(fmt.Sprintf("%d", stat.Gid))
	fileState.userGroup = *groupPointer

	userPointer, _ := user.LookupId(fmt.Sprintf("%d", stat.Uid))
	fileState.user = *userPointer
	return fileState
}

/*
确认两个文件是否相同
*/
func (fs FileID) IsSame(state FileID) bool {
	return fs.Inode == state.Inode && fs.Device == state.Device
}

/*
FileID.__str__
*/
func (fs FileID) String() string {
	var buf [64]byte
	current := strconv.AppendUint(buf[:0], fs.Inode, 10)
	current = append(current, '-')
	current = strconv.AppendUint(current, fs.Device, 10)
	return string(current)
}

/*
只读模式打开文件
*/
func ReadOpen(path string) (*os.File, error) {
	flag := os.O_RDONLY
	perm := os.FileMode(0644)
	return os.OpenFile(path, flag, perm)
}

/*
读写方式打开文件
*/
func WriteOpen(path string) (*os.File, error) {
	flag := os.O_CREATE | os.O_RDWR
	return os.OpenFile(path, flag, 0644)
}

// IsRemoved checks whether the file held by f is removed.
func IsRemoved(f *os.File) bool {
	_, err := os.Stat(f.Name())
	return err != nil
}

// GBK bytes转UTF-8
func ConvertByte2String(byte []byte, charset string) string {
	const (
		UTF8    = string("UTF-8")
		GB18030 = string("GB18030")
	)

	var str string
	switch charset {
	case GB18030:
		var decodeBytes, _ = simplifiedchinese.GB18030.NewDecoder().Bytes(byte)
		str = string(decodeBytes)
	case UTF8:
		fallthrough
	default:
		str = string(byte)
	}

	return str
}

// DropCR drops a terminal \r from the data.
func DropCR(data []byte) []byte {
	if len(data) > 0 && data[len(data)-1] == '\r' {
		return data[0 : len(data)-1]
	}
	return data
}

/**
正则切割行
*/
func RegexSplitLines(data []byte, atEOF bool, regexDelimiter string) (advance int, token []byte, err error) {
	t := regexp.MustCompile(regexDelimiter)
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if infoSlice := t.FindStringIndex(string(data)); len(infoSlice) == 2 {
		// We have a full newline-terminated line.
		pos := infoSlice[0]
		return pos + 1, DropCR(data[0:pos]), nil
	}
	// If we're at EOF, we have a final, non-terminated line. Return it.
	if atEOF {
		return len(data), DropCR(data), nil
	}
	// Request more data.
	return 0, nil, nil
}

// 普通分割行
func SplitLines(data []byte, atEOF bool, delimiter string) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if i := strings.Index(string(data), delimiter); i >= 0 {
		// We have a full newline-terminated line.
		return i + 1, DropCR(data[0:i]), nil
	}
	// If we're at EOF, we have a final, non-terminated line. Return it.
	if atEOF {
		return len(data), DropCR(data), nil
	}
	// Request more data.
	return 0, nil, nil
}

func GetOSName() string {
	return runtime.GOOS
}

/*
	检查URL是否存活，检查端口是否存活
*/
func CheckIPPortAlive(ip string, port int) bool {
	url := fmt.Sprintf("%s:%d", ip, port)
	timeout := time.Duration(1 * time.Second)
	conn, err := net.DialTimeout("tcp", url, timeout)
	defer conn.Close()

	if err, ok := err.(*net.OpError); ok && err.Timeout() {
		p.Warnf("Dial to %s timeout: %s", url, err)
		return false
	}

	if err != nil {
		// Log or report the error here
		p.Errorf("Dial to %s failed : %s", url, err)
		return false
	}
	return true
}
