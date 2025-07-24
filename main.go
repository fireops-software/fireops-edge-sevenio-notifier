package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/fireops-software/fireops-edge-sevenio-notifier/api"
	"github.com/fireops-software/fireops-edge-sevenio-notifier/services"
	"github.com/uoul/go-common/config"
	"github.com/uoul/go-common/log"
	"github.com/uoul/go-common/messaging"
)

const (
	VERSION      = "{VERSION}"
	SERVICE_NAME = "fireops-edge-sevenio-notifier"
)

func main() {
	// Create AppContext
	ctx, cancel := context.WithCancel(context.Background())

	// Create ConfigProvider
	cp := config.NewEnvVarProvider()

	// Create Logger
	logger := log.NewConsoleLogger(
		log.StringToLogLevel(
			cp.StringOrDefault("LOG_LEVEL", "INFO"),
			log.INFO,
		),
	)

	// Create SevenIoApi
	sevenIoClient := api.NewSevenIoApi(
		ctx,
		logger,
		cp.StringOrDefault("SEVENIO_API_URL", ""),
		cp.StringOrDefault("SEVENIO_TOKEN", ""),
	)

	// Create RabbitMqClient
	rabbitMq := messaging.NewRabbitMqMessenger(
		ctx,
		logger,
		cp.StringOrDefault("RABBITMQ_HOST", ""),
		cp.UInt16OrDefault("RABBITMQ_PORT", 5672),
		cp.StringOrDefault("RABBITMQ_USER", ""),
		cp.StringOrDefault("RABBITMQ_PW", ""),
	)

	// Create SevenIoNotifier service
	services.NewSevenIoNotifier(
		ctx,
		logger,
		rabbitMq,
		messaging.RabbitMqExchange{
			Type:       "topic",
			Exchange:   cp.StringOrDefault("RABBITMQ_EXCHANGE", "fireops-edge-events"),
			RoutingKey: cp.StringOrDefault("RABBITMQ_ROUTING_KEY", "alu2g.new"),
		},
		sevenIoClient,
		services.WithSevenIoAlertText(
			cp.StringOrDefault("SMS_ALERT_TEXT", "EINSATZ von FireOps"),
		),
		services.WithSevenIoGroupDefault(
			cp.StringOrDefault("SMS_GROUP_DEFAULT", "Mannschaft"),
		),
		services.WithSevenIoGroupFull(
			cp.StringOrDefault("SMS_GROUP_EXTENDED", "Kommando"),
		),
	)

	// Create HealthReporter
	services.NewHealthReporter(
		ctx,
		logger,
		rabbitMq,
		messaging.RabbitMqExchange{
			Type:       "topic",
			Exchange:   cp.StringOrDefault("RABBITMQ_HEALTH_EXCHANGE", "fireops-edge-health"),
			RoutingKey: cp.StringOrDefault("RABBITMQ_HEALTH_ROUTING_KEY", ""),
		},
		SERVICE_NAME,
	)

	// Show run message
	logger.Info("Running...")

	// Wait until stop
	osSig := make(chan os.Signal, 1)
	signal.Notify(osSig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-osSig
	cancel()
	logger.Info("Shutting down...")
}
