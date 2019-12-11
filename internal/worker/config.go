package worker

import (
	"stress-test/internal/amqp"
	_ "stress-test/internal/amqp"
)

type Config struct {
	Rabbit           RabbitOptions
	ParallelRequests ParallelRequestsOptions
}

type ParallelRequestsOptions struct {
	GoNum       int `required:"true" long:"go-num" env:"CONCURRENCY" description:"The number of goroutines"`
	RequestsNum int `required:"true" long:"requests-num" env:"MESSAGES_NUM" description:"The number of requests"`
}

type RabbitOptions struct {
	Uri                string `required:"true" long:"amqp-url" env:"AMQP_CONSUMER_URI" description:"RabbitMQ URI"`
	ConsumerTag        string `long:"amqp-consumer-tag" env:"AMQP_CONSUMER_TAG" description:"See RabbitMQ docs"`
	ConsumerAutoAck    bool   `long:"amqp-consumer-auto-ack" env:"AMQP_CONSUMER_AUTO_ACK" description:"See RabbitMQ docs"`
	ConsumerExclusive  bool   `long:"amqp-consumer-exclusive" env:"AMQP_CONSUMER_EXCLUSIVE" description:"See RabbitMQ docs"`
	ConsumerNoLocal    bool   `long:"amqp-consumer-no-local" env:"AMQP_CONSUMER_NO_LOCAL" description:"See RabbitMQ docs"`
	ConsumerNoWait     bool   `long:"amqp-consumer-no-wait" env:"AMQP_CONSUMER_NO_WAIT" description:"See RabbitMQ docs"`
	ConsumerRequeue    bool   `long:"amqp-consumer-requeue" env:"AMQP_CONSUMER_REQUEUE" description:"See RabbitMQ docs"`
	ExchangeName       string `required:"true" long:"amqp-exchange-name" env:"AMQP_EXCHANGE_NAME" description:"See RabbitMQ docs"`
	ExchangeType       string `long:"amqp-exchange-type" env:"AMQP_EXCHANGE_TYPE" default:"direct" description:"See RabbitMQ docs"`
	ExchangeDurable    bool   `long:"amqp-exchange-durable" env:"AMQP_EXCHANGE_DURABLE" description:"See RabbitMQ docs"`
	ExchangeAutoDelete bool   `long:"amqp-exchange-auto-delete" env:"AMQP_EXCHANGE_AUTO_DELETE" description:"See RabbitMQ docs"`
	ExchangeInternal   bool   `long:"amqp-exchange-internal" env:"AMQP_EXCHANGE_INTERNAL" description:"See RabbitMQ docs"`
	ExchangeNoWait     bool   `long:"amqp-exchange-no-wait" env:"AMQP_EXCHANGE_NO_WAIT" description:"See RabbitMQ docs"`
	QueueName          string `long:"amqp-queue-name" env:"AMQP_QUEUE_NAME" description:"See RabbitMQ docs"`
	QueueTag           string `long:"amqp-queue-tag" env:"AMQP_QUEUE_TAG" description:"See RabbitMQ docs"`
	QueueDurable       bool   `long:"amqp-queue-durable" env:"AMQP_QUEUE_DURABLE" description:"See RabbitMQ docs"`
	QueueAutoDelete    bool   `long:"amqp-queue-auto-delete" env:"AMQP_QUEUE_AUTO_DELETE" description:"See RabbitMQ docs"`
	QueueExclusive     bool   `long:"amqp-queue-exclusive" env:"AMQP_QUEUE_EXCLUSIVE" description:"See RabbitMQ docs"`
	QueueNoWait        bool   `long:"amqp-queue-no-wait" env:"AMQP_QUEUE_NO_WAIT" description:"See RabbitMQ docs"`
	QosPrefetchCount   int    `long:"amqp-qos-prefetch-count" env:"AMQP_QOS_PREFETCH_COUNT" default:"0" description:"See RabbitMQ docs"`
	QosPrefetchSize    int    `long:"amqp-qos-prefetch-size" env:"AMQP_QOS_PREFETCH_SIZE" default:"0" description:"See RabbitMQ docs"`
	QosGlobal          bool   `long:"amqp-qos-global" env:"AMQP_QOS_GLOBAL" description:"See RabbitMQ docs"`
	Mandatory          bool   `long:"amqp-mandatory" env:"AMQP_MANDATORY" description:"See RabbitMQ docs"`
	Immediate          bool   `long:"amqp-immediate" env:"AMQP_IMMEDIATE" description:"See RabbitMQ docs"`
}

func (r RabbitOptions) GetConsumerSettings() *amqp.ConsumerSettings {
	return &amqp.ConsumerSettings{
		Uri: r.Uri,
		Exchange: amqp.ExchangeSettings{
			Name:       r.ExchangeName,
			Type:       r.ExchangeType,
			Durable:    r.ExchangeDurable,
			AutoDelete: r.ExchangeAutoDelete,
			Internal:   r.ExchangeInternal,
			NoWait:     r.ExchangeNoWait,
		},
		Queue: amqp.QueueSettings{
			Name:       r.QueueName,
			Tag:        r.QueueTag,
			Durable:    r.QueueDurable,
			AutoDelete: r.QueueAutoDelete,
			Exclusive:  r.QueueExclusive,
			NoWait:     r.QueueNoWait,
		},
		Qos: amqp.QosSettings{
			PrefetchCount: r.QosPrefetchCount,
			PrefetchSize:  r.QosPrefetchSize,
			Global:        r.QosGlobal,
		},
		Tag:       r.ConsumerTag,
		AutoAck:   r.ConsumerAutoAck,
		Exclusive: r.ConsumerExclusive,
		NoLocal:   r.ConsumerNoLocal,
		NoWait:    r.ConsumerNoWait,
		Requeue:   r.ConsumerRequeue,
	}
}

func (r RabbitOptions) GetProducerSettings() *amqp.ProducerSettings {
	return &amqp.ProducerSettings{
		Uri: r.Uri,
		Exchange: amqp.ExchangeSettings{
			Name:       r.ExchangeName,
			Type:       r.ExchangeType,
			Durable:    r.ExchangeDurable,
			AutoDelete: r.ExchangeAutoDelete,
			Internal:   r.ExchangeInternal,
			NoWait:     r.ExchangeNoWait,
		},
		Tag:       r.ConsumerTag,
		Mandatory: r.Mandatory,
		Immediate: r.Immediate,
	}
}
