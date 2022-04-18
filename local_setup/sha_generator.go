package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/Shopify/sarama"
)

// Exit codes
const (
	// ExitStatusOK - the tool finished with success
	ExitStatusOK = iota
	// ExitStatusKafkaBrokerError - kafka broker connection establishment errors
	ExitStatusKafkaBrokerError
	// ExitStatusKafkaProducerError - kafka event production failures
	ExitStatusKafkaProducerError
)

// CliFlags holds all the allowed command line arguments and flags.
type CliFlags struct {
	NumMessages int
	Produce     bool
	KafkaBroker string
	KafkaTopic  string
}

// JSONContent represents any JSON object as key-value mapping
type JSONContent map[string]*json.RawMessage

// KafkaMessage is the structure of JSON messages produced
type KafkaMessage struct {
	ImageCount int         `json:"imageCount"`
	Digests    JSONContent `json:"images"`
}

func generateSHA256(count int) (sha []string) {
	sha = make([]string, count)
	secret := []byte("mysecret")
	for i := 0; i < count; i++ {
		// Create a new HMAC by defining the hash type and the key (as byte array)
		h := hmac.New(sha256.New, secret)

		// Write Data to it
		h.Write([]byte(time.Now().String()))

		// Get result and encode as hexadecimal string
		sha[i] = "sha256:" + hex.EncodeToString(h.Sum(nil))
	}
	return
}

func produce(shas []string, broker, topic string) {
	producer, err := sarama.NewSyncProducer([]string{broker}, nil)
	if err != nil {
		fmt.Printf("couldn't connect to Kafka broker %v\n", broker)
		os.Exit(ExitStatusKafkaBrokerError)
	}

	kafkaMsg := KafkaMessage{
		ImageCount: len(shas),
	}

	images := make(JSONContent, len(shas))

	//empty content, as long as it is a valid JSON object
	content := json.RawMessage("{}")
	for _, msg := range shas {
		images[msg] = &content
	}

	kafkaMsg.Digests = images

	fmt.Println(kafkaMsg)

	jsonBytes, err := json.Marshal(kafkaMsg)
	if err != nil {
		fmt.Println("couldn't turn Kafka message into valid JSON")
		fmt.Printf("error: %v\n", err)
		fmt.Printf("kafka message: %v\n", kafkaMsg)
		os.Exit(ExitStatusKafkaProducerError)
	}

	producerMsg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(jsonBytes),
	}

	partitionID, offset, err := producer.SendMessage(producerMsg)
	if err != nil {
		fmt.Println("failed to produce message to Kafka")
	} else {
		fmt.Printf("message sent to partition %d at offset %d\n", partitionID, offset)
	}
}

func main() {
	if len(os.Args) > 1 {
		var flags CliFlags

		flag.IntVar(&flags.NumMessages, "num-messages", 1, "number of SHA256 messages to generate")
		flag.BoolVar(&flags.Produce, "produce", false, "send generated SHAs to configured Kafka topic")
		flag.StringVar(&flags.KafkaBroker, "kafka-broker", "localhost:9092", "Kafka broker in the <host>:<port> format")
		flag.StringVar(&flags.KafkaTopic, "kafka-topic", "test_sha", "Kafka topic for producer")
		flag.Parse()

		shas := generateSHA256(flags.NumMessages)

		if !flags.Produce {
			fmt.Println(shas)
			os.Exit(ExitStatusOK)
		}

		produce(shas, flags.KafkaBroker, flags.KafkaTopic)
		os.Exit(ExitStatusOK)
	}

	shas := generateSHA256(1)
	fmt.Println(shas)
}
