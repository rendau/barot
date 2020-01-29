package mq

import (
	"context"
	"github.com/rendau/barot/internal/interfaces"
	"github.com/streadway/amqp"
	"time"
)

const (
	connectionWaitTimeout = 30 * time.Second
	queueName             = "banner_notify"
)

// Rmq - is type for rabbit-mq client
type Rmq struct {
	lg  interfaces.Logger
	con *amqp.Connection
	ch  *amqp.Channel
}

// NewRmq - creates new Rmq instance
func NewRmq(dsn string, lg interfaces.Logger) (*Rmq, error) {
	var err error

	res := &Rmq{
		lg: lg,
	}

	connectionCtx, _ := context.WithTimeout(context.Background(), connectionWaitTimeout)
	res.con, err = res.connectionWait(connectionCtx, dsn)
	if err != nil {
		res.lg.Errorw("Fail to connect", err)
		return nil, err
	}

	res.ch, err = res.con.Channel()
	if err != nil {
		res.lg.Errorw("Fail to open channel", err)
		return nil, err
	}

	_, err = res.ch.QueueDeclare(
		queueName,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		res.lg.Errorw("Fail to create queue", err)
		return nil, err
	}

	return res, nil
}

func (r *Rmq) connectionWait(ctx context.Context, dsn string) (*amqp.Connection, error) {
	var err error
	var res *amqp.Connection

	for {
		res, err = amqp.Dial(dsn)
		if err == nil || ctx.Err() != nil {
			break
		}
		time.Sleep(time.Second)
	}

	return res, err
}

//// PublishEventNotification - publishes event to mq
//func (r *Rmq) PublishEventNotification(event *entities.Event) error {
//	eventBytes, err := json.Marshal(event)
//	if err != nil {
//		return err
//	}
//
//	err = r.ch.Publish(
//		"",
//		queueName,
//		false,
//		false,
//		amqp.Publishing{
//			ContentType: "application/json",
//			Body:        eventBytes,
//		},
//	)
//
//	return err
//}

// Stop - stops mq
func (r *Rmq) Stop() {
	err := r.ch.Close()
	if err != nil {
		r.lg.Errorw("Fail to close channel", err)
	}
	err = r.con.Close()
	if err != nil {
		r.lg.Errorw("Fail to close connection", err)
	}
}
