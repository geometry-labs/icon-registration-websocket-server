package websockets

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
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

	// close signal
	client_close_sig := make(chan bool)

	// Read for registrations
	// broadcaster_ids := make([]string, 0, 0)
	go func() {
		for {
			_, msg_raw, err := c.ReadMessage()
			if err != nil {
				client_close_sig <- true
				break
			}

			log.Println(string(msg_raw))

		}
	}()

	for {
		// Read
		msg := <-ws.TopicChan

		// Broadcast
		err = c.WriteMessage(websocket.TextMessage, msg.Value)
		if err != nil {
			break
		}

		// check for client close
		select {
		case _ = <-client_close_sig:
			break
		default:
			continue
		}
	}
}
