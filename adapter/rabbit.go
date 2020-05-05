// @Datetime  : 2019-06-25 17:12
// @Author    : psyduck
// @Purpose   :
// @TODO      : 设置连接保存时间
//
package adapter

import (
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"rhizobium/common"
	"sync"
)

var pr = common.GetLoggerNoFields("rabbit")

var RabbitRouteMap = map[string]string{
	"ddb_log":        "ddb_log_queue",
	"agent_info":     "rhizobium_info_queue",
	"rhizobium_info": "rhizobium_info_queue",
	"ddb_qs_config":  "ddb_qs_config",
	"mongo_log":      "mongo_log_queue",
}

type RabbitMessage struct {
	RouteKey     string
	QueueName    string
	BytesMessage []byte
}

type RabbitMQ struct {
	rabbitConnect         *amqp.Connection
	rabbitChannel         *amqp.Channel
	messagePublishChannel chan RabbitMessage
	messageConsumeChannel chan RabbitMessage
	basicInfo             struct {
		host     string
		port     int
		user     string
		password string
		vhost    string
		exchange string
	}
}

// func NewRabbitMQ(CONF model.Config, pubChan chan RabbitMessage) RabbitMQ {
func NewRabbitMQ() (RabbitMQ, error) {
	var err error
	var CONF = common.CONF
	var r RabbitMQ
	r.basicInfo.host = CONF.RabbitMQ.Host
	r.basicInfo.port = CONF.RabbitMQ.Port
	r.basicInfo.password = CONF.RabbitMQ.Password
	r.basicInfo.vhost = CONF.RabbitMQ.VHost
	r.basicInfo.user = CONF.RabbitMQ.User
	r.basicInfo.exchange = CONF.RabbitMQ.Exchange
	r.messagePublishChannel = RabbitPublishChan

	err = r.dialConnection()
	if err != nil {
		return r, err
	}
	err = r.openChannel()
	if err != nil {
		return r, err
	}
	err = r.declareExchange()
	if err != nil {
		return r, err
	}
	// 先声明队列试试
	err = r.declareQueue()
	if err != nil {
		return r, err
	}
	err = r.bindQueueRoute()
	if err != nil {
		return r, err
	}
	// r.declareQueue()
	pr.Debug("A rabbit is been caught")
	return r, err
}

func (r *RabbitMQ) dialConnection() error {
	var err error
	uri := fmt.Sprintf("amqp://%s:%s@%s:%d/%s",
		r.basicInfo.user,
		r.basicInfo.password,
		r.basicInfo.host,
		r.basicInfo.port,
		r.basicInfo.vhost)

	r.rabbitConnect, err = amqp.Dial(uri)
	if err != nil {
		pr.Warnf("拨号AMQP时出错,%s", err)
	}
	return err
}

func (r *RabbitMQ) openChannel() error {
	var err error
	r.rabbitChannel, err = r.rabbitConnect.Channel()
	if err != nil {
		pr.Warnf("初始化频道时出错,%s", err)
	}
	return err
}

func (r *RabbitMQ) declareExchange() error {
	var err error
	err = r.rabbitChannel.ExchangeDeclare(
		r.basicInfo.exchange, // name
		"direct",             // type
		true,                 // durable
		false,                // auto-deleted
		false,                // internal
		false,                // no-wait
		nil,                  // arguments
	)
	if err != nil {
		pr.Warnf("定义消息交换机时出错,%s", err)
		return err
	}
	return err
}

func (r *RabbitMQ) bindQueueRoute() error {
	var err error
	for routeKey := range RabbitRouteMap {
		err = r.rabbitChannel.QueueBind(
			RabbitRouteMap[routeKey], // queue name
			routeKey,                 // routing key
			r.basicInfo.exchange,     // exchange
			false,
			nil)
		if err != nil {
			pr.Warnf("映射队列与路由键失败,%s", err)
			return err
		}
	}
	return err
}

