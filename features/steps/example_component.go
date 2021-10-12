package steps

import (
	"context"
	"fmt"
	"net/http"
	
	"github.com/ONSdigital/dp-cantabular-metadata-exporter/config"
	"github.com/ONSdigital/dp-cantabular-metadata-exporter/service"
	"github.com/ONSdigital/dp-cantabular-metadata-exporter/service/mock"
	componenttest "github.com/ONSdigital/dp-component-test"
	"github.com/ONSdigital/dp-healthcheck/healthcheck"
)

type Component struct {
	componenttest.ErrorFeature
	svc            *service.Service
	errorChan      chan error
	Config         *config.Config
	HTTPServer     *http.Server
	ServiceRunning bool
	apiFeature     *componenttest.APIFeature
}

func NewComponent() (*Component, error) {

	c := &Component{
		HTTPServer:     &http.Server{},
		errorChan:      make(chan error),
		ServiceRunning: false,
	}

	var err error

	c.Config, err = config.Get()
	if err != nil {
		return nil, err
	}

	c.apiFeature = componenttest.NewAPIFeature(c.InitialiseService)

	return c, nil
}

func (c *Component) Reset() *Component {
	c.apiFeature.Reset()
	return c
}

func (c *Component) Close() error {
	if c.svc != nil && c.ServiceRunning {
		c.svc.Close(context.Background())
		c.ServiceRunning = false
	}
	return nil
}

func (c *Component) InitialiseService() (http.Handler, error) {
	var err error
	c.svc = service.New()
	if err = c.svc.Init(context.Background(), c.Config, "1", "", ""); err != nil {
		return nil, fmt.Errorf("failed to initialise service: %w", err)
	}

	c.ServiceRunning = true
	return c.HTTPServer.Handler, nil
}

func (c *Component) DoGetHealthcheckOk(cfg *config.Config, buildTime string, gitCommit string, version string) (service.HealthChecker, error) {
	return &mock.HealthCheckerMock{
		AddCheckFunc: func(name string, checker healthcheck.Checker) error { return nil },
		StartFunc:    func(ctx context.Context) {},
		StopFunc:     func() {},
	}, nil
}

func (c *Component) DoGetHTTPServer(bindAddr string, router http.Handler) service.HTTPServer {
	c.HTTPServer.Addr = bindAddr
	c.HTTPServer.Handler = router
	return c.HTTPServer
}
