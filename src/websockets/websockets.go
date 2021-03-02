package websockets

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"

	"kafka-websocket-server/registration"
)

type KafkaWebsocketServer struct {
	Broadcaster *TopicBroadcaster
	Port        string
	Prefix      string
}

func (ws *KafkaWebsocketServer) ListenAndServe() {

	endpoint_path := ws.Prefix + "/"

	http.HandleFunc(endpoint_path, ws.readAndFilterKafkaTopic)

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

	// Add broadcaster
	topic_chan := make(chan *kafka.Message)
	id := ws.Broadcaster.AddWebsocketChannel(topic_chan)
	defer func() {
		// Remove broadcaster
		ws.Broadcaster.RemoveWebsocketChannel(id)
	}()

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
			registered_msg := fmt.Sprintf(`{"broadcaster_id": "%s"}`, string(broadcaster_id))
			_ = c.WriteMessage(websocket.TextMessage, []byte(registered_msg))

			defer func() {
				err := registration.UnregisterBroadcaster(broadcaster_id)
				log.Printf("UNREG: %s", string(broadcaster_id))
				if err != nil {
					log.Printf("Error unregistering: %s", err.Error())
				}
			}()
		}
	}()

	// Write to websocket connection
	go func() {
		for {
			// Read
			msg := <-topic_chan

			var broadcaster_ids_key []string
			_ = json.Unmarshal(msg.Key, &broadcaster_ids_key)

			// Compare local registered ids to msg registered ids
			wrote_flag := false
			for _, b := range broadcaster_ids {
				if wrote_flag == true {
					break
				}
				for _, bk := range broadcaster_ids_key {
					if string(b) == bk {
						// Broadcast
						_ = c.WriteMessage(websocket.TextMessage, msg.Value)

						// Raise flag to cancel duplicate messages
						wrote_flag = true
						break
					}
				}
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
