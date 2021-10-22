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
					BindAddr:                      ":26700",
					GracefulShutdownTimeout:       5 * time.Second,
					HealthCheckInterval:           30 * time.Second,
					HealthCheckCriticalTimeout:    90 * time.Second,
					VaultPath:                    "secret/shared/psk",
					VaultAddress:                 "http://localhost:8200",
					VaultToken:                   "",
					EncryptionDisabled:           false,
					Kafka: KafkaConfig {
						Addr:                         []string{"localhost:9092"},
						Version:                      "1.0.2",
						OffsetOldest:                 true,
						NumWorkers:                   1,
						MaxBytes:                     2000000,
						SecProtocol:                  "",
						SecCACerts:                   "",
						SecClientKey:                 "",
						SecClientCert:                "",
						SecSkipVerify:                false,
						CantabularMetadataExportGroup: "cantabular-metadata-export",
						CantabularMetadataExportTopic: "cantabular-metadata-export",
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
