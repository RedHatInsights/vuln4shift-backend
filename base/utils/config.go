package utils

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"strings"

	"github.com/Shopify/sarama"
	clowder "github.com/redhatinsights/app-common-go/pkg/api/v1"
)

var (
	Cfg = Config{}
)

type Config struct {
	// Database config
	DbHost          string
	DbPort          int
	DbName          string
	DbAdminUser     string
	DbAdminPassword string
	DbUser          string
	DbPassword      string

	// Port and metrics config
	PublicPort  int
	PrivatePort int
	MetricsPort int
	MetricsPath string

	// Common app config
	LoggingLevel              string
	APIRetries                int
	PrometheusPushGateway     string
	CloudWatchAccessKeyID     string
	CloudWatchSecretAccesskey string
	CloudWatchRegion          string
	CloudWatchLogGroup        string

	// DB admin config
	ArchiveDbWriterPass string
	PyxisGathererPass   string
	VmaasGathererPass   string
	CveAggregatorPass   string
	ManagerPass         string
	CleanerPass         string
	ExploitGathererPass string
	SchemaMigration     int

	// Manager config
	AmsEnabled      bool
	AmsAPIURL       string
	AmsAPIPagesize  int
	AmsClientID     string
	AmsClientSecret string

	// Digest writer config
	KafkaBroker              clowder.BrokerConfig
	KafkaBrokerAddress       string
	KafkaBrokerConsumerGroup string
	KafkaBrokerIncomingTopic string
	KafkaConsumerTimeout     string
	KafkaProducerTimeout     string
	KafkaPayloadTrackerTopic string
	PayloadTrackerEnabled    bool

	// VMaaS sync config
	VmaasBaseURL   string
	VmaasBatchSize int
	VmaasPageSize  int

	// Pyxis gatherer config
	PyxisBaseURL     string
	PyxisProfile     string
	ForceSync        bool
	SkipImageCveSync bool

	// Exploits sync config
	ExploitFileURL string
	GitToken       string
}

func init() {
	initConfig()
}

