package config

import (
	"os"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestConfig(t *testing.T) {
	os.Clearenv()
	var err error
	var configuration *Config

	Convey("Given an environment with no environment variables set", t, func() {
		Convey("Then cfg should be nil", func() {
			So(cfg, ShouldBeNil)
		})

		Convey("When the config values are retrieved", func() {
			Convey("Then there should be no error returned, and values are as expected", func() {
				configuration, err = Get() // This Get() is only called once, when inside this function
				So(err, ShouldBeNil)
				So(configuration, ShouldResemble, &Config{
					BindAddr:                   ":26700",
					GracefulShutdownTimeout:    5 * time.Second,
					HealthCheckInterval:        30 * time.Second,
					HealthCheckCriticalTimeout: 90 * time.Second,
					DefaultRequestTimeout:      10 * time.Second,
					VaultPath:                  "secret/shared/psk",
					VaultAddress:               "http://localhost:8200",
					VaultToken:                 "",
					PublicBucket:               "dp-cantabular-metadata-exporter",
					PrivateBucket:              "dp-cantabular-metadata-exporter",
					DatasetAPIURL:              "http://localhost:22000",
					FilterAPIURL:               "http://localhost:22100",
					PopulationTypesAPIURL:      "http://localhost:27300",
					ExternalPrefixURL:          "http://localhost:22000",
					DownloadServiceURL:         "http://localhost:23600",
					CantabularURL:              "http://localhost:8491",
					CantabularExtURL:           "http://localhost:8492",
					AWSRegion:                  "eu-west-1",
					StopConsumingOnUnhealthy:   true,
					S3PublicURL:                "http://public-bucket",
					Kafka: KafkaConfig{
						Addr:                          []string{"localhost:9092", "localhost:9093", "localhost:9094"},
						Version:                       "1.0.2",
						OffsetOldest:                  true,
						NumWorkers:                    1,
						MaxBytes:                      2000000,
						SecProtocol:                   "",
						SecCACerts:                    "",
						SecClientKey:                  "",
						SecClientCert:                 "",
						SecSkipVerify:                 false,
						ConsumerMinBrokersHealthy:     1,
						ProducerMinBrokersHealthy:     2,
						CantabularMetadataExportGroup: "cantabular-metadata-export",
						CantabularCSVCreatedTopic:     "cantabular-csv-created",
						CantabularCSVWCreatedTopic:    "cantabular-csvw-created",
					},
				})
			})

			Convey("Then a second call to config should return the same config", func() {
				// This achieves code coverage of the first return in the Get() function.
				newCfg, newErr := Get()
				So(newErr, ShouldBeNil)
				So(newCfg, ShouldResemble, cfg)
			})
		})
	})
}
