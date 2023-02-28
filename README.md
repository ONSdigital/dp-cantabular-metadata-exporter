dp-cantabular-metadata-exporter
================
dp-cantabular-metadata-exporter

### Getting started

* Run `make debug`

### Dependencies

* No further dependencies other than those defined in `go.mod`

### Configuration

| Environment variable                | Default                              | Description
| ----------------------------------- | ------------------------------------ | -----------
| BIND_ADDR                           | localhost:26300                      | The host and port to bind to
| GRACEFUL_SHUTDOWN_TIMEOUT           | 5s                                   | The graceful shutdown timeout in seconds (`time.Duration` format)
| HEALTHCHECK_INTERVAL                | 30s                                  | Time between self-healthchecks (`time.Duration` format)
| HEALTHCHECK_CRITICAL_TIMEOUT        | 90s                                  | Time to wait until an unhealthy dependent propagates its state to make this app unhealthy (`time.Duration` format)
| DEFAULT_REQUEST_TIMEOUT             | 10s                                  | Default timeout for graphQL queries against Cantabular API extension and HTTP requests for the Cantabular Client only
| CANTABULAR_URL                      | http://localhost:8491                | The Cantabular server URL
| CANTABULAR_EXT_API_URL              | http://localhost:8492                | The Cantabular API extension URL
| DATASET_API_URL                     | http://localhost:22000               | The Dataset API URL
| POPULATION_TYPES_API_URL            | http://localhost:27300               | The Population Types API URL
| DOWNLOAD_SERVICE_URL                | http://localhost:23600               | The Download Service URL, only used to generate download links
| FILTER_API_URL                      | http://localhost:22100               | The Fiilter API URL
| EXTERNAL_PREFIX_URL                 | http://localhost:22000               | The string used to add environment prefixes to URL's
| SERVICE_AUTH_TOKEN                  |                                      | The service token for this app
| CANTABULAR_HEALTHCHECK_ENABLED      | false                                | Flag to enable/disable healthchecks against Cantabular server and API extension
| AWS_REGION                          | eu-west-1                            | The AWS region to use
| PUBLIC_BUCKET                       | public-bucket                        | The name of the S3 bucket to store published files
| PRIVATE_BUCKET                      | private-bucket                       | The name of the S3 bucket to store un-published files
| PUBLIC_URL                          |                                      | The S3 bucket url
| S3_PUBLIC_URL                       |                                      | The S3 public bucket url
| LOCAL_OBJECT_STORE                  |                                      | Used during feature tests
| MINIO_ACCESS_KEY                    |                                      | Used during feature tests
| MINIO_SECRET_KEY                    |                                      | Used during feature tests
| COMPONENT_TEST_USE_LOG_FILE         |                                      | Used during feature tests
| VAULT_ADDR                          | http://localhost:8200                | The address of vault
| VAULT_TOKEN                         | -                                    | Use `make debug` to set a vault token
| VAULT_PATH                          | secret/shared/psk                    | The vault path to store psks
| ENCRYPTION_DISABLED                 | false                                | Flag to enable/disable encryption for un-published files
| STOP_CONSUMING_ON_UNHEALTHY         | true                                 | Flag to enable/disable kafka-consumer consumption depending on health status. If true, the consumer will stop consuming on 'WARNING' and 'CRITICAL' and it will start consuming on 'OK'
| KAFKA_ADDR                          | localhost:9092                       | The kafka broker addresses (can be comma separated)
| KAFKA_VERSION                       | "1.0.2"                              | The kafka version that this service expects to connect to
| KAFKA_OFFSET_OLDEST                 | true                                 | Start processing Kafka messages in order from the oldest in the queue
| KAFKA_NUM_WORKERS                   | 1                                    | The maximum number of parallel kafka consumers
| KAFKA_MAX_BYTES                     | 2000000                              | the maximum number of bytes per kafka message
| KAFKA_TOPIC_CSV_CREATED             | cantabular-csv-created               | The name of the topic that is produced after a CSV file has been successfully generated
| KAFKA_TOPIC_CSVW_CREATED            | cantabular-csvw-created               | The name of the topic that is produced after a CSVW file has been successfully generated
| KAFKA_SEC_PROTO                     | _unset_                              | if set to `TLS`, kafka connections will use TLS [[1]](#notes_1)
| KAFKA_SEC_CA_CERTS                  | _unset_                              | CA cert chain for the server cert [[1]](#notes_1)
| KAFKA_SEC_CLIENT_KEY                | _unset_                              | PEM for the client key [[1]](#notes_1)
| KAFKA_SEC_CLIENT_CERT               | _unset_                              | PEM for the client certificate [[1]](#notes_1)
| KAFKA_SEC_SKIP_VERIFY               | false                                | ignores server certificate issues if `true` [[1]](#notes_1)

**Notes:**

	1. <a name="notes_1">For more info, see the [kafka TLS examples documentation](https://github.com/ONSdigital/dp-kafka/tree/main/examples#tls)</a>


### Healthcheck

 The `/health` endpoint returns the current status of the service. Dependent services are health checked on an interval defined by the `HEALTHCHECK_INTERVAL` environment variable.

 On a development machine a request to the health check endpoint can be made by:

 `curl localhost:8125/health`

### Contributing

See [CONTRIBUTING](CONTRIBUTING.md) for details.

### License

Copyright © 2021, Office for National Statistics (https://www.ons.gov.uk)

Released under MIT license, see [LICENSE](LICENSE.md) for details.

