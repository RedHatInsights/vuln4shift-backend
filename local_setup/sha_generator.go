package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"math/rand"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/Shopify/sarama"
	_ "github.com/lib/pq" // PostgreSQL database driver
	"github.com/spf13/viper"
)

// Exit codes
const (
	// ExitStatusOK - the tool finished with success
	ExitStatusOK = iota
	// ExitStatusKafkaBrokerError - kafka broker connection establishment error
	ExitStatusKafkaBrokerError
	// ExitStatusKafkaProducerError - kafka event production failure
	ExitStatusKafkaProducerError
	// ExitStatusDBConfigError - DB configuration retrieval error
	ExitStatusDBConfigError
	// ExitStatusDBConnectionError - DB connection establishment error
	ExitStatusDBConnectionError
)

type DatabaseUsernamePerTable struct {
	PGUsername string `mapstructure:"pg_username" toml:"pg_username"`
	PGPassword string `mapstructure:"pg_password" toml:"pg_password"`
}

// DBConfig holds config for connecting to a postgreSQL DB
type DatabaseConfig struct {
	PGHost   string `mapstructure:"pg_host" toml:"pg_host"`
	PGPort   int    `mapstructure:"pg_port" toml:"pg_port"`
	PGDBName string `mapstructure:"pg_db_name" toml:"pg_db_name"`
	PGParams string `mapstructure:"pg_params" toml:"pg_params"`

	PGTableParams map[string]DatabaseUsernamePerTable `mapstructure:"pg_table_params" toml:"pg_table_params"`
}

// CliFlags holds all the allowed command line arguments and flags.
type CliFlags struct {
	OrgID         int
	AccountNumber int
	ClusterName   string
	NumMessages   int
	Produce       bool
	Store         bool
	StoreAccount  bool
	KafkaBroker   string
	KafkaTopic    string
	Digests       string
}

// JSONContent represents any JSON object as key-value mapping
type JSONContent map[string]*json.RawMessage

// Image data structure is representation of Images JSON object
type Image struct {
	Pods       int          `json:"-"`
	ImageCount int          `json:"imageCount"`
	Digests    *JSONContent `json:"images"`
	Namespaces *JSONContent `json:"namespaces"`
}

// KafkaMessage is the structure of JSON messages produced
type KafkaMessage struct {
	Organization  int    `json:"OrgID"`
	AccountNumber int    `json:"AccountNumber"`
	ClusterName   string `json:"ClusterName"`
	Images        *Image `json:"Images"`
}

var (
	VERBOSE  = false
	DBCONFIG *DatabaseConfig
)

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

func produce(shas []string, account, org int, cluster, broker, topic string) {
	producer, err := sarama.NewSyncProducer([]string{broker}, nil)
	if err != nil {
		fmt.Printf("couldn't connect to Kafka broker %v\n", broker)
		os.Exit(ExitStatusKafkaBrokerError)
	}

	kafkaMsg := KafkaMessage{
		AccountNumber: account,
		Organization:  org,
		ClusterName:   cluster,
		Images: &Image{
			ImageCount: len(shas),
		},
	}

	images := make(JSONContent, len(shas))

	//empty content, as long as it is a valid JSON object
	content := json.RawMessage("{}")
	for _, msg := range shas {
		images[msg] = &content
	}

	kafkaMsg.Images.Digests = &images
	kafkaMsg.Images.Namespaces = &images

	if VERBOSE {
		fmt.Printf("content of Kafka message to produce:\n\torg: %v\n\taccount: %v\n\tcluster:%v\n\tDigests:%v\n", kafkaMsg.Organization, kafkaMsg.AccountNumber, kafkaMsg.ClusterName, kafkaMsg.Images.Digests)
	}

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
		fmt.Println(err.Error())
		fmt.Println("failed to produce message to Kafka")
	} else {
		if VERBOSE {
			fmt.Printf("message sent to partition %d at offset %d\n", partitionID, offset)
		}
	}
}

func readYamlConfig() {
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.SetConfigType("yml")

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println(err.Error())
		fmt.Println("Error reading config file")
		os.Exit(ExitStatusDBConfigError)
	}

	if err := viper.Unmarshal(&DBCONFIG); err != nil {
		fmt.Println(err.Error())
		fmt.Println("Unable to decode into config struct")
		os.Exit(ExitStatusDBConfigError)
	}
}

// Setup DB gets the environment variables ONLY from config.yml
func setupDB(table string) (dataSource string) {
	dataSource = fmt.Sprintf(
		"postgresql://%v:%v@%v:%v/%v?%v",
		DBCONFIG.PGTableParams[table].PGUsername,
		url.QueryEscape(DBCONFIG.PGTableParams[table].PGPassword),
		DBCONFIG.PGHost,
		DBCONFIG.PGPort,
		DBCONFIG.PGDBName,
		DBCONFIG.PGParams,
	)
	return
}

