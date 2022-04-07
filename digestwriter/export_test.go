package digestwriter

import (
	"github.com/Shopify/sarama"
)

// Export for testing
//
// This source file contains name aliases of all package-private functions
// that need to be called from unit tests. Aliases should start with uppercase
// letter because unit tests belong to different package.
//
// Please look into the following blogpost:
// https://medium.com/@robiplus/golang-trick-export-for-test-aa16cbd7b8cd
// to see why this trick is needed.

var (
	// functions from consumer.go source file
	ExtractDigestsFromMessage = extractDigestsFromMessage
	ParseMessage              = parseMessage
	SetupLogger               = setupLogger
)

func HandleKafkaMessage(c *KafkaConsumer, msg *sarama.ConsumerMessage) {
	c.handleMessage(msg)
}
