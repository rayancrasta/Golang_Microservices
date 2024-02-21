package event

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	conn      *amqp.Connection
	queueName string
}

func NewConsumer(conn *amqp.Connection) (Consumer, error) {
	consumer := Consumer{
		conn: conn,
	}

	err := consumer.setup()
	if err != nil {
		return Consumer{}, err
	}

	return consumer, nil
}

func (consumer *Consumer) setup() error {
	channel, err := consumer.conn.Channel() // AMQ channel
	if err != nil {
		return err
	}

	return declareExchange(channel)
}

type Payload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (consumer *Consumer) Listen(topics []string) error {
	ch, err := consumer.conn.Channel()
	if err != nil {
		return err
	}

	defer ch.Close()

	//Get a random queue
	queue, err := declareRandomQueue(ch)
	if err != nil {
		return err
	}

	for _, topic := range topics {
		//Bind channel to each topics
		ch.QueueBind(
			queue.Name,
			topic,
			"logs_topic",
			false, // wait ?
			nil,   // arguments
		)
	}
	if err != nil {
		return err
	}

	// Look for messages
	messages, err := ch.Consume(queue.Name, "", true, false, false, false, nil)

	if err != nil {
		return err
	}

	// Consume till I exit from app
	forever := make(chan bool)

	go func() {
		for d := range messages {
			var payload Payload
			_ = json.Unmarshal(d.Body, &payload)

			go handlePayload(payload)

		}
	}()

	fmt.Printf("Waiting for message on [Exchang,Queue] [logs_topic,%s]\n", queue.Name)
	<-forever

	return nil

}

func handlePayload(payload Payload) {
	switch payload.Name {
	case "log", "event":
		//Log what we get
		err := logEvent(payload)
		if err != nil {
			log.Println(err)
		}

	case "auth":
		// Authentication

	default:
		err := logEvent(payload)
		if err != nil {
			log.Println(err)
		}
	}

}

func logEvent(entry Payload) error {
	log.Println("Inside Broker-logItem function")
	jsonData, _ := json.MarshalIndent(entry, "", "\t")
	logServiceURL := "http://logger-service/log"

	request, err := http.NewRequest("POST", logServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{}

	response, err := client.Do(request)

	if err != nil {
		return err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		return err
	}

	return nil

}
