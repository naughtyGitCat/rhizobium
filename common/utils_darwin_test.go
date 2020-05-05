// @Datetime  : 2019-06-13 17:15
// @Author    : psyduck
// @Purpose   :
// @TODO      : Pair programming
//
package common

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"runtime"
	"syscall"
	"testing"
	"time"
)

/*
获取文件的stat信息，MacOS
*/
func GetDarwinFileState1(info os.FileInfo) *FileStat1 {
	stat := info.Sys().(*syscall.Stat_t)
	fileState := &FileStat1{
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

func TestGetDarwinFileState1(t *testing.T) {
	f, _ := os.Stat("/Users/psyduck/GitLab/Rhizobium/README.md")
	x := GetDarwinFileState1(f)
	fmt.Printf("%+v\n", x)
}

type FileStat1 struct {
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

func TestGetOSName(t *testing.T) {
	log.Println(runtime.GOOS)
}
