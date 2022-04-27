package ginCore

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
)

type Producer interface {
	MsgContent() string
}

// 定义接收者接口
type Receiver interface {
	Consumer([]byte) error
}

// 消息队列配置
type MQConfig struct {
	Type     string `json:"type"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Vhost    string `json:"vhost"`
}

type MQService struct {
	Config     *MQConfig
	Connection *amqp.Connection
	Channel    *amqp.Channel
	Done       chan error
}

type HandlerMsgFunc func(string) error

func NewMQService(c *MQConfig) (*MQService, error) {
	service := &MQService{Config: c}
	return service.Connect()
}

func (s *MQService) Connect() (*MQService, error) {
	var err error
	//RabbitUrl := fmt.Sprintf("amqp://%s:%s@%s:%d/", s.Config.User, s.Config.Password, s.Config.Host, s.Config.Port)
	mqConn, err := amqp.Dial(amqp.URI{Scheme: "amqp", Host: s.Config.Host, Port: s.Config.Port, Username: s.Config.User, Password: s.Config.Password, Vhost: "/"}.String())
	if err != nil {
		fmt.Printf("MQ打开链接失败:%s \n", err)
		return nil, err
	}
	mqChan, err := mqConn.Channel()
	if err != nil {
		fmt.Printf("MQ打开管道失败:%s \n", err)
		return nil, err
	}
	s.Connection = mqConn
	s.Channel = mqChan
	return s, nil
}

func (s *MQService) ExchangeDeclare(exchangeName, exchangeType string, durable, autoDelete, internal, noWait bool) (*MQService, error) {
	err := s.Channel.ExchangeDeclare(exchangeName, exchangeType, durable, autoDelete, internal, noWait, nil)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return s, nil
}

func (s *MQService) Publish(exchangeName, routingKey, body string, reliable bool) error {
	if reliable {
		if err := s.Channel.Confirm(false); err != nil {
			fmt.Println(err)
			return err
		}
		confirms := s.Channel.NotifyPublish(make(chan amqp.Confirmation, 1))
		defer s.confirmOne(confirms)
	}
	err := s.Channel.Publish(exchangeName, routingKey, false, false, amqp.Publishing{
		Headers:         amqp.Table{},
		ContentType:     "text/plain",
		ContentEncoding: "",
		Body:            []byte(body),
		DeliveryMode:    amqp.Transient, // 1=non-persistent, 2=persistent
		Priority:        0,              // 0-9
		// a bunch of application/implementation-specific fields
	})
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (s *MQService) Customer(exchangeName, queueName, key string, function HandlerMsgFunc) error {
	go func() {
		fmt.Printf("closing: %s", <-s.Channel.NotifyClose(make(chan *amqp.Error)))
	}()

	queue, err := s.Channel.QueueDeclare(
		queueName, // name of the queue
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // noWait
		nil,       // arguments
	)
	if err != nil {
		fmt.Println(err)
		return err
	}
	if err = s.Channel.QueueBind(
		queue.Name,   // name of the queue
		key,          // bindingKey
		exchangeName, // sourceExchange
		false,        // noWait
		nil,          // arguments
	); err != nil {
		return err
	}

	deliveries, err := s.Channel.Consume(
		queue.Name, // name
		"",         // consumerTag,
		false,      // noAck
		false,      // exclusive
		false,      // noLocal
		false,      // noWait
		nil,        // arguments
	)
	if err != nil {
		return err
	}
	go s.handle(deliveries, s.Done, function)

	return nil
}

func (s *MQService) Close() {
	err := s.Channel.Close()
	if err != nil {
		fmt.Printf("MQ管道关闭失败:%s \n", err)
	}
	err = s.Connection.Close()
	if err != nil {
		fmt.Printf("MQ链接关闭失败:%s \n", err)
	}
}
func (s *MQService) confirmOne(confirms <-chan amqp.Confirmation) {
	log.Printf("waiting for confirmation of one publishing")

	if confirmed := <-confirms; confirmed.Ack {
		log.Printf("confirmed delivery with delivery tag: %d", confirmed.DeliveryTag)
	} else {
		log.Printf("failed delivery of delivery tag: %d", confirmed.DeliveryTag)
	}
}

func (s *MQService) handle(deliveries <-chan amqp.Delivery, done chan error, function HandlerMsgFunc) {
	for d := range deliveries {
		log.Printf(
			"got %dB delivery: [%v] %q",
			len(d.Body),
			d.DeliveryTag,
			d.Body,
		)
		err2 := function(string(d.Body))
		if err2 != nil {
			err := d.Ack(false)
			if err != nil {
				return
			}
		} else {
			err := d.Ack(true)
			if err != nil {
				return
			}
		}

	}
	log.Printf("handle: deliveries channel closed")
	done <- nil
}
