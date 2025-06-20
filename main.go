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
			Exchange:   cp.StringOrDefault("RABBITMQ_EXCHANGE", "fireops-edge-alerts"),
			RoutingKey: cp.StringOrDefault("RABBITMQ_ROUTING_KEY", "new"),
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

	// Show run message
	logger.Info("Running...")

	// Wait until stop
	osSig := make(chan os.Signal, 1)
	signal.Notify(osSig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-osSig
	cancel()
	logger.Info("Shutting down...")
}
