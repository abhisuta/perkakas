package sse

import (
	"context"
	"github.com/Shopify/sarama"
	"github.com/kitabisa/perkakas/v2/queue/kafka"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewSseClient(t *testing.T) {
	client := NewSseClient("localhost:9092", kafka.WithClientID("unit-test"))
	assert.NotNil(t, client)
}

func TestClient_GetSetKafkaVersion(t *testing.T) {
	client := NewSseClient("localhost:9092", kafka.WithClientID("unit-test"))
	client.SetKafkaVersion(context.Background(), "2.5.0")
	version := client.GetKafkaVersion(context.Background())
	assert.Equal(t, "2.5.0", version)
}

func TestClient_PublishEvent(t *testing.T) {
	mockBroker := sarama.NewMockBroker(t, 1)
	mockBroker.SetHandlerByMap(map[string]sarama.MockResponse{
		"MetadataRequest": sarama.NewMockMetadataResponse(t).
			SetBroker(mockBroker.Addr(), mockBroker.BrokerID()),
		"ProduceRequest": sarama.NewMockProduceResponse(t),
	})

	client := NewSseClient(mockBroker.Addr(), kafka.WithClientID("katresnan"), kafka.WithRetryMax(5))
	data := map[string]interface{} {
		"name": "test",
	}

	for i := 0; i < 5; i++ {
		err := client.PublishEvent(context.Background(), "testing", "", data)
		assert.NoError(t, err)
	}
}
