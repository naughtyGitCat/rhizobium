// @Datetime  : 2019-08-11 15:50
// @Author    : psyduck
// @Purpose   :
// @TODO      : Pair programming
//
package rpc

import (
	"bytes"
	"context"
	"fmt"
	p "github.com/sirupsen/logrus"
	"os/exec"
	"time"
)

func (s *Server) RunLinuxCmd(ctx context.Context, in *RunLinuxCmdRequest) (*RunLinuxCmdResponse, error) {
	var err error
	var progress *exec.Cmd
	var retCode int64 = -2
	var nomBuffer, errBuffer bytes.Buffer
	var retContent string
	// 熔丝
	fuse, cancel := context.WithTimeout(context.Background(), time.Duration(in.ExecTimeout)*time.Second)
	defer cancel()
	p.Infow(fmt.Sprintf("Receive a linux cmd request: %#v", in), "service", "linux")
	// 初始化进程
	if in.ExecUser != "" {
		// TODO: 先检查用户存不存在,若不存在返回错误信息为不存在用户。通过`id -u username`判断返回码是否为0，或者在/etc/passwd中查找
		progress = exec.CommandContext(fuse, "runuser", in.ExecUser, "-c", in.Cmd)
	} else {
		progress = exec.CommandContext(fuse, "bash", "-c", in.Cmd)
	}
	// Path指程序的执行位置,Dir指设定CWD。也就是会在Path下寻找可执行文件位置,但是ls的时候还是对Dir进行ls
	progress.Dir = in.ExecDir
	progress.Stdout = &nomBuffer
	progress.Stderr = &errBuffer
	// 为什么没有直接使用Run()？因为可以更详细的打印每个步骤的错误
	// 启动进程
	if err = progress.Start(); err != nil {
		p.Errorw(fmt.Sprintf("Initialize cmd process %s failed, %s", in.Cmd, err), "service", "linux")
	} else {
		// 等待进程执行结果
		if err = progress.Wait(); err != nil {
			p.Errorf("Get cmd %s failed info, %s", in.Cmd, err)
		}
	}
	// 求出返回编码
	if exitError, ok := err.(*exec.ExitError); ok {
		retCode = int64(exitError.ExitCode())
	} else {
		retCode = int64(progress.ProcessState.ExitCode())
	}
	// 算出返回信息
	// 当err不为nil时要清理err,避免直接抛异常给调用端，让调用方的返回更加整洁
	if err != nil {
		retContent = errBuffer.String() + err.Error()
		err = nil
	} else {
		retContent = nomBuffer.String()
	}
	p.Infow(fmt.Sprintf("Finished linux cmd execution, retCode: %d, retContent: %s", retCode, retContent),
		"service", "linux")
	return &RunLinuxCmdResponse{ReqID: in.ReqID, RetCode: retCode, RetContent: retContent}, err
}
