// @Datetime  : 2019-06-20 13:42
// @Author    : psyduck
// @Purpose   :
// @TODO      : Pair programming
//
package common

import (
	"fmt"
	"os"
	"os/user"
	"syscall"
	"testing"
	"time"
)

func TestGetFileState(t *testing.T) {
	f, _ := os.Stat("/home/dba/DDB.README")
	x := GetFileState1(f)
	fmt.Printf("%+v\n", f)
	fmt.Printf("%+v", x)
}

func GetFileState1(info os.FileInfo) FileStat1 {
	stat := info.Sys().(*syscall.Stat_t)
	fileState := FileStat1{
		Inode:   uint64(stat.Ino),
		Device:  uint64(stat.Dev),
		Name:    info.Name(),
		Size:    info.Size(),
		Mode:    info.Mode(),
		CrtTime: time.Unix(int64(stat.Ctim.Sec), int64(stat.Ctim.Nsec)),
		ModTime: time.Unix(stat.Mtim.Sec, stat.Mtim.Nsec),
		AcsTime: time.Unix(stat.Atim.Sec, stat.Atim.Nsec),
		isDir:   info.IsDir(),
	}
	groupPointer, _ := user.LookupGroupId(fmt.Sprintf("%d", stat.Gid))
	fileState.userGroup = *groupPointer

	userPointer, _ := user.LookupId(fmt.Sprintf("%d", stat.Uid))
	fileState.user = *userPointer
	return fileState
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
