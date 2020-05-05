// @Datetime  : 2020/4/23 9:01 下午
// @Author    : psyduck
// @Purpose   :
// @TODO      : Pair programming
//
package common

import (
	"github.com/sirupsen/logrus"
)

var LocalIP = GetLocalIP().String()

var p = GetLogger("common", logrus.Fields{})
