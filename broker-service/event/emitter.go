package event

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Emitter struct {
	connection *amqp.Connection
}

func (e *Emitter) setup() error {
	channel, err := e.connection.Channel()
	//a Channel is a communication link between a client (producer or consumer) and the RabbitMQ server. It is used to perform operations such as declaring queues, exchanges, publishing messages, and consuming messages.
	if err != nil {
		return err
	}

	defer channel.Close()

	return declareExchange(channel)
}

//Push message to the exchange
func (e *Emitter) Push(event string, severty string) error {
	channel, err := e.connection.Channel()

	if err != nil {
		return err
	}

	defer channel.Close()

	log.Println("Pushing to Channel")

	err = channel.Publish(
		"logs_topic", // exchange name
		severty,      // routing key
		false,        // mandatory (not used here, set to false)
		false,        // immediate (not used here, set to false)
		amqp.Publishing{
			ContentType: "text/plain",  // content type of the message
			Body:        []byte(event), // message body as a byte slice
		},
	)

	if err != nil {
		return err
	}

	return nil
}

// Call setup and return Emiiter struct ( amqp connection )
func NewEventEmitter(conn *amqp.Connection) (Emitter, error) {
	emitter := Emitter{
		connection: conn,
	}

	err := emitter.setup()
	if err != nil {
		return Emitter{}, err
	}

	return emitter, nil
}
