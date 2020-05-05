// @Datetime  : 2019-06-18 21:15
// @Author    : psyduck
// @Purpose   :
// @TODO      : Pair programming
//
package common

import (
	"fmt"
	"syscall"
	"testing"
	"time"
)

func TestGetUpSince(t *testing.T) {
	x := syscall.Sysinfo_t{}
	if err := syscall.Sysinfo(&x); err != nil {
		fmt.Println(err)
	}

	// var nowSeconds = time.Now().Unix()
	fmt.Printf("%+v", x)
	fmt.Println(time.Unix(time.Now().Unix()-x.Uptime, 0))
}

func TestGetMemoryTotal(t *testing.T) {
	x := syscall.Sysinfo_t{}
	if err := syscall.Sysinfo(&x); err != nil {
		fmt.Println(err)
	}
	fmt.Println(x.Totalram)
}
