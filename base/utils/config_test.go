package utils

import (
	"testing"

	"github.com/IBM/sarama"

	clowder "github.com/redhatinsights/app-common-go/pkg/api/v1"
	"github.com/stretchr/testify/assert"
)

func createTestBroker(saslMechanism string) clowder.BrokerConfig {
	nilStr := "nil"
	authType := clowder.BrokerConfigAuthtypeSasl

	return clowder.BrokerConfig{
		Authtype: &authType,
		Cacert:   nil,
		Hostname: "",
		Port:     nil,
		Sasl: &clowder.KafkaSASLConfig{
			Password:         &nilStr,
			SaslMechanism:    &saslMechanism,
			SecurityProtocol: &nilStr,
			Username:         &nilStr,
		},
	}
}

func TestSetKafkaSSLConfig(t *testing.T) {
	for _, s := range []string{sarama.SASLTypePlaintext, sarama.SASLTypeSCRAMSHA256, sarama.SASLTypeSCRAMSHA512} {
		Cfg.KafkaBroker = createTestBroker(s)
		assert.Nil(t, SetKafkaSSLConfig(sarama.NewConfig()))
	}
}

func TestSetKafkaSSLConfigNil(t *testing.T) {
	err := SetKafkaSSLConfig(nil)
	assert.Equal(t, "sarama config is required", err.Error())
}

func TestSetKafkaSSLConfigNilBroker(t *testing.T) {
	Cfg.KafkaBroker = createTestBroker("")
	Cfg.KafkaBroker.Sasl = nil
	err := SetKafkaSSLConfig(sarama.NewConfig())
	assert.Equal(t, "sasl config on kafka broker does not exist", err.Error())
}

func TestSetKafkaSSLConfigInvalidMechanism(t *testing.T) {
	Cfg.KafkaBroker = createTestBroker("")
	Cfg.KafkaBroker.Sasl.SaslMechanism = nil
	err := SetKafkaSSLConfig(sarama.NewConfig())
	assert.Equal(t, "sasl mechanism not specified", err.Error())
}

func TestSetKafkaSSLConfigUnknownMechanism(t *testing.T) {
	Cfg.KafkaBroker = createTestBroker("unknown")
	err := SetKafkaSSLConfig(sarama.NewConfig())
	assert.Equal(t, "unknown sasl mechanism: unknown", err.Error())
}

func TestSetKafkaSSLConfigInvalidUsr(t *testing.T) {
	Cfg.KafkaBroker = createTestBroker(sarama.SASLTypePlaintext)
	Cfg.KafkaBroker.Sasl.Username = nil
	err := SetKafkaSSLConfig(sarama.NewConfig())
	assert.Equal(t, "sasl username not specified", err.Error())
}

func TestSetKafkaSSLConfigInvalidPswd(t *testing.T) {
	Cfg.KafkaBroker = createTestBroker(sarama.SASLTypePlaintext)
	Cfg.KafkaBroker.Sasl.Password = nil
	err := SetKafkaSSLConfig(sarama.NewConfig())
	assert.Equal(t, "sasl password not specified", err.Error())
}

// It does not test cert correctness.
func TestSetKafkaTLSConfig(t *testing.T) {
	cert := "-----BEGIN CERTIFICATE-----\nMIIIxlc2VhcVAE3v7tfRnW9\n-----END CERTIFICATE-----"
	Cfg.KafkaBroker.Cacert = &cert

	scfg := sarama.NewConfig()

	assert.Nil(t, SetKafkaTLSConfig(scfg))
	assert.True(t, true, scfg.Net.TLS.Config.RootCAs)
}
