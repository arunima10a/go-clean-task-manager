package v1

import (
	"context"
	"strconv"
	"strings"

	"github.com/arunima10a/task-manager/internal/usecase"
	"github.com/arunima10a/task-manager/pkg/logger"
	amqp "github.com/rabbitmq/amqp091-go"
)

type taskRabbitMQ struct {
	t usecase.TaskUseCase
	l *logger.Logger
}

func NewTaskConsumer(t usecase.TaskUseCase, l *logger.Logger) *taskRabbitMQ {
	return &taskRabbitMQ{t: t, l: l}
}

func (r *taskRabbitMQ) Start(conn *amqp.Connection) {
	ch, err := conn.Channel()
	if err != nil {
		r.l.Error(err, "AMQP - taskRabbitMQ - Start - Channel")
		return
	}

	_, err = ch.QueueDeclare(
		"task_created",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		r.l.Error(err, "AMQP - taskRabbitMQ - Start - QueueDeclare")
		return
	}

	msgs, err := ch.Consume(
		"task_created",
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		r.l.Error(err, "AMQP - Consume failed")

	}

	go func() {
		for d := range msgs {
			r.l.Info("AMQP: Message received: %s", string(d.Body))

			parts := strings.Split(string(d.Body), "|")
			if len(parts) < 2 {
				continue
			}

			id, _ := strconv.Atoi(parts[0])

			err := r.t.EnrichTaskWithQuote(context.Background(), id)
			if err != nil {
				r.l.Error(err, "AMQP: Failed to enrich task")
			}
		}
	}()
}
