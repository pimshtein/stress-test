package amqp

import (
	"errors"
	"github.com/streadway/amqp"
	"log"
	"strings"
)

type ConsumerSettings struct {
	Uri       string           `json:"uri"`
	Exchange  ExchangeSettings `json:"exchange"`
	Queue     QueueSettings    `json:"queue"`
	Qos       QosSettings      `json:"qos"`
	Tag       string           `json:"tag"`
	AutoAck   bool             `json:"auto_ack"`
	Exclusive bool             `json:"internal"`
	NoLocal   bool             `json:"no_local"`
	NoWait    bool             `json:"no_wait"`
	Requeue   bool             `json:"requeue"`
}

type ProducerSettings struct {
	Uri       string           `json:"uri"`
	Exchange  ExchangeSettings `json:"exchange"`
	Tag       string           `json:"tag"`
	Mandatory bool             `json:"mandatory"`
	Immediate bool             `json:"immediate"`
}

type ExchangeSettings struct {
	Name       string `json:"name"`
	Type       string `json:"type"`
	Durable    bool   `json:"durable"`
	AutoDelete bool   `json:"auto_delete"`
	Internal   bool   `json:"internal"`
	NoWait     bool   `json:"no_wait"`
}

type QueueSettings struct {
	Name       string `json:"name"`
	Tag        string `json:"tag"`
	Durable    bool   `json:"durable"`
	AutoDelete bool   `json:"auto_delete"`
	Exclusive  bool   `json:"exclusive"`
	NoWait     bool   `json:"no_wait"`
}

type QosSettings struct {
	PrefetchCount int  `json:"prefetch_count"`
	PrefetchSize  int  `json:"prefetch_size"`
	Global        bool `json:"global"`
}

type Consumer struct {
	settings *ConsumerSettings
	handler  func(msg amqp.Delivery) error
}

type Producer struct {
	settings *ProducerSettings
}

func NewConsumer(settings *ConsumerSettings, handler func(msg amqp.Delivery) error) *Consumer {
	return &Consumer{settings: settings, handler: handler}
}

func NewProducer(settings *ProducerSettings) *Producer {
	return &Producer{settings: settings}
}

func (c Consumer) Consume(stop <-chan struct{}) error {
	log.Println("starting consumer")
	if nil == c.settings {
		return errors.New("RabbitMQ consumer must be initialized")
	}
	uri := strings.Trim(c.settings.Uri, "/ ")
	conn, err := amqp.Dial(uri)
	if nil != err {
		return err
	}
	defer func() {
		_ = conn.Close()
	}()

	ch, err := conn.Channel()
	if nil != err {
		return err
	}
	defer func() {
		_ = ch.Close()
	}()

	err = ch.ExchangeDeclare(
		c.settings.Exchange.Name,
		c.settings.Exchange.Type,
		c.settings.Exchange.Durable,
		c.settings.Exchange.AutoDelete,
		c.settings.Exchange.Internal,
		c.settings.Exchange.NoWait,
		nil,
	)
	if nil != err {
		return err
	}

	q, err := ch.QueueDeclare(
		c.settings.Queue.Name,       // name
		c.settings.Queue.Durable,    // durable
		c.settings.Queue.AutoDelete, // delete when unused
		c.settings.Queue.Exclusive,  // exclusive
		c.settings.Queue.NoWait,     // no-wait
		nil,                         // arguments
	)
	if nil != err {
		return err
	}

	err = ch.QueueBind(
		q.Name,
		c.settings.Queue.Tag,
		c.settings.Exchange.Name,
		c.settings.Queue.NoWait,
		nil,
	)
	if nil != err {
		return err
	}

	err = ch.Qos(
		c.settings.Qos.PrefetchCount, // prefetch count
		c.settings.Qos.PrefetchSize,  // prefetch size
		c.settings.Qos.Global,        // global
	)
	if nil != err {
		return err
	}

	messages, err := ch.Consume(
		q.Name,               // queue
		c.settings.Tag,       // consumer
		c.settings.AutoAck,   // auto-ack
		c.settings.Exclusive, // exclusive
		c.settings.NoLocal,   // no-local
		c.settings.NoWait,    // no-wait
		nil,                  // args
	)
	if nil != err {
		return err
	}

	log.Println("consuming")
	defer log.Println("consumer exit")

	h := func(d amqp.Delivery) {
		err := c.handler(d)
		if nil != err {
			log.Println(err)
			if !c.settings.AutoAck {
				_ = d.Reject(c.settings.Requeue)
			}
		} else if !c.settings.AutoAck {
			_ = d.Ack(false)
		}
	}

	channelClose := make(chan *amqp.Error)
	ch.NotifyClose(channelClose)
	for {
		select {
		case delivery := <-messages:
			go h(delivery)
		case <-stop:
			log.Println("consumer stopped")
			return nil
		case err = <-channelClose:
			log.Println("channel closed")
			return err
		}
	}
}

func (p Producer) PublishMessage(message []byte, routingKey string, contentType string) error {
	if nil == p.settings {
		return errors.New("RabbitMQ consumer must be initialized")
	}

	uri := strings.Trim(p.settings.Uri, "/ ")
	conn, err := amqp.Dial(uri)

	if nil != err {
		return err
	}

	defer func() {
		_ = conn.Close()
	}()

	ch, err := conn.Channel()
	if nil != err {
		return err
	}

	defer func() {
		_ = ch.Close()
	}()

	err = ch.ExchangeDeclare(
		p.settings.Exchange.Name,
		p.settings.Exchange.Type,
		p.settings.Exchange.Durable,
		p.settings.Exchange.AutoDelete,
		p.settings.Exchange.Internal,
		p.settings.Exchange.NoWait,
		nil,
	)
	if nil != err {
		return err
	}

	return ch.Publish(
		p.settings.Exchange.Name,
		routingKey,
		p.settings.Mandatory,
		p.settings.Immediate,
		amqp.Publishing{
			Headers:         amqp.Table{},
			ContentType:     contentType,
			ContentEncoding: "",
			Body:            message,
			DeliveryMode:    amqp.Transient,
			Priority:        0,
		},
	)
}
