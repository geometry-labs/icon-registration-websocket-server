package websockets

import (
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

type BroadcasterID int

var LAST_BROADCASTER_ID BroadcasterID = 0

type TopicBroadcaster struct {

	// Input
	TopicChan chan *kafka.Message

	// Output
	WebsocketChans map[BroadcasterID]chan *kafka.Message
}

func (tb *TopicBroadcaster) AddWebsocketChannel(topic_chan chan *kafka.Message) BroadcasterID {
	id := LAST_BROADCASTER_ID
	LAST_BROADCASTER_ID++

	tb.WebsocketChans[id] = topic_chan

	return id
}

func (tb *TopicBroadcaster) RemoveWebsocketChannel(id BroadcasterID) {
	_, ok := tb.WebsocketChans[id]
	if ok {
		delete(tb.WebsocketChans, id)
	}
}

func (tb *TopicBroadcaster) Broadcast() {
	for {
		msg := <-tb.TopicChan

		for _, channel := range tb.WebsocketChans {
			select {
			case channel <- msg:
				continue
			default:
				continue
			}
		}
	}
}
