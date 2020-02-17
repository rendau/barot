package rmq

import (
	"context"
	"encoding/json"
	"time"

	"github.com/rendau/barot/internal/domain/entities"
	"github.com/rendau/barot/internal/interfaces"
	"github.com/streadway/amqp"
)

const (
	connectionWaitTimeout = 30 * time.Second
	queueName             = "banner_notify"
)

// RabbitMQ - is type for rabbit-mq client
type RabbitMQ struct {
	lg  interfaces.Logger
	con *amqp.Connection
	ch  *amqp.Channel
}

// NewRabbitMQ - creates new RabbitMQ instance
func NewRabbitMQ(dsn string, lg interfaces.Logger) (*RabbitMQ, error) {
	var err error

	res := &RabbitMQ{
		lg: lg,
	}

	connectionCtx, connectionCtxCancel := context.WithTimeout(context.Background(), connectionWaitTimeout)
	defer connectionCtxCancel()

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

func (q *RabbitMQ) connectionWait(ctx context.Context, dsn string) (*amqp.Connection, error) {
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

// PublishBannerEvent - publishes event to mq
func (q *RabbitMQ) PublishBannerEvent(event *entities.BannerEvent) error {
	eventBytes, err := json.Marshal(struct {
		Type      string    `json:"type"`
		BannerID  int64     `json:"banner_id"`
		SlotID    int64     `json:"slot_id"`
		UsrTypeID int64     `json:"usr_type_id"`
		DateTime  time.Time `json:"datetime"`
	}{
		event.Type.String(),
		event.BannerID,
		event.SlotID,
		event.UsrTypeID,
		event.DateTime,
	})
	if err != nil {
		q.lg.Errorw("Fail to encode to json", err)
		return err
	}

	err = q.ch.Publish(
		"",
		queueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        eventBytes,
		},
	)

	return err
}

// Stop - stops mq
func (q *RabbitMQ) Stop() {
	err := q.ch.Close()
	if err != nil {
		q.lg.Errorw("Fail to close channel", err)
	}

	err = q.con.Close()
	if err != nil {
		q.lg.Errorw("Fail to close connection", err)
	}
}
