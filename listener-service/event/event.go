package event

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

func declareExchange(ch *amqp.Channel) error {
	//The exchange is responsible for routing the message to one or more queues based on rules defined by the exchange type.
	return ch.ExchangeDeclare(
		"logs_topic", //name
		"topic",      // type
		true,         // is it durable ?
		false,        // auto-deleted ?
		false,        // internal ?
		false,        // no-wait
		nil,          // arguments
	)
}

func declareRandomQueue(ch *amqp.Channel) (amqp.Queue, error) {
	return ch.QueueDeclare(
		"",    //name ?
		false, // durable ?
		false, // autodelete when unused?
		true,  // exclusive channel for current operations
		false, //no-wait?
		nil,   // arguments?
	)
}
