help: `go run sha_generator.go -h ` or `go run sha_generator.go -help`

Please see script's help for usage explanation and arguments' default values. All arguments can be combined.

- generate one SHA and print on console (default behavior) - `go run sha_generator.go`
- generate <NUM> SHAs and print on console - `go run sha_generator.go -num NUM`

- specify which org ID and account - `go run sha_generator.go -org 12345 -account 12345678`
  - This will create one SHA-256 as `-num` is not specified.
  - `org` and `account` do not influence the result as it is not sent to Kafka nor stored in DB.

- send generated data to Kafka (see help for related configuration) - `go run sha_generator.go -produce`
- store  generated data in DB (see help for related configuration) - `go run sha_generator.go -store`

- generate <NUM> SHAs and send to kafka - `go run sha_generator.go -num NUM -produce`
- generate <NUM> SHAs and store in DB - `go run sha_generator.go -num NUM -store`
