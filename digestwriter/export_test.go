package digestwriter

import "app/base/utils"

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
)

// kafka-related functions

// NewDummyConsumerWithProcessor has the same arguments as NewConsumer
// but returns a non-configured utils.KafkaConsumer and the DigestConsumer object for testing purposes
func NewDummyConsumerWithProcessor(storage Storage) (*utils.KafkaConsumer, *DigestConsumer) {
	processor := DigestConsumer{
		storage,
		0,
	}
	consumer := utils.KafkaConsumer{
		Processor: &processor,
	}
	return &consumer, &processor
}

// storage-related functions

func LinkDigestsToCluster(s *DBStorage, clusterStr string, clusterID, archID int64, digests []string) error {
	tx := s.connection.Begin()
	defer tx.Rollback()
	err := s.linkDigestsToCluster(tx, clusterStr, clusterID, archID, digests)
	if err != nil {
		return err
	}
	return tx.Commit().Error
}
