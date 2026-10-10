// File: subscribers/subscriber.go
package subscribers

import (
	"context"
	"fmt"

	"scav/infra"
	"scav/infra/logger"
	"scav/infra/mq"
)

// registration describes a single subscription to a subject.
type registration struct {
	Subject string
	Handler mq.MessageHandler
}

// RegisterAll registers all application MQ subscriptions.
func RegisterAll(
	ctx context.Context,
	app *infra.Deps,
) error {
	if ctx == nil {
		return fmt.Errorf("context is nil")
	}
	if app == nil {
		return fmt.Errorf("app is nil")
	}
	if app.MQ == nil {
		return fmt.Errorf("MQ is not initialized")
	}

	registrations := []registration{
		{
			Subject: "chat.created",
			Handler: handleChatCreated,
		},
		{
			Subject: "chat.message.created",
			Handler: handleChatMessageCreated,
		},
		{
			Subject: "user.created",
			Handler: handleUserCreated,
		},
		{
			Subject: "notification.created",
			Handler: handleNotificationCreated,
		},
	}

	for _, reg := range registrations {
		if reg.Subject == "" {
			return fmt.Errorf("MQ registration has empty subject")
		}
		if reg.Handler == nil {
			return fmt.Errorf("MQ registration %q has nil handler", reg.Subject)
		}

		sub, err := app.MQ.Subscribe(ctx, reg.Subject, reg.Handler)
		if err != nil {
			return fmt.Errorf("subscribe subject=%q: %w", reg.Subject, err)
		}

		// Cleanly unsubscribe when application lifecycle context terminates
		go func(s mq.Subscription, subject string) {
			<-ctx.Done()
			if err := s.Unsubscribe(); err != nil {
				logger.L.Sugar().Warnw("failed to unsubscribe during context teardown", "subject", subject, "error", err)
			}
		}(sub, reg.Subject)

		logger.L.Sugar().Infow(
			"MQ subscriber registered",
			"subject", reg.Subject,
		)
	}

	logger.L.Sugar().Infow(
		"all MQ subscribers registered",
		"count", len(registrations),
	)

	return nil
}

func handleChatCreated(ctx context.Context, msg mq.Message) error {
	return processEvent(ctx, msg, "chat.created")
}

func handleChatMessageCreated(ctx context.Context, msg mq.Message) error {
	return processEvent(ctx, msg, "chat.message.created")
}

func handleUserCreated(ctx context.Context, msg mq.Message) error {
	return processEvent(ctx, msg, "user.created")
}

func handleNotificationCreated(ctx context.Context, msg mq.Message) error {
	return processEvent(ctx, msg, "notification.created")
}

func processEvent(
	ctx context.Context,
	msg mq.Message,
	subject string,
) error {
	if ctx == nil {
		return fmt.Errorf("context is nil")
	}

	event, err := mq.UnpackEnvelope(msg.Data)
	if err != nil {
		return fmt.Errorf("unpack %s event: %w", subject, err)
	}

	// Re-inject context metadata for downstream propagation
	if event.TraceID != "" {
		ctx = mq.WithTraceID(ctx, event.TraceID)
	}
	if event.Source != "" {
		ctx = mq.WithServiceName(ctx, event.Source)
	}

	logger.L.Sugar().Infow(
		"processing MQ event",
		"subject", subject,
		"event_id", event.ID,
		"trace_id", event.TraceID,
		"source", event.Source,
	)

	return nil
}
