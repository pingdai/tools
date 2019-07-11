package amqp_producer

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

type Producer struct {
	Conn    *amqp.Connection
	Channel *amqp.Channel
}

func NewProducer(amqpURI, exchange, exchangeType string) (*Producer, error) {
	c := &Producer{
		Conn:    nil,
		Channel: nil,
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

	return c, nil
}

// One would typically keep a channel of publishings, a sequence number, and a
// set of unacknowledged sequence numbers and loop until the publishing channel
// is closed.
func ConfirmOne(confirms <-chan amqp.Confirmation) {
	logrus.Infof("waiting for confirmation of one publishing")

	if confirmed := <-confirms; confirmed.Ack {
		logrus.Infof("confirmed delivery with delivery tag: %d", confirmed.DeliveryTag)
	} else {
		logrus.Infof("failed delivery of delivery tag: %d", confirmed.DeliveryTag)
	}
}
