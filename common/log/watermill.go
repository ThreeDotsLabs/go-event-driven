package log

import (
	"github.com/ThreeDotsLabs/watermill/message"
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
