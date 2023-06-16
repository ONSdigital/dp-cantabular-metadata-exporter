# dp-cantabular-metadata-exporter

================

dp-cantabular-metadata-exporter

## Getting started

- Run `make debug`

### Dependencies

- No further dependencies other than those defined in `go.mod`

### Configuration

| Environment variable                   | Default                    | Description                                                                                                                                                                             |
| -------------------------------------- | -------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| AWS_REGION                             | eu-west-1                  | The AWS region to use                                                                                                                                                                   |
| BIND_ADDR                              | localhost:26700            | The host and port to bind to                                                                                                                                                            |
| CANTABULAR_EXT_API_URL                 | <http://localhost:8492>    | The Cantabular API extension URL                                                                                                                                                        |
| CANTABULAR_URL                         | <http://localhost:8491>    | The Cantabular server URL                                                                                                                                                               |
| COMPONENT_TEST_USE_LOG_FILE            |                            | Used during feature tests                                                                                                                                                               |
| DATASET_API_URL                        | <http://localhost:22000>   | The Dataset API URL                                                                                                                                                                     |
| DEFAULT_REQUEST_TIMEOUT                | 10s                        | Default timeout for graphQL queries against Cantabular API extension and HTTP requests for the Cantabular Client only                                                                   |
| DOWNLOAD_SERVICE_URL                   | <http://localhost:23600>   | The Download Service URL, only used to generate download links                                                                                                                          |
| ENCRYPTION_DISABLED                    | false                      | Flag to enable/disable encryption for un-published files                                                                                                                                |
| EXTERNAL_PREFIX_URL                    | <http://localhost:22000>   | The string used to add environment prefixes to URL's                                                                                                                                    |
| FILTER_API_URL                         | <http://localhost:22100>   | The Fiilter API URL                                                                                                                                                                     |
| GRACEFUL_SHUTDOWN_TIMEOUT              | 5s                         | The graceful shutdown timeout in seconds (`time.Duration` format)                                                                                                                       |
| HEALTHCHECK_CRITICAL_TIMEOUT           | 90s                        | Time to wait until an unhealthy dependent propagates its state to make this app unhealthy (`time.Duration` format)                                                                      |
| HEALTHCHECK_INTERVAL                   | 30s                        | Time between self-healthchecks (`time.Duration` format)                                                                                                                                 |
| KAFKA_ADDR                             | localhost:9092             | The kafka broker addresses (can be comma separated)                                                                                                                                     |
| KAFKA_CONSUMER_MIN_BROKERS_HEALTHY     | 1                          | the minimum number of healthy brokers                                                                                                                                                   |
| KAFKA_GROUP_CANTABULAR_METADATA_EXPORT | cantabular-metadata-export | The cantabular metadata export group                                                                                                                                                    |
| KAFKA_MAX_BYTES                        | 2000000                    | the maximum number of bytes per kafka message                                                                                                                                           |
| KAFKA_NUM_WORKERS                      | 1                          | The maximum number of parallel kafka consumers                                                                                                                                          |
| KAFKA_OFFSET_OLDEST                    | true                       | Start processing Kafka messages in order from the oldest in the queue                                                                                                                   |
| KAFKA_PRODUCER_MIN_BROKERS_HEALTHY     | 1                          | The minimum number of healthy brokers                                                                                                                                                   |
| KAFKA_SEC_CA_CERTS                     | _unset_                    | CA cert chain for the server cert [^1]                                                                                                                                                  |
| KAFKA_SEC_CLIENT_CERT                  | _unset_                    | PEM for the client certificate [^1]                                                                                                                                                     |
| KAFKA_SEC_CLIENT_KEY                   | _unset_                    | PEM for the client key [^1]                                                                                                                                                             |
| KAFKA_SEC_PROTO                        | _unset_                    | if set to `TLS`, kafka connections will use TLS [^1]                                                                                                                                    |
| KAFKA_SEC_SKIP_VERIFY                  | false                      | ignores server certificate issues if `true` [^1]                                                                                                                                        |
| KAFKA_TOPIC_CSV_CREATED                | cantabular-csv-created     | The name of the topic that is produced after a CSV file has been successfully generated                                                                                                 |
| KAFKA_TOPIC_CSVW_CREATED               | cantabular-csvw-created    | The name of the topic that is produced after a CSVW file has been successfully generated                                                                                                |
| KAFKA_VERSION                          | "1.0.2"                    | The kafka version that this service expects to connect to                                                                                                                               |
| LOCAL_OBJECT_STORE                     |                            | Used during feature tests                                                                                                                                                               |
| MINIO_ACCESS_KEY                       |                            | Used during feature tests                                                                                                                                                               |
| MINIO_SECRET_KEY                       |                            | Used during feature tests                                                                                                                                                               |
| POPULATION_TYPES_API_URL               | <http://localhost:27300>   | The Population Types API URL                                                                                                                                                            |
| PRIVATE_BUCKET                         | private-bucket             | The name of the S3 bucket to store un-published files                                                                                                                                   |
| PUBLIC_BUCKET                          | public-bucket              | The name of the S3 bucket to store published files                                                                                                                                      |
| PUBLIC_URL                             |                            | The S3 bucket url                                                                                                                                                                       |
| SERVICE_AUTH_TOKEN                     |                            | The service token for this app                                                                                                                                                          |
| STOP_CONSUMING_ON_UNHEALTHY            | true                       | Flag to enable/disable kafka-consumer consumption depending on health status. If true, the consumer will stop consuming on 'WARNING' and 'CRITICAL' and it will start consuming on 'OK' |
| S3_PUBLIC_URL                          |                            | The S3 public bucket url                                                                                                                                                                |
| VAULT_ADDR                             | <http://localhost:8200>    | The address of vault                                                                                                                                                                    |
| VAULT_PATH                             | secret/shared/psk          | The vault path to store psks                                                                                                                                                            |
| VAULT_TOKEN                            | -                          | Use `make debug` to set a vault token                                                                                                                                                   |

[^1]: For more info, see the [kafka TLS examples documentation](https://github.com/ONSdigital/dp-kafka/tree/main/examples#tls)

### Healthcheck

The `/health` endpoint returns the current status of the service. Dependent services are health checked on an interval defined by the `HEALTHCHECK_INTERVAL` environment variable.

On a development machine a request to the health check endpoint can be made by:

`curl localhost:8125/health`

### Contributing

See [CONTRIBUTING](CONTRIBUTING.md) for details.

### License

Copyright © 2023, Office for National Statistics (<https://www.ons.gov.uk>)

Released under MIT license, see [LICENSE](LICENSE.md) for details
