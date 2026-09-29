package mq

import (
	"context"
	"fmt"

	"github.com/nats-io/nats.go"
)

type jetStreamMQ struct {
	js nats.JetStreamContext
}

type jsSubscription struct {
	sub *nats.Subscription
}

func (s *jsSubscription) Unsubscribe() error {
	if s.sub == nil {
		return nil
	}
	return s.sub.Unsubscribe()
}

// NewJetStreamMQ creates an MQ implementation backed by NATS JetStream.
func NewJetStreamMQ(jsctx nats.JetStreamContext) MQ {
	return &jetStreamMQ{
		js: jsctx,
	}
}

func (j *jetStreamMQ) Publish(ctx context.Context, subject string, data []byte) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	_, err := j.js.Publish(subject, data, nats.Context(ctx))
	if err != nil {
		return fmt.Errorf("js publish: %w", err)
	}
	return nil
}

func (j *jetStreamMQ) Ping(ctx context.Context) error {
	// Query JetStream account information to verify connectivity
	_, err := j.js.AccountInfo(nats.Context(ctx))
	if err != nil {
		return fmt.Errorf("js ping failed: %w", err)
	}
	return nil
}

func (j *jetStreamMQ) Subscribe(ctx context.Context, subject string, handler MessageHandler) (Subscription, error) {
	cb := j.makeCallback(ctx, handler)

	sub, err := j.js.Subscribe(subject, cb, nats.ManualAck())
	if err != nil {
		return nil, fmt.Errorf("js subscribe: %w", err)
	}

	return &jsSubscription{sub: sub}, nil
}

func (j *jetStreamMQ) QueueSubscribe(ctx context.Context, subject, queue string, handler MessageHandler) (Subscription, error) {
	cb := j.makeCallback(ctx, handler)

	sub, err := j.js.QueueSubscribe(subject, queue, cb, nats.ManualAck())
	if err != nil {
		return nil, fmt.Errorf("js queue subscribe: %w", err)
	}

	return &jsSubscription{sub: sub}, nil
}

// makeCallback binds your custom MessageHandler signature to NATS's msg.Ack() / msg.Nak() behavior.
func (j *jetStreamMQ) makeCallback(parentCtx context.Context, handler MessageHandler) nats.MsgHandler {
	return func(msg *nats.Msg) {
		reqCtx := parentCtx
		if reqCtx == nil {
			reqCtx = context.Background()
		}

		m := Message{
			Subject: msg.Subject,
			Data:    msg.Data,
		}

		err := handler(reqCtx, m)
		if err != nil {
			_ = msg.Nak()
			return
		}

		_ = msg.Ack()
	}
}
