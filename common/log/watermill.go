package log

import (
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/sirupsen/logrus"
)

const correlationIDMessageMetadataKey = "correlation_id"

type CorrelationPublisherDecorator struct {
	message.Publisher
}

func (c CorrelationPublisherDecorator) Publish(topic string, messages ...*message.Message) error {
	for i := range messages {
		// if correlation_id is already set, let's not override
		if messages[i].Metadata.Get(correlationIDMessageMetadataKey) != "" {
			continue
		}

		// correlation_id as const
		messages[i].Metadata.Set(correlationIDMessageMetadataKey, CorrelationIDFromContext(messages[i].Context()))
	}

	return c.Publisher.Publish(topic, messages...)
}

type WatermillLogrusAdapter struct {
	Log *logrus.Entry
}

func NewWatermill(log *logrus.Entry) *WatermillLogrusAdapter {
	return &WatermillLogrusAdapter{Log: log}
}

func (w WatermillLogrusAdapter) Error(msg string, err error, fields watermill.LogFields) {
	w.Log.WithError(err).WithFields(logrus.Fields(fields)).Error(msg)
}

func (w WatermillLogrusAdapter) Info(msg string, fields watermill.LogFields) {
	// Watermill info logs are too verbose
	w.Log.WithFields(logrus.Fields(fields)).Debug(msg)
}

func (w WatermillLogrusAdapter) Debug(msg string, fields watermill.LogFields) {
	w.Log.WithFields(logrus.Fields(fields)).Debug(msg)
}

func (w WatermillLogrusAdapter) Trace(msg string, fields watermill.LogFields) {
	w.Log.WithFields(logrus.Fields(fields)).Trace(msg)
}

func (w WatermillLogrusAdapter) With(fields watermill.LogFields) watermill.LoggerAdapter {
	return WatermillLogrusAdapter{w.Log.WithFields(logrus.Fields(fields))}
}
