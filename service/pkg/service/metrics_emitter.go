package service

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"go.ajitem.com/zapdriver"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var MINUTES_PER_OSCILATION = 2

func NewMetricsEmitterService(cfg MetricsEmitterConfig) MetricsEmitterService {
	configureLogger(cfg)
	return MetricsEmitterService{config: cfg}
}

type MetricsEmitterConfig struct {
	Port        int
	Debug       bool
	Development bool
}

type MetricsEmitterService struct {
	config MetricsEmitterConfig
}

func (s *MetricsEmitterService) Run() {
	e := echo.New()

	start := time.Now()

	e.GET("/metrics", func(c echo.Context) error {
		return c.String(http.StatusOK, promHandler(c, start))
	})
	zap.L().Info("server started.")
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", s.config.Port)))
}

func promHandler(c echo.Context, start time.Time) string {
	minutes_duration := int(time.Now().Sub(start).Minutes())
	zap.L().Debug(fmt.Sprintf("debug: minutes_duration: %d\n", minutes_duration))

	var value int
	var metric string
	if (minutes_duration/MINUTES_PER_OSCILATION)%2 == 0 {
		zap.L().Debug("Returning 0.")
		value = 0
		metric = fmt.Sprintf("metric_a{reason=\"OOMKilled\", pod=\"tempo-ingester-0\", one=\"one\"} %d", value)
	} else {
		zap.L().Debug("Returning 1.")
		value = 1
		metric = fmt.Sprintf("metric_b{reason=\"OOMKilled\", pod=\"tempo-ingester-0\", one=\"one\"} %d", value)
	}

	return metric
}

func configureLogger(cfg MetricsEmitterConfig) {
	var logConfig zap.Config
	if cfg.Development {
		// Logs debug and above
		logConfig = zapdriver.NewDevelopmentConfig()
	} else {
		// Logs info and above
		logConfig = zapdriver.NewProductionConfig()
	}

	if cfg.Debug {
		logConfig.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	} else {
		logConfig.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	}

	logConfig.OutputPaths = []string{"stdout"}

	logger, err := logConfig.Build(zapdriver.WrapCore(
		zapdriver.ReportAllErrors(true),
		zapdriver.ServiceName("MetricsEmitter"),
		zapdriver.ServiceVersion("0.0.1"),
	))
	if err != nil {
		panic("could not start logger.")
	}

	defer logger.Sync()

	zap.ReplaceGlobals(logger)
}