func (r *RabbitMQ) declareQueue() error {
	var err error
	args := make(amqp.Table)
	args["x-message-ttl"] = 3600000
	for routeKey := range RabbitRouteMap {
		_, err := r.rabbitChannel.QueueDeclare(
			RabbitRouteMap[routeKey], // name
			true,                     // durable
			false,                    // delete when unused
			false,                    // exclusive
			false,                    // no-wait
			args,                     // arguments,
		)
		if err != nil {
			pr.Warnf("声明队列失败,%s", err)
			return err
		}
	}
	return err
}

func (r *RabbitMQ) PublishJson() {
	for msg := range r.messagePublishChannel {
		err := r.rabbitChannel.Publish(
			r.basicInfo.exchange, // exchange
			msg.RouteKey,         // routing key
			false,                // mandatory
			false,                // immediate
			amqp.Publishing{
				ContentType: "application/json",
				Body:        msg.BytesMessage,
			})
		pr.Debugf(" [x] Sent %s\n", string(msg.BytesMessage))
		if err != nil {
			pr.Warnf("消息发送失败,%s", err)
		}
	}
}

type RabbitMan struct {
	RabbitHole            chan *amqp.Channel
	messagePublishChannel chan RabbitMessage
	RabbitHoleLatch       sync.Mutex
	logger                *logrus.Logger
}

func HireRabbitMan() *RabbitMan {
	var m = RabbitMan{}
	m.RabbitHole = make(chan *amqp.Channel, 8)
	m.messagePublishChannel = RabbitPublishChan
	m.logger = common.GetLoggerNoFields("rabbit")
	return &m
}

/**
抓一只兔子
*/
func (m *RabbitMan) catchARabbit() (*amqp.Channel, error) {
	m.logger.Debug("Now try to catch a rabbit")
	rabbit, err := NewRabbitMQ()
	if err != nil {
		return nil, err
	}
	return rabbit.rabbitChannel, err
}

/**

 */
func (m *RabbitMan) borrowARabbit() (*amqp.Channel, error) {
	m.logger.Debug("We need a health rabbit")
	m.RabbitHoleLatch.Lock()
	if len(m.RabbitHole) == 0 {
		pr.Debug("There is no rabbit, try to catch one")
		for i := 0; i < 8; i++ {
			rabbit, err := NewRabbitMQ()
			if err != nil {
				pr.Debugf("It's a dead rabbit....")
			} else {
				pr.Debug("Now push the rabbit to rabbit hole")
				m.RabbitHole <- rabbit.rabbitChannel
			}
		}
	}
	if len(m.RabbitHole) == 0 {
		m.RabbitHoleLatch.Unlock()
		return nil, errors.New("failed to catch rabbits")
	}
	for rabbit := range m.RabbitHole {
		m.logger.Debug("Now pull a rabbit from hole")
		m.RabbitHoleLatch.Unlock()
		return rabbit, nil
	}
	m.RabbitHoleLatch.Unlock()
	return nil, errors.New("failed to catch rabbits")
}

func (m *RabbitMan) giveBackARabbit(r *amqp.Channel) {
	m.logger.Debugf("Put the rabbit back to rabbit hole")
	m.RabbitHole <- r
}

func (m *RabbitMan) CarryCarrots() {
	for msg := range m.messagePublishChannel {
		pr.Debugf("Now get a %s msg from carrot stack %s", msg.RouteKey, msg.BytesMessage)
		t, err := m.borrowARabbit()
		if err == nil {
			err = t.Publish(
				common.CONF.RabbitMQ.Exchange, // exchange
				msg.RouteKey,                  // routing key
				false,                         // mandatory
				false,                         // immediate
				amqp.Publishing{
					ContentType: "application/json",
					Body:        msg.BytesMessage,
				})
			if err == nil {
				pr.Debugf(" [x] Sent %s\n", string(msg.BytesMessage))
				m.giveBackARabbit(t)
			} else {
				pr.Warnf("%s: %s", "Rabbit failed to carry a carrot, discard the carrot", err)
			}
		} else {
			m.logger.Warnf("Call rabbit to carry a carrot failed, discard the carrot")
		}
	}
}
