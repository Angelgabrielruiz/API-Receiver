
package controllers

import (
    "fmt"
    "github.com/streadway/amqp"
    "Receive/src/pago/menssage/domain/repository"
)

type RabbitMQAdapter struct {
    conn *amqp.Connection
}


func NewRabbitMQAdapter(conn *amqp.Connection) repository.RabbitMQRepository {
    return &RabbitMQAdapter{conn: conn}
}

func (r *RabbitMQAdapter) PublishMessage(message string) error {
    channel, err := r.conn.Channel()
    if err != nil {
        return fmt.Errorf("Error al crear canal RabbitMQ: %w", err)
    }
    defer channel.Close()

    // Declarar la cola (puedes hacerla dinámica si gustas)
    queueName := "my_queue"
    _, err = channel.QueueDeclare(
        queueName,
        true,  // durable
        false, // autoDelete
        false, // exclusive
        false, // noWait
        nil,   // args
    )
    if err != nil {
        return fmt.Errorf("Error al declarar la cola: %w", err)
    }

    // Publicar el mensaje
    err = channel.Publish(
        "",
        queueName,
        false,
        false,
        amqp.Publishing{
            ContentType: "application/json",
            Body:        []byte(message),
        },
    )
    if err != nil {
        return fmt.Errorf("Error al publicar mensaje: %w", err)
    }

    return nil
}
