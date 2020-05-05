// @Datetime  : 2020/4/23 9:09 下午
// @Author    : psyduck
// @Purpose   :
// @TODO      : Pair programming
//
package common

import (
	"fmt"
	rotateLogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	path2 "path"
	"time"
)

func GetLogger(fileNamePrefix string, fields logrus.Fields) *logrus.Logger {
	path := path2.Join(CONF.FilePath.LogPath, fmt.Sprintf("%s.log", fileNamePrefix))
	writer, _ := rotateLogs.New(
		path+".%Y%m%d%H%M",
		rotateLogs.WithLinkName(path),
		rotateLogs.WithMaxAge(time.Duration(86400)*time.Second),
		rotateLogs.WithRotationTime(time.Duration(604800)*time.Second),
	)
	p := logrus.New()
	if CONF.Misc.Debug == true {
		p.SetLevel(logrus.DebugLevel)
	}
	p.WithFields(fields)
	p.AddHook(lfshook.NewHook(
		lfshook.WriterMap{
			logrus.DebugLevel: writer,
			logrus.InfoLevel:  writer,
			logrus.WarnLevel:  writer,
			logrus.ErrorLevel: writer,
			logrus.FatalLevel: writer,
			logrus.PanicLevel: writer,
		},
		&logrus.TextFormatter{
			FullTimestamp: true,
			ForceColors:   true,
			ForceQuote:    false,
			PadLevelText:  true,
		},
	))
	return p
}

func GetLoggerNoFields(fileNamePrefix string) *logrus.Logger {
	return GetLogger(fileNamePrefix, logrus.Fields{})
}
