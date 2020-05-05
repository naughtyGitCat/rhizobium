// @Datetime  : 2019/10/5 11:55 上午
// @Author    : psyduck
// @Purpose   :
// @TODO      : Pair programming
//
package adapter

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/plain"
	"github.com/sirupsen/logrus"
	"rhizobium/common"
)

var pk = common.GetLogger("kafka", logrus.Fields{})

var KafkaTopicMap = map[string]string{
	"DDBServiceLogTopic":     "ddb_log",
	"MongoServiceLogTopic":   "mongo_log",
	"RedisServiceLogTopic":   "redis_log",
	"MySQLServiceErrorTopic": "mysql_error",
	"MySQLServiceSlowTopic":  "mysql_slow",
}

type KafkaRawMSG struct {
	Topic      string
	KeyBytes   []byte
	ValueBytes []byte
}

type Ungeziefer map[string]*kafka.Writer

/*
	初始化KAFKA
*/
func InitKafka() Ungeziefer {
	ungeziefer := make(map[string]*kafka.Writer)
	for key, value := range KafkaTopicMap {
		pk.Debugf("Now spawn kafka %s writer %s", key, value)
		ungeziefer[value] = newKafkaWriter(value)
	}
	return ungeziefer
}

/*
	发布消息
*/
func (u Ungeziefer) PublishLog() {
	for msg := range KafkaPublishChan {
		w := u[msg.Topic]
		err := w.WriteMessages(context.Background(), kafka.Message{Key: msg.KeyBytes, Value: msg.ValueBytes})
		pk.Debugf(" [x] Sent %#s\n", msg)
		if err != nil {
			pk.Fatalf("Failed to publish a message: %s", msg, err)
		}
	}
}

/*
	生成kafka writer
*/
func newKafkaWriter(topic string) *kafka.Writer {
	d := &kafka.Dialer{SASLMechanism: plain.Mechanism{Username: common.CONF.KafkaLQ.User,
		Password: common.CONF.KafkaLQ.Password}}
	w := kafka.NewWriter(kafka.WriterConfig{
		Dialer:     d,
		Brokers:    []string{fmt.Sprintf("%s:%d", common.CONF.KafkaLQ.Host, common.CONF.KafkaLQ.Port)},
		Topic:      topic,
		BatchBytes: 8 * 1024 * 1024, // 8M
	})
	return w
}
