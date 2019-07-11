package amqp_consumer

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

type ConsumerHandler func([]byte) error

type Consumer struct {
	Conn    *amqp.Connection
	Channel *amqp.Channel
	tag     string
	done    chan error
}

func (c *Consumer) Shutdown() error {
	// will close() the deliveries channel
	if err := c.Channel.Cancel(c.tag, true); err != nil {
		return fmt.Errorf("Consumer cancel failed: %s", err)
	}

	if err := c.Conn.Close(); err != nil {
		return fmt.Errorf("AMQP connection close error: %s", err)
	}

	defer logrus.Warnf("AMQP shutdown OK")

	// wait for handle() to exit
	return <-c.done
}

func handle(deliveries <-chan amqp.Delivery, done chan error, cb ConsumerHandler) {
	for d := range deliveries {
		logrus.Debugf(
			"got %dB delivery: [%v] %q",
			len(d.Body),
			d.DeliveryTag,
			d.Body,
		)

		if err := cb(d.Body); err != nil {
			logrus.Warnf("mq lib 回调函数返回失败，err:%v", err)
		} else {
			d.Ack(false)
		}
	}
	logrus.Warnf("handle: deliveries channel closed")
	done <- nil
}

func NewConsumer(amqpURI, exchange, exchangeType, queueName, routingKey, ctag string, cb ConsumerHandler) (*Consumer, error) {
	if cb == nil {
		panic("作为消费者，需要注册回调函数")
	}

	c := &Consumer{
		Conn:    nil,
		Channel: nil,
		tag:     ctag,
		done:    make(chan error),
	}

	var err error

	logrus.Debugf("dialing %s", amqpURI)
	c.Conn, err = amqp.Dial(amqpURI)
	if err != nil {
		return nil, fmt.Errorf("Dial: %s", err)
	}

	go func() {
		logrus.Infof("closing: %s", <-c.Conn.NotifyClose(make(chan *amqp.Error)))
	}()

	logrus.Debugf("got Connection, getting Channel")
	c.Channel, err = c.Conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("Channel: %s", err)
	}

	logrus.Debugf("got Channel, declaring Exchange (%s)", exchange)
	if err = c.Channel.ExchangeDeclare(
		exchange,     // name of the exchange
		exchangeType, // type
		true,         // durable
		false,        // delete when complete
		false,        // internal
		false,        // noWait
		nil,          // arguments
	); err != nil {
		return nil, fmt.Errorf("Exchange Declare: %v", err)
	}

	logrus.Infof("declared Exchange, declaring Queue %q", queueName)
	queue, err := c.Channel.QueueDeclare(
		queueName, // name of the queue
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // noWait
		nil,       // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("Queue Declare: %s", err)
	}

	logrus.Infof("declared Queue (%q %d messages, %d consumers), binding to Exchange (key %s)",
		queue.Name, queue.Messages, queue.Consumers, routingKey)

	if err = c.Channel.QueueBind(
		queue.Name, // name of the queue
		routingKey, // bindingKey
		exchange,   // sourceExchange
		false,      // noWait
		nil,        // arguments
	); err != nil {
		return nil, fmt.Errorf("Queue Bind: %s", err)
	}

	logrus.Infof("Queue bound to Exchange, starting Consume (consumer tag %s)", c.tag)
	deliveries, err := c.Channel.Consume(
		queue.Name, // name
		c.tag,      // consumerTag,
		false,      // noAck
		false,      // exclusive
		false,      // noLocal
		false,      // noWait
		nil,        // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("Queue Consume: %s", err)
	}

	go handle(deliveries, c.done, cb)

	return c, nil
}
