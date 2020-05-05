// @Datetime  : 2019-05-05 21:04
// @Author    : psyduck
// @Purpose   : 读取写入配置文件等操作
// @TODO      : Pair programming
//
package common

import (
	"flag"
	"github.com/BurntSushi/toml"
	"log"
	"os"
	"path/filepath"
)

//func ReadConf(confPath string) model.Config {
func ReadConf() Config {
	confPathPtr := flag.String("config", "/Users/psyduck/Github/rhizobium/config.toml",
		"Config file path(absolute path)")
	flag.Parse()
	// fmt.Printf("Get confPath from cmd flag,%+v",*confPathPtr)

	var conf Config
	//_, err := toml.DecodeFile(confPath, &conf)
	_, err := toml.DecodeFile(*confPathPtr, &conf)
	if err != nil {
		log.Fatalf("%s: %s", err, "Read config file failed")
	}
	conf.FilePath.HomePath, _ = filepath.Abs(filepath.Dir(os.Args[0]))
	conf.FilePath.LogPath = conf.FilePath.HomePath + "/logs/"
	conf.FilePath.AmidePath = conf.FilePath.HomePath + "/amides/"
	return conf
}

//var CONF = ReadConf("./config.toml")
var CONF = ReadConf()

type MySQL struct {
	User      string
	Password  string
	UseSocket bool
}

type RabbitMQ struct {
	Host     string
	Port     int
	VHost    string
	User     string
	Password string
	Exchange string
}

type KafkaLQ struct {
	Host             string
	Port             int
	User             string
	Password         string
	DefaultTopic     string `toml:"default_topic"`
	DefaultPartition int    `toml:"default_partition"`
}

type FilePath struct {
	HomePath  string
	AmidePath string
	LogPath   string
}

type Config struct {
	MySQL    MySQL
	RabbitMQ RabbitMQ
	KafkaLQ  KafkaLQ
	FilePath FilePath
	Misc     Misc
}

type Misc struct {
	Debug bool
}
