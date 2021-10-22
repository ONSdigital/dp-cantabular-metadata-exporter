module github.com/ONSdigital/dp-cantabular-metadata-exporter

go 1.17

replace github.com/coreos/etcd => github.com/coreos/etcd v3.3.24+incompatible

require (
	github.com/ONSdigital/dp-api-clients-go v1.41.1
	github.com/ONSdigital/dp-api-clients-go/v2 v2.3.0
	github.com/ONSdigital/dp-component-test v0.6.0
	github.com/ONSdigital/dp-healthcheck v1.1.3
	github.com/ONSdigital/dp-kafka/v2 v2.4.2
	github.com/ONSdigital/dp-net v1.2.0
	github.com/ONSdigital/dp-s3 v1.7.0
	github.com/ONSdigital/go-ns v0.0.0-20191104121206-f144c4ec2e58
	github.com/ONSdigital/log.go/v2 v2.0.6
	github.com/aws/aws-sdk-go v1.41.6
	github.com/cucumber/godog v0.12.1
	github.com/go-chi/chi/v5 v5.0.4
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/pkg/errors v0.9.1
	github.com/smartystreets/goconvey v1.6.6
)