func initConfig() {
	// Database config
	if clowder.IsClowderEnabled() {
		Cfg.DbHost = clowder.LoadedConfig.Database.Hostname
		Cfg.DbPort = clowder.LoadedConfig.Database.Port
		Cfg.DbName = clowder.LoadedConfig.Database.Name
		Cfg.DbAdminUser = clowder.LoadedConfig.Database.AdminUsername
		Cfg.DbAdminPassword = clowder.LoadedConfig.Database.AdminPassword
	} else {
		Cfg.DbHost = GetEnv("POSTGRES_HOST", "unknown_host")
		Cfg.DbPort = GetEnv("POSTGRES_PORT", 0)
		Cfg.DbName = GetEnv("POSTGRES_DB", "unknown_database")
		Cfg.DbAdminUser = GetEnv("POSTGRES_ADMIN_USER", "unknown_admin")
		Cfg.DbAdminPassword = GetEnv("POSTGRES_ADMIN_PASSWORD", "unknown_admin_pwd")
	}
	Cfg.DbUser = GetEnv("POSTGRES_USER", "unknown_user")
	Cfg.DbPassword = GetEnv("POSTGRES_PASSWORD", "unknown_user_pwd")

	// Port and metrics config
	if clowder.IsClowderEnabled() {
		Cfg.PublicPort = *clowder.LoadedConfig.PublicPort
		Cfg.PrivatePort = *clowder.LoadedConfig.PrivatePort
		Cfg.MetricsPort = clowder.LoadedConfig.MetricsPort
		Cfg.MetricsPath = clowder.LoadedConfig.MetricsPath
	} else {
		Cfg.PublicPort = 8000
		Cfg.PrivatePort = 10000
		Cfg.MetricsPort = 9000
		Cfg.MetricsPath = "/metrics"
	}

	// Common app config
	Cfg.LoggingLevel = GetEnv("LOGGING_LEVEL", "INVALID")
	Cfg.APIRetries = GetEnv("API_RETRIES", 0)
	Cfg.PrometheusPushGateway = GetEnv("PROMETHEUS_PUSHGATEWAY", "pushgateway")
	if clowder.IsClowderEnabled() {
		cwCfg := clowder.LoadedConfig.Logging.Cloudwatch
		if cwCfg != nil {
			Cfg.CloudWatchAccessKeyID = cwCfg.AccessKeyId
			Cfg.CloudWatchSecretAccesskey = cwCfg.SecretAccessKey
			Cfg.CloudWatchRegion = cwCfg.Region
			Cfg.CloudWatchLogGroup = cwCfg.LogGroup
		}
	} else {
		Cfg.CloudWatchAccessKeyID = GetEnv("CLOUDWATCH_ACCESS_KEY_ID", "")
		Cfg.CloudWatchSecretAccesskey = GetEnv("CLOUDWATCH_SECRET_ACCESS_KEY", "")
		Cfg.CloudWatchRegion = GetEnv("CLOUDWATCH_REGION", "")
		Cfg.CloudWatchLogGroup = GetEnv("CLOUDWATCH_LOG_GROUP", "")
	}

	// DB admin config
	Cfg.ArchiveDbWriterPass = GetEnv("USER_ARCHIVE_DB_WRITER_PASS", "")
	Cfg.PyxisGathererPass = GetEnv("USER_PYXIS_GATHERER_PASS", "")
	Cfg.VmaasGathererPass = GetEnv("USER_VMAAS_GATHERER_PASS", "")
	Cfg.CveAggregatorPass = GetEnv("USER_CVE_AGGREGATOR_PASS", "")
	Cfg.ManagerPass = GetEnv("USER_MANAGER_PASS", "")
	Cfg.CleanerPass = GetEnv("USER_CLEANER_PASS", "")
	Cfg.ExploitGathererPass = GetEnv("USER_EXPLOIT_GATHERER_PASS", "")
	Cfg.SchemaMigration = GetEnv("SCHEMA_MIGRATION", 0)

	// Manager config
	Cfg.AmsEnabled = GetEnv("AMS_ENABLED", false)
	Cfg.AmsAPIURL = GetEnv("AMS_API_URL", "http://ams_api_url")
	Cfg.AmsAPIPagesize = GetEnv("AMS_API_PAGESIZE", -1)
	Cfg.AmsClientID = GetEnv("AMS_CLIENT_ID", "")
	Cfg.AmsClientSecret = GetEnv("AMS_CLIENT_SECRET", "")

	// Digest writer config
	requestedKafkaBrokerTopic := GetEnv("KAFKA_BROKER_INCOMING_TOPIC", "")
	requestedKafkaPayloadTrackerTopic := GetEnv("KAFKA_PAYLOAD_TRACKER_TOPIC", "")
	if clowder.IsClowderEnabled() {
		if len(clowder.LoadedConfig.Kafka.Brokers) > 0 {
			Cfg.KafkaBroker = clowder.LoadedConfig.Kafka.Brokers[0]
			Cfg.KafkaBrokerAddress = fmt.Sprintf("%s:%d", Cfg.KafkaBroker.Hostname, *Cfg.KafkaBroker.Port)
		}
		if _, ok := clowder.KafkaTopics[requestedKafkaBrokerTopic]; ok {
			Cfg.KafkaBrokerIncomingTopic = clowder.KafkaTopics[requestedKafkaBrokerTopic].Name
		}
		if _, ok := clowder.KafkaTopics[requestedKafkaPayloadTrackerTopic]; ok {
			Cfg.KafkaPayloadTrackerTopic = clowder.KafkaTopics[requestedKafkaPayloadTrackerTopic].Name
		}
	} else {
		Cfg.KafkaBrokerAddress = GetEnv("KAFKA_BROKER_ADDRESS", "")
		Cfg.KafkaBrokerIncomingTopic = requestedKafkaBrokerTopic
		Cfg.KafkaPayloadTrackerTopic = requestedKafkaPayloadTrackerTopic
	}
	Cfg.KafkaBrokerConsumerGroup = GetEnv("KAFKA_BROKER_CONSUMER_GROUP", "")
	Cfg.KafkaConsumerTimeout = GetEnv("KAFKA_CONSUMER_TIMEOUT", "")
	Cfg.KafkaProducerTimeout = GetEnv("KAFKA_PRODUCER_TIMEOUT", "")
	Cfg.PayloadTrackerEnabled = GetEnv("PAYLOAD_TRACKER_ENABLED", false)

	// VMaaS sync config
	Cfg.VmaasBaseURL = GetEnv("VMAAS_BASE_URL", "http://localhost")
	Cfg.VmaasBatchSize = GetEnv("VMAAS_BATCH_SIZE", 0)
	Cfg.VmaasPageSize = GetEnv("VMAAS_PAGE_SIZE", 0)

	// Pyxis gatherer config
	Cfg.PyxisBaseURL = GetEnv("PYXIS_BASE_URL", "http://localhost")
	Cfg.PyxisProfile = GetEnv("PYXIS_PROFILE", "unknown_profile")
	Cfg.ForceSync = GetEnv("FORCE_SYNC", false)
	Cfg.SkipImageCveSync = GetEnv("SKIP_IMAGE_CVE_SYNC", false)

	// Exploits sync config
	Cfg.ExploitFileURL = GetEnv("EXPLOIT_FILE_URL", "http://localhost")
	Cfg.GitToken = GetEnv("GIT_TOKEN", "")
}

