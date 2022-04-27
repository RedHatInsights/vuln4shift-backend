package utils

import (
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
	LoggingLevel string
	APIRetries   int

	// DB admin config
	ArchiveDbWriterPass string
	PyxisGathererPass   string
	VmaasGathererPass   string
	CveAggregatorPass   string
	ManagerPass         string
	SchemaMigration     int

	// Digest writer config
	KafkaBrokerAddress       string
	KafkaBrokerConsumerGroup string
	KafkaBrokerIncomingTopic string
	KafkaConsumerTimeout     string

	// VMaaS sync config
	VmaasBaseURL   string
	VmaasBatchSize int
	VmaasPageSize  int

	// Pyxis gatherer config
	PyxisBaseURL string
	PyxisProfile string
}

func init() {
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

	// DB admin config
	Cfg.ArchiveDbWriterPass = GetEnv("USER_ARCHIVE_DB_WRITER_PASS", "")
	Cfg.PyxisGathererPass = GetEnv("USER_PYXIS_GATHERER_PASS", "")
	Cfg.VmaasGathererPass = GetEnv("USER_VMAAS_GATHERER_PASS", "")
	Cfg.CveAggregatorPass = GetEnv("USER_CVE_AGGREGATOR_PASS", "")
	Cfg.ManagerPass = GetEnv("USER_MANAGER_PASS", "")
	Cfg.SchemaMigration = GetEnv("SCHEMA_MIGRATION", 0)

	// Digest writer config
	requestedKafkaBrokerTopic := GetEnv("KAFKA_BROKER_INCOMING_TOPIC", "")
	if clowder.IsClowderEnabled() {
		if len(clowder.KafkaServers) > 0 {
			Cfg.KafkaBrokerAddress = clowder.KafkaServers[0]
		}
		if _, ok := clowder.KafkaTopics[requestedKafkaBrokerTopic]; ok {
			Cfg.KafkaBrokerIncomingTopic = clowder.KafkaTopics[requestedKafkaBrokerTopic].Name
		}
	} else {
		Cfg.KafkaBrokerAddress = GetEnv("KAFKA_BROKER_ADDRESS", "")
		Cfg.KafkaBrokerIncomingTopic = requestedKafkaBrokerTopic
	}
	Cfg.KafkaBrokerConsumerGroup = GetEnv("KAFKA_BROKER_CONSUMER_GROUP", "")
	Cfg.KafkaConsumerTimeout = GetEnv("KAFKA_CONSUMER_TIMEOUT", "")

	// VMaaS sync config
	Cfg.VmaasBaseURL = GetEnv("VMAAS_BASE_URL", "http://localhost")
	Cfg.VmaasBatchSize = GetEnv("VMAAS_BATCH_SIZE", 0)
	Cfg.VmaasPageSize = GetEnv("VMAAS_PAGE_SIZE", 0)

	// Pyxis gatherer config
	Cfg.PyxisBaseURL = GetEnv("PYXIS_BASE_URL", "http://localhost")
	Cfg.PyxisProfile = GetEnv("PYXIS_PROFILE", "unknown_profile")
}
