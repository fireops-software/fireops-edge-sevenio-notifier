package services

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/fireops-software/fireops-edge-sevenio-notifier/api"
	"github.com/fireops-software/fireops-edge-sevenio-notifier/domain"
	"github.com/rabbitmq/amqp091-go"
	"github.com/uoul/go-common/log"
	"github.com/uoul/go-common/messaging"
)

type SevenIoNotifier struct {
	messenger messaging.IMessenger[messaging.RabbitMqExchange, amqp091.Delivery]
	ctx       context.Context
	logger    log.ILogger

	exchange   messaging.RabbitMqExchange
	sevenIoApi api.ISevenIoApi

	sender           string
	groupFullInfo    string
	groupDefaultInfo string
	alertText        string
}

func (s *SevenIoNotifier) run() {
	// Subscribe for new alerts
	sub := s.messenger.Subscribe(s.exchange)
	defer s.messenger.Unsubscribe(sub)
	for {
		select {
		case <-s.ctx.Done():
			return
		case msg := <-sub:
			s.logger.Debugf("New incomming message from rabbitmq: %v", msg.Result.Body)
			if msg.Error != nil {
				s.logger.Errorf("failed to receive message from messenger - %v", msg.Error)
				continue
			}
			events := []domain.Event{}
			err := json.Unmarshal(msg.Result.Body, &events)
			if err != nil {
				s.logger.Errorf("failed to parse incomming message as json - %v", err)
				continue
			}
			if err = s.sendSms(events); err != nil {
				s.logger.Errorf("failed to send sms - %v", err)
				continue
			}
		}
	}

}

func (s *SevenIoNotifier) sendSms(events []domain.Event) error {
	// Send sms with default information
	s.logger.Debugf("Send sms to default group(%s)...", s.groupDefaultInfo)
	r := <-s.sevenIoApi.SendSms(s.sender, s.groupDefaultInfo, s.alertText)
	if r.Error != nil {
		return r.Error
	}
	// Send sms with full information
	s.logger.Debugf("Send sms to extended group(%s)...", s.groupFullInfo)
	r = <-s.sevenIoApi.SendSms(s.sender, s.groupFullInfo, createAlertInfoText(events))
	if r.Error != nil {
		return r.Error
	}
	return nil
}

func createAlertInfoText(events []domain.Event) string {
	var sb strings.Builder
	templateStr := "Nr.: %s\nArt: %s\nAnrufer: %s\nTel.: %s\nOrt: %s\nInfo: %s\n"
	for _, e := range events {
		sb.WriteString(fmt.Sprintf(
			templateStr,
			valueOrDefault(e.Num1),
			valueOrDefault(e.Category),
			valueOrDefault(e.CallerName),
			valueOrDefault(e.CallerNumber),
			valueOrDefault(e.Location),
			valueOrDefault(e.EventAlarmtext),
		))
	}
	return sb.String()
}

func WithSevenIoSender(sender string) func(*SevenIoNotifier) {
	return func(sin *SevenIoNotifier) {
		sin.sender = sender
	}
}

func WithSevenIoGroupDefault(defaultGroup string) func(*SevenIoNotifier) {
	return func(sin *SevenIoNotifier) {
		sin.groupDefaultInfo = defaultGroup
	}
}

func WithSevenIoGroupFull(fullGroup string) func(*SevenIoNotifier) {
	return func(sin *SevenIoNotifier) {
		sin.groupFullInfo = fullGroup
	}
}

func WithSevenIoAlertText(alertText string) func(*SevenIoNotifier) {
	return func(sin *SevenIoNotifier) {
		sin.alertText = alertText
	}
}

func valueOrDefault[T any](v *T) T {
	if v != nil {
		return *v
	}
	return *new(T)
}

func NewSevenIoNotifier(
	ctx context.Context,
	logger log.ILogger,
	messenger messaging.IMessenger[messaging.RabbitMqExchange, amqp091.Delivery],
	exchange messaging.RabbitMqExchange,
	sevenIoApi api.ISevenIoApi,

	opts ...func(*SevenIoNotifier)) *SevenIoNotifier {

	notifier := &SevenIoNotifier{
		ctx:        ctx,
		logger:     logger,
		messenger:  messenger,
		exchange:   exchange,
		sevenIoApi: sevenIoApi,

		sender:           "FireOps",
		groupFullInfo:    "Kommando",
		groupDefaultInfo: "Mannschaft",
		alertText:        "EINSATZBEFEHL",
	}
	for _, o := range opts {
		o(notifier)
	}
	go notifier.run()
	return notifier
}
