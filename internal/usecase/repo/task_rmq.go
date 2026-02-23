package repo

import (
	"fmt"

	"github.com/arunima10a/task-manager/pkg/rabbitmq"
	amqp "github.com/rabbitmq/amqp091-go"
)

type TaskRMQ struct {
	rmq *rabbitmq.RabbitMQ
}

func NewTaskRMQ(rmq *rabbitmq.RabbitMQ) *TaskRMQ {
	return &TaskRMQ{rmq: rmq}
}

func (r *TaskRMQ) PublishTaskCreated(taskID int, description string) error {
	ch, err := r.rmq.Conn.Channel()
	if err != nil {
		return err
	}
	defer func() { _ = ch.Close() }()

	q, _ := ch.QueueDeclare("task_created", true, false, false, false, nil)

	body := fmt.Sprintf("%d|%s", taskID, description)
	return ch.Publish("", q.Name, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(body),
	})
}