// CreateKafkaConfig adds SSL kafka sarama configuration based on clowder
func SetKafkaSSLConfig(config *sarama.Config) error {
	if config == nil {
		return errors.New("sarama config is required")
	}

	broker := Cfg.KafkaBroker

	if broker.Sasl == nil {
		return errors.New("sasl config on kafka broker does not exist")
	}

	saslMechanism := broker.Sasl.SaslMechanism
	if saslMechanism == nil || *saslMechanism == "" {
		return errors.New("sasl mechanism not specified")
	}

	switch strings.ToLower(*saslMechanism) {
	case "plain":
		config.Net.SASL.Mechanism = sarama.SASLTypePlaintext
	case "scram-sha-256":
		config.Net.SASL.Mechanism = sarama.SASLTypeSCRAMSHA256
		config.Net.SASL.SCRAMClientGeneratorFunc = func() sarama.SCRAMClient { return &XDGSCRAMClient{HashGeneratorFcn: SHA256} }
	case "scram-sha-512":
		config.Net.SASL.Mechanism = sarama.SASLTypeSCRAMSHA512
		config.Net.SASL.SCRAMClientGeneratorFunc = func() sarama.SCRAMClient { return &XDGSCRAMClient{HashGeneratorFcn: SHA512} }
	default:
		return fmt.Errorf("unknown sasl mechanism: %s", *saslMechanism)
	}

	if broker.Sasl.Username == nil {
		return errors.New("sasl username not specified")
	}
	if broker.Sasl.Password == nil {
		return errors.New("sasl password not specified")
	}

	config.Net.SASL.User = *broker.Sasl.Username
	config.Net.SASL.Password = *broker.Sasl.Password
	config.Net.SASL.Handshake = true
	config.Net.SASL.Enable = true

	return nil
}

// SetKafkaTLSConfig adds TLS kafka sarama configuration based on clowder
func SetKafkaTLSConfig(config *sarama.Config) error {
	broker := Cfg.KafkaBroker
	tlsConfig := tls.Config{}

	if broker.Cacert != nil && *broker.Cacert != "" {
		certPool := x509.NewCertPool()
		certPool.AppendCertsFromPEM([]byte(*broker.Cacert))

		tlsConfig.RootCAs = certPool
	}
	config.Net.TLS.Enable = true
	config.Net.TLS.Config = &tlsConfig

	return nil
}