func getSQLConnection(dataSource string) *sql.DB {
	connection, err := sql.Open("postgres", dataSource)
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("failed to connect to provided database")
		os.Exit(ExitStatusDBConnectionError)
	}
	if err = connection.Ping(); err != nil {
		fmt.Println(err.Error())
		fmt.Println("failed to connect to provided database")
		os.Exit(ExitStatusDBConnectionError)
	}
	fmt.Println("connected to PostgreSQL database successfully")
	return connection
}

func executeSqlStatement(connection *sql.DB, statement string, args []interface{}) error {
	tx, err := connection.Begin()
	if _, err = tx.Exec(statement, args...); err != nil {
		fmt.Println(err.Error())
		return tx.Rollback()
	}
	return tx.Commit()
}

func prepareInsertDigestsStatement(shas []string) (statement string, statementArgs []interface{}) {
	statement = `INSERT INTO image (digest, pyxis_id, modified_date) VALUES %s`

	rand.Seed(time.Now().UnixNano())
	var valuesIdx []string
	statementIdx := 0
	modifiedDate := time.Now()

	for _, sha := range shas {
		statementArgs = append(statementArgs, sha, rand.Int(), modifiedDate)
		statementIdx = len(statementArgs)
		valuesIdx = append(valuesIdx, "($"+fmt.Sprint(statementIdx-2)+
			", $"+fmt.Sprint(statementIdx-1)+", $"+fmt.Sprint(statementIdx)+")")
	}
	statement = fmt.Sprintf(statement, strings.Join(valuesIdx, ","))
	return
}

func store(shas []string, dataSource string) {
	connection := getSQLConnection(dataSource)
	statement, args := prepareInsertDigestsStatement(shas)
	if VERBOSE {
		fmt.Printf("insert digests SQL statement:\n\t%v\n\t%v\n", statement, args)
	}

	if err := executeSqlStatement(connection, statement, args); err != nil {
		fmt.Println("Something went wrong while inserting digests")
		fmt.Println(err.Error())
	}
	fmt.Printf("inserted %d digests in the 'image' table\n", len(shas))
	_ = connection.Close()
}

func writeAccountAndOrg(orgId int, dataSource string) {
	connection := getSQLConnection(dataSource)
	statement := `INSERT INTO account (org_id) VALUES ($1) ON CONFLICT DO NOTHING`
	args := []interface{}{orgId}
	if VERBOSE {
		fmt.Printf("insert account SQL statement:\n\t%v\n\t%v\n", statement, args)
	}
	if err := executeSqlStatement(connection, statement, args); err != nil {
		fmt.Println("Something went wrong while writing account data")
		fmt.Println(err.Error())
	}

	_ = connection.Close()
}

func main() {
	if len(os.Args) > 1 {
		var flags CliFlags
		flag.IntVar(&flags.OrgID, "org", 1, "organization ID to include in generated message")
		flag.IntVar(&flags.AccountNumber, "account", 1, "account number to include in generated message")
		flag.StringVar(&flags.ClusterName, "cluster", "84f7eedc-0000-0000-9d4d-000000000000", "cluster name to include in generated message")
		flag.IntVar(&flags.NumMessages, "num", 1, "number of SHA256 messages to generate")
		flag.BoolVar(&flags.Produce, "produce", false, "send generated SHAs to configured Kafka topic")
		flag.StringVar(&flags.KafkaBroker, "kafka-broker", "localhost:9093", "Kafka broker in the <host>:<port> format")
		flag.StringVar(&flags.KafkaTopic, "kafka-topic", "test_sha", "Kafka topic for producer")
		flag.BoolVar(&flags.Store, "store", false, "store generated SHAs in the 'image' table of the given DB")
		flag.BoolVar(&flags.StoreAccount, "store-account", false, "update 'account' table of the given account and org-id")
		flag.BoolVar(&VERBOSE, "verbose", false, "print additional information during execution")
		flag.Parse()

		shas := generateSHA256(flags.NumMessages)
		fmt.Println(shas)

		readYamlConfig()
		if flags.StoreAccount {
			writeAccountAndOrg(flags.OrgID, setupDB("account"))
		}
		if flags.Store {
			store(shas, setupDB("image"))
		}

		if !flags.Produce {
			os.Exit(ExitStatusOK)
		}

		produce(shas, flags.OrgID, flags.AccountNumber, flags.ClusterName, flags.KafkaBroker, flags.KafkaTopic)
		os.Exit(ExitStatusOK)
	}

	shas := generateSHA256(1)
	fmt.Println(shas)
}
