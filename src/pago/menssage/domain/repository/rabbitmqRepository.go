package repository

type RabbitMQRepository interface {
	PublishMessage(message string) error
}