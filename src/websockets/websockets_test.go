package websockets

import (
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

func TestKafkaWebsocketServer(t *testing.T) {

	topic_chan := make(chan *kafka.Message)

	websocket_server := KafkaWebsocketServer{
		topic_chan,
		"8080",
		"",
	}

	// Start websocket server
	go websocket_server.ListenAndServe()

	// Start mock channel data
	go func() {
		for {
			msg := &(kafka.Message{})
			msg.Value = []byte("Test Data")

			topic_chan <- msg

			time.Sleep(1 * time.Second)
		}
	}()

	// Validate message
	websocket_client, _, err := websocket.DefaultDialer.Dial("ws://localhost:8080/data", nil)
	if err != nil {
		t.Logf("Failed to connect to KafkaWebsocketServer")
		t.Fail()
	}
	defer websocket_client.Close()

	_, message, err := websocket_client.ReadMessage()
	if err != nil {
		t.Logf("Failed to read websocket")
		t.Fail()
	}

	if string(message) != "Test Data" {
		t.Logf("Failed to validate data")
		t.Fail()
	}
}
