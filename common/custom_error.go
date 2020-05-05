// @Datetime  : 2019-06-03 16:02
// @Author    : psyduck
// @Purpose   :
// @TODO      : Pair programming
//
package common

import (
	"errors"
	"fmt"
)

var (
	FatalError  = errors.New("unrecoverable error occurs")
	SlightError = errors.New("some error occurs, but it does not affect the progress running")
)

/*
自定义错误
*/
func CustomError(errMsg string, errLevel string) error {
	return errors.New(fmt.Sprintf("[%s ERROR], %s", errLevel, errMsg))
}

/*
功能，属性，接口未定义错误
*/
func UnimplementedError(targetName string, interfaceName string) error {
	return errors.New(fmt.Sprintf("[%s ERROR], %s haven't implement interface %s", "Fatal", targetName, interfaceName))
}

/*
操作文件错误
*/
func FileError(fileName string, err error, operation string) error {
	errMsg := fmt.Sprintf("%s file %s failed, %s", operation, fileName, err)
	return errors.New(errMsg)
}

func CustomWarn(errMsg string, errLevel string) error {
	return errors.New(fmt.Sprintf("[%s Warn], %s", errLevel, errMsg))
}

func RabbitMQError(errMsg string, errLevel string) error {
	return errors.New(fmt.Sprintf("[%s ERROR], %s", errLevel, errMsg))
}
