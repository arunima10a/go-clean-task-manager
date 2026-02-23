package rabbitmq

import (
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	Conn    *amqp.Connection
	connErr chan *amqp.Error
}

func New(url string) (*RabbitMQ, error) {
	r := &RabbitMQ{
		connErr: make(chan *amqp.Error),
	}

	if err := r.connect(url); err != nil {
		return nil, err
	}

	return r, nil
}

func (r *RabbitMQ) connect(url string) error {
	var err error
	for i := 0; i < 5; i++ {
		r.Conn, err = amqp.Dial(url)
		if err == nil {
			r.Conn.NotifyClose(r.connErr)

			go r.reconnect(url)
			return nil
		}
		time.Sleep(2 * time.Second)
	}
	return fmt.Errorf("rabbitmq - connect - dial: %w", err)
}
func (r *RabbitMQ) reconnect(url string) {
	for range r.connErr {
		log.Println("RabbitMQ connection lost. Retrying...")
		for {
			err := r.connect(url)
			if err == nil {
				log.Println("RabbitMQ reconnected successfully")
				return
			}
			time.Sleep(5 * time.Second)
		}
	}
}

func (r *RabbitMQ) Close() error {
	return r.Conn.Close()
}
