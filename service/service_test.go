package service_test

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/ONSdigital/dp-cantabular-metadata-exporter/config"
	"github.com/ONSdigital/dp-cantabular-metadata-exporter/service"
	"github.com/ONSdigital/dp-cantabular-metadata-exporter/service/mock"

	"github.com/ONSdigital/dp-healthcheck/healthcheck"
	kafka "github.com/ONSdigital/dp-kafka/v3"
	"github.com/ONSdigital/dp-kafka/v3/kafkatest"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	ctx           = context.Background()
	testBuildTime = "12"
	testGitCommit = "GitCommit"
	testVersion   = "Version"
)

var errHealthcheck = errors.New("could not get healthcheck")

func TestInit(t *testing.T) {
	Convey("Having a set of mocked dependencies", t, func() {
		cfg, err := config.Get()
		So(err, ShouldBeNil)

		hcMock := &mock.HealthCheckerMock{
			AddAndGetCheckFunc: func(name string, checker healthcheck.Checker) (*healthcheck.Check, error) {
				return &healthcheck.Check{}, nil
			},
			StartFunc:        func(ctx context.Context) {},
			StopFunc:         func() {},
			SubscribeAllFunc: func(s healthcheck.Subscriber) {},
		}

		service.GetHealthCheck = func(cfg *config.Config, buildTime, gitCommit, version string) (service.HealthChecker, error) {
			return hcMock, nil
		}

		service.GetKafkaProducer = func(ctx context.Context, cfg *config.Config) (kafka.IProducer, error) {
			return &kafkatest.IProducerMock{
				ChannelsFunc: kafka.CreateProducerChannels,
			}, nil
		}

		service.GetKafkaConsumer = func(ctx context.Context, cfg *config.Config) (kafka.IConsumerGroup, error) {
			return &kafkatest.IConsumerGroupMock{
				ChannelsFunc: func() *kafka.ConsumerGroupChannels {
					return kafka.CreateConsumerGroupChannels(1, 1)
				},
				RegisterHandlerFunc: func(ctx context.Context, h kafka.Handler) error {
					return nil
				},
			}, nil
		}

		serverMock := &mock.HTTPServerMock{
			ListenAndServeFunc: func() error {
				return nil
			},
			ShutdownFunc: func(ctx context.Context) error {
				return nil
			},
		}
		service.GetHTTPServer = func(bindAddr string, router http.Handler) service.HTTPServer {
			return serverMock
		}

		svc := &service.Service{}

		Convey("Given that initialising healthcheck returns an error", func() {
			service.GetHealthCheck = func(cfg *config.Config, buildTime, gitCommit, version string) (service.HealthChecker, error) {
				return nil, errHealthcheck
			}
			// setup (run before each `Convey` at this scope / indentation):
			newSvc := service.New()
			err := newSvc.Init(ctx, cfg, testBuildTime, testGitCommit, testVersion)

			Convey("Then service Init fails with an error", func() {
				So(errors.Is(err, errHealthcheck), ShouldBeTrue)
			})

			Reset(func() {
				// This reset is run after each `Convey` at the same scope (indentation)
			})
		})

		Convey("Given that all dependencies are successfully initialised", func() {
			// setup (run before each `Convey` at this scope / indentation):
			err := svc.Init(ctx, cfg, testBuildTime, testGitCommit, testVersion)

			Convey("Then service Init succeeds", func() {
				So(err, ShouldBeNil)
			})

			Reset(func() {
				// This reset is run after each `Convey` at the same scope (indentation)
			})
		})
	})
}

