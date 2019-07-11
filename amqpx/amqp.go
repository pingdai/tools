package amqpx

import (
	"fmt"
	"github.com/pingdai/tools/amqpx/amqp_consumer"
	"github.com/pingdai/tools/amqpx/amqp_producer"
	"github.com/pingdai/tools/constants"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"time"
)

type Amqpx struct {
	Type            constants.MQType `json:"type"` // mq类型，1生产者，2消费者
	Uri             string           `json:"uri"`
	Exchange        string           `json:"exchange"`
	ExchangeType    string           `json:"exchange_type"`
	Queue           string           `json:"queue"`
	RoutingKey      string           `json:"routing_key"`
	ConsumerTag     string           `json:"consumer_tag"`
	ConsumerHandler amqp_consumer.ConsumerHandler

	lifeTime time.Duration // ==0的话将会声明周期将会为永远
	init     bool
	conn     *amqp.Connection
	channel  *amqp.Channel
}

func (amqpx *Amqpx) MarshalDefaults() {
	if amqpx.lifeTime < 0 {
		amqpx.lifeTime = 5 * time.Second
	}

	if amqpx.ConsumerTag == "" {
		amqpx.ConsumerTag = "simple-consumer"
	}

	if amqpx.Uri == "" {
		panic("MQ connect uri cannot be empty")
	}

	if amqpx.Type != constants.MQ_TYPE_PRODUCER &&
		amqpx.Type != constants.MQ_TYPE_CONSUMER {
		panic("MQ Type should be 1:producer 2:consumer")
	}
}

func (amqpx *Amqpx) Init() {
	if !amqpx.init {
		go amqpx.New()
	}
}

func (amqpx *Amqpx) New() {
	if amqpx.init {
		logrus.Infof("RabbitMQ Cli has been initialized")
		return
	}

	logrus.Debugf("开始进行MQ链接")

	amqpx.MarshalDefaults()

	if amqpx.Type == constants.MQ_TYPE_PRODUCER {
		amqpx.doProducer()
	} else {
		amqpx.doConsumer()
	}
}

func (amqpx *Amqpx) doConsumer() {
	c, err := amqp_consumer.NewConsumer(
		amqpx.Uri,
		amqpx.Exchange,
		amqpx.ExchangeType,
		amqpx.Queue,
		amqpx.RoutingKey,
		amqpx.ConsumerTag,
		amqpx.ConsumerHandler,
	)
	if err != nil {
		panic(fmt.Sprintf("%v", err))
	}

	amqpx.conn = c.Conn
	amqpx.channel = c.Channel

	if amqpx.lifeTime > 0 {
		logrus.Infof("running for %s", amqpx.lifeTime)
		time.Sleep(amqpx.lifeTime)
	} else {
		logrus.Infof("running forever")
		amqpx.init = true
		select {}
	}

	logrus.Infof("shutting down")

	if err = c.Shutdown(); err != nil {
		logrus.Errorf("error during shutdown: %s", err)
	}
}

func (amqpx *Amqpx) doProducer() {
	c, err := amqp_producer.NewProducer(
		amqpx.Uri,
		amqpx.Exchange,
		amqpx.ExchangeType,
	)
	if err != nil {
		panic(fmt.Sprintf("%v", err))
	}

	amqpx.conn = c.Conn
	amqpx.channel = c.Channel
}

// 针对生产者调用，发布
func (amqpx *Amqpx) Publish(body string) error {
	// 是否发送确认
	var reliable = true
	if reliable {
		logrus.Infof("enabling publishing confirms.")
		if err := amqpx.channel.Confirm(false); err != nil {
			return fmt.Errorf("Channel could not be put into confirm mode: %s", err.Error())
		}

		confirms := amqpx.channel.NotifyPublish(make(chan amqp.Confirmation, 1))

		defer amqp_producer.ConfirmOne(confirms)
	}

	logrus.Infof("declared Exchange, publishing %dB body (%s)", len(body), body)
	if err := amqpx.channel.Publish(
		amqpx.Exchange,   // publish to an exchange
		amqpx.RoutingKey, // routing to 0 or more queues
		false,            // mandatory
		false,            // immediate
		amqp.Publishing{
			Headers:         amqp.Table{},
			ContentType:     "text/plain",
			ContentEncoding: "",
			Body:            []byte(body),
			DeliveryMode:    amqp.Persistent, // 1=non-persistent, 2=persistent
			Priority:        0,               // 0-9
			// a bunch of application/implementation-specific fields
		},
	); err != nil {
		return fmt.Errorf("Exchange Publish: %s", err.Error())
	}

	return nil
}
