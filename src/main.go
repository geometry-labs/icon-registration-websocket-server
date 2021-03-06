package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"

	"kafka-websocket-server/consumer"
	"kafka-websocket-server/registration"
	"kafka-websocket-server/websockets"
)

func main() {

	output_topic_env := os.Getenv("ICON_REGISTRATION_WEBSOCKET_OUTPUT_TOPIC")
	broker_url_env := os.Getenv("ICON_REGISTRATION_WEBSOCKET_BROKER_URL")
	registration_url_env := os.Getenv("ICON_REGISTRATION_WEBSOCKET_REGISTRATION_URL")
	port_env := os.Getenv("ICON_REGISTRATION_WEBSOCKET_PORT")
	prefix_env := os.Getenv("ICON_REGISTRATION_WEBSOCKET_PREFIX")

	if broker_url_env == "" {
		log.Println("ERROR: required enviroment variable missing: ICON_REGISTRATION_WEBSOCKET_BROKER_URL")
		os.Exit(1)
	}
	if registration_url_env == "" {
		log.Println("ERROR: required enviroment variable missing: ICON_REGISTRATION_WEBSOCKET_REGISTRATION_URL")
		os.Exit(1)
	}
	if output_topic_env == "" {
		output_topic_env = "outputs"
	}
	if port_env == "" {
		port_env = "3000"
	}
	if prefix_env == "" {
		prefix_env = ""
	}

	// Set registration url
	registration.SetRegistrationURL(registration_url_env)

	// Create channel
	output_topic_chan := make(chan *kafka.Message)

	// Create consumer
	kafka_consumer := consumer.KafkaTopicConsumer{
		output_topic_env,
		output_topic_chan,
		broker_url_env,
	}

	// Start consumer
	go kafka_consumer.ConsumeAndBroadcastTopics()
	log.Println("Kafka consumer started...")

	// Create broadcaster
	broadcaster := &websockets.TopicBroadcaster{
		output_topic_chan,
		make(map[websockets.BroadcasterID]chan *kafka.Message),
	}

	// Start broadcaster
	go broadcaster.Broadcast()
	log.Printf("Topic broadcaster started...")

	// Create server
	websocket_server := websockets.KafkaWebsocketServer{
		broadcaster,
		port_env,
		prefix_env,
	}

	// Start server
	go websocket_server.ListenAndServe()
	log.Printf("Server listening on port :%s%s/...", port_env, prefix_env)

	// Listen for close sig
	sigCh := make(chan os.Signal)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM, syscall.SIGKILL)

	// Keep main thread alive
	for sig := range sigCh {
		log.Printf("Stopping registration websocket server...%s", sig.String())
	}
}
