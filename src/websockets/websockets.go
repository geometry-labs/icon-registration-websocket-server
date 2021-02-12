package websockets

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"

	"kafka-websocket-server/registration"
)

type KafkaWebsocketServer struct {
	TopicChan chan *kafka.Message
	Port      string
}

func (ws *KafkaWebsocketServer) ListenAndServe() {

	http.HandleFunc("/", ws.readAndFilterKafkaTopic)

	log.Fatal(http.ListenAndServe(":"+ws.Port, nil))
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (ws *KafkaWebsocketServer) readAndFilterKafkaTopic(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()

	// Close signal from client
	client_close_sig := make(chan bool)

	// Read for registrations
	broadcaster_ids := make([]registration.BroadcasterID, 0, 0)
	go func() {
		for {
			_, msg_raw, err := c.ReadMessage()
			if err != nil {
				client_close_sig <- true
				break
			}

			broadcaster_id, err := registration.RegisterBroadcaster(msg_raw)
			if err != nil {
				_ = c.WriteMessage(websocket.TextMessage, []byte(`{"error": "failed to register"}`))
			}

			broadcaster_ids = append(broadcaster_ids, broadcaster_id)
			_ = c.WriteMessage(websocket.TextMessage, []byte(`{"error": ""}`))

			defer func() {
				log.Printf("Unregister: %s\n", string(broadcaster_id))
			}()
		}
	}()

	// Write to websocket connection
	go func() {
		for {
			// Read
			msg := <-ws.TopicChan

			// TODO filter

			// Broadcast
			err = c.WriteMessage(websocket.TextMessage, msg.Value)
			if err != nil {
				break
			}
		}
	}()

	for {
		// check for client close
		select {
		case _ = <-client_close_sig:
			return
		default:
			continue
		}
	}
}
