package main

import (
	"context"
	"github.com/segmentio/kafka-go"
	"log"
	"strconv"
)

const (
	topic         = "message-log"
	brokerAddress = "localhost:9092"
)

type kafkaMessage struct {
	id int
	content string
}

func createTopic() {
	c, _ := kafka.Dial("tcp", brokerAddress)
	kt := kafka.TopicConfig{
		Topic: topic,
		NumPartitions: 1,
		ReplicationFactor: 1,
	}
	err := c.CreateTopics(kt)
	if err != nil {
		panic("Could Not Create Topic " + err.Error())
	}
	log.Println("Topic Created Successfully")
}
// The Producer
func producer(ctx context.Context, message *kafkaMessage) {
	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{brokerAddress},
		Topic: topic,
	})
	err := w.WriteMessages(ctx, kafka.Message{
		Key: []byte(strconv.Itoa((message.id))),
		Value: []byte(message.content),
	})
	if err != nil {
		panic("Could not write message " + err.Error())
	}
	log.Println("Message Written")
}

// The Consumer
func consumer(ctx context.Context) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string {brokerAddress},
		Topic: topic,
	})
	for {
		msg, err := r.ReadMessage(ctx)
		if err != nil {
			panic("could not read message " + err.Error())
		}
		log.Printf("The Message Received is { %s }", string(msg.Value))
	}
}


func main() {
	ctx := context.Background()
	
	consumer(ctx)
}