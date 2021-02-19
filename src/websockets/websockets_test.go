package websockets

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"

	"kafka-websocket-server/registration"
)

func TestRegistrationWebsocketServer(t *testing.T) {

	topic_chan := make(chan *kafka.Message)

	websocket_server := KafkaWebsocketServer{
		topic_chan,
		"8080",
		"",
	}

	// Start websocket server
	go websocket_server.ListenAndServe()

	// Set Register URL
	registration_url_env := os.Getenv("ICON_REGISTRATION_WEBSOCKET_REGISTRATION_URL")
	registration.SetRegistrationURL(registration_url_env)

	// Test json config
	register_json := `
	{
					"connection_type": "ws",
					"endpoint": "wss://test",
					"transaction_events": [
							{
									"to_address": "cx0000000000000000000000000000000000000000"
							}
					]
	}
	`

	// Open websocket
	websocket_client, _, err := websocket.DefaultDialer.Dial("ws://localhost:8080/", nil)
	if err != nil {
		t.Logf("Failed to connect to KafkaWebsocketServer")
	}
	defer websocket_client.Close()

	// Send registration
	err = websocket_client.WriteMessage(websocket.TextMessage, []byte(register_json))
	if err != nil {
		t.Logf("Failed to write to websocket")
		t.Fail()
	}

	// Read broadcaster_id
	_, msg_raw, err := websocket_client.ReadMessage()
	if err != nil {
		t.Logf("Failed to read websocket")
		t.Fail()
	}

	msg_json := make(map[string]interface{})
	err = json.Unmarshal(msg_raw, &msg_json)
	if err != nil {
		t.Logf("Failed to parse broadcaster_id")
		t.Fail()
	}

	broadcaster_id := msg_json["broadcaster_id"].(string)

	// Start mock channel data
	topic_value := fmt.Sprintf(`{"test_val": "%s"}`, broadcaster_id)
	topic_key := fmt.Sprintf(`["%s", "other-broadcaster-id"]`, broadcaster_id)
	go func() {
		for {
			msg := &(kafka.Message{})
			msg.Value = []byte(topic_value)
			msg.Key = []byte(topic_key)

			topic_chan <- msg

			time.Sleep(1 * time.Second)
		}
	}()

	_, message, err := websocket_client.ReadMessage()
	if err != nil {
		t.Logf("Failed to read websocket")
		t.Fail()
	}

	if string(message) != topic_value {
		t.Logf("Failed to validate data")
		t.Fail()
	}

	// Pass
	return
}
