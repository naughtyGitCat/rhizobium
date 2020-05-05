// @Datetime  : 2019-08-11 16:23
// @Author    : psyduck
// @Purpose   :
// @TODO      : Pair programming
//
package rpc

import (
	"Rhizobium/common"
	"Rhizobium/common/logp"
	"fmt"
)

type Server struct {
}

func getLogger(name string) *logp.Logger {
	conf := logp.DefaultConfig()
	conf.Level = logp.DebugLevel
	conf.Files.Path = common.CONF.FilePath.LogPath
	conf.Files.RotateOnStartup = false
	conf.Files.Name = "rpc.log"
	if err := logp.Configure(conf); err != nil {
		err := fmt.Errorf("error initializing logging: %v", err)
		panic(err)
	}
	return logp.NewLogger(name)
}

var p = getLogger("rpc")
