# FireOps Edge SevenIO Notifier
This repository contains a service, that notifies about new incomming alerts through sms, that are send via [SevenIO](https://www.seven.io/). SevenIO is an online SMS Gateway, that is used for sending alert messages.

## Configuration
All configuration is done via environmental variables because the intended form of running the project is in a Docker container.

| Variable | Default | Description |
|----------|---------|-------------|
| RABBITMQ_HOST |  | RabbitMQ host (e.g. 192.168.x.x or Hostname) |
| RABBITMQ_PORT | 5672 | RabbitMQ port |
| RABBITMQ_USER |  | RabbitMQ user |
| RABBITMQ_PW |  | RabbitMQ password |
| RABBITMQ_EXCHANGE | fireops-edge-events | RabbitMQ Exchange, where alerts will be published |
| RABBITMQ_ROUTING_KEY | alu2g.new | RabbitMQ routing key for all changes on currently active alerts |
| RABBITMQ_HEALTH_EXCHANGE | fireops-edge-health | RabbitMQ exchange for health messages |
| RABBITMQ_HEALTH_ROUTING_KEY |  | RabbitMQ routing key for health messages |
||||
| SEVENIO_API_URL |  | SevenIO base url (e.g. https://gateway.seven.io) |
| SEVENIO_TOKEN |  | SevenIO API Key (from developer console) |
||||
| SMS_ALERT_TEXT | EINSATZ von FireOps | Alert text, that will be send as sms text (body) |
| SMS_GROUP_DEFAULT | Mannschaft | Contact group of SevenIO contacts, that will receive default message data |
| SMS_GROUP_EXTENDED | Kommando | Contact group of SevenIO contacts, that will receive extended message data |
||||
| LOG_LEVEL | INFO | TRACE, DEBUG, INFO, WARNING, ERROR, FATAL, OFF |