func TestClose(t *testing.T) {
	Convey("Having a correctly initialised service", t, func() {
		cfg, err := config.Get()
		So(err, ShouldBeNil)

		hcStopped := false

		// healthcheck Stop does not depend on any other service being closed/stopped
		hcMock := &mock.HealthCheckerMock{
			AddAndGetCheckFunc: func(name string, checker healthcheck.Checker) (*healthcheck.Check, error) {
				return &healthcheck.Check{}, nil
			},
			StartFunc:        func(ctx context.Context) {},
			StopFunc:         func() { hcStopped = true },
			SubscribeAllFunc: func(s healthcheck.Subscriber) {},
		}
		service.GetHealthCheck = func(cfg *config.Config, buildTime, gitCommit, version string) (service.HealthChecker, error) {
			return hcMock, nil
		}

		// server Shutdown will fail if healthcheck is not stopped
		serverMock := &mock.HTTPServerMock{
			ListenAndServeFunc: func() error { return nil },
			ShutdownFunc: func(ctx context.Context) error {
				if !hcStopped {
					return errors.New("server stopped before healthcheck")
				}
				return nil
			},
		}

		service.GetHTTPServer = func(bindAddr string, router http.Handler) service.HTTPServer {
			return serverMock
		}

		pc := kafka.CreateProducerChannels()
		service.GetKafkaProducer = func(ctx context.Context, cfg *config.Config) (kafka.IProducer, error) {
			return &kafkatest.IProducerMock{
				ChannelsFunc: func() *kafka.ProducerChannels {
					return pc
				},
				CloseFunc: func(context.Context) error {
					return nil
				},
				LogErrorsFunc: func(ctx context.Context) {},
			}, nil
		}

		cgc := kafka.CreateConsumerGroupChannels(1, 1)
		cgc.State = &kafka.ConsumerStateChannels{
			Consuming: kafka.NewStateChan(),
		}
		service.GetKafkaConsumer = func(ctx context.Context, cfg *config.Config) (kafka.IConsumerGroup, error) {
			return &kafkatest.IConsumerGroupMock{
				ChannelsFunc: func() *kafka.ConsumerGroupChannels {
					return cgc
				},
				LogErrorsFunc: func(ctx context.Context) {},
				StartFunc:     func() error { return nil },
				CloseFunc: func(ctx context.Context, optFuncs ...kafka.OptFunc) error {
					return nil
				},
				RegisterHandlerFunc: func(ctx context.Context, h kafka.Handler) error {
					return nil
				},
				StopAndWaitFunc: func() error { return nil },
			}, nil
		}

		Convey("Closing the service results in all the dependencies being closed in the expected order", func() {
			svcErrors := make(chan error, 1)
			svc := service.New()
			err := svc.Init(ctx, cfg, testBuildTime, testGitCommit, testVersion)
			So(err, ShouldBeNil)

			err = svc.Start(context.Background(), svcErrors)
			So(err, ShouldBeNil)

			err = svc.Close(context.Background())
			So(err, ShouldBeNil)
			So(len(hcMock.StopCalls()), ShouldEqual, 1)
			So(len(serverMock.ShutdownCalls()), ShouldEqual, 1)
		})

		Convey("If services fail to stop, the Close operation tries to close all dependencies and returns an error", func() {
			failingserverMock := &mock.HTTPServerMock{
				ListenAndServeFunc: func() error { return nil },
				ShutdownFunc: func(ctx context.Context) error {
					return errors.New("failed to stop http server")
				},
			}

			service.GetHTTPServer = func(bindAddr string, router http.Handler) service.HTTPServer {
				return failingserverMock
			}

			svcErrors := make(chan error, 1)
			svc := service.New()
			err := svc.Init(ctx, cfg, testBuildTime, testGitCommit, testVersion)
			So(err, ShouldBeNil)

			err = svc.Start(context.Background(), svcErrors)
			So(err, ShouldBeNil)

			err = svc.Close(context.Background())
			So(err, ShouldNotBeNil)
			So(len(hcMock.StopCalls()), ShouldEqual, 1)
			So(len(failingserverMock.ShutdownCalls()), ShouldEqual, 1)
		})

		Convey("If service times out while shutting down, the Close operation fails with the expected error", func() {
			cfg.GracefulShutdownTimeout = 1 * time.Millisecond
			timeoutServerMock := &mock.HTTPServerMock{
				ListenAndServeFunc: func() error { return nil },
				ShutdownFunc: func(ctx context.Context) error {
					time.Sleep(2 * time.Millisecond)
					return nil
				},
			}

			svc := service.Service{
				Config:      cfg,
				Server:      timeoutServerMock,
				HealthCheck: hcMock,
			}

			err = svc.Close(context.Background())
			So(err, ShouldNotBeNil)
			So(errors.Is(err, context.DeadlineExceeded), ShouldBeTrue)
			So(len(hcMock.StopCalls()), ShouldEqual, 1)
			So(len(timeoutServerMock.ShutdownCalls()), ShouldEqual, 1)
		})
	})
}
