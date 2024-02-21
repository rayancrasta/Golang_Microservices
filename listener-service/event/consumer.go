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

// Take the connection string and call setup
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

//RabbitMQ-Exchange values are set in declareExchange func call ,channel values are set
func (consumer *Consumer) setup() error {
	channel, err := consumer.conn.Channel() // AMQ channel
	if err != nil {
		return err
	}

	return declareExchange(channel) //
}

type Payload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

//Declare a random queue, do the config of queue with an exchange, look for messages in that queue forever, based called handlepayload to take action based on the payload.Name ( catgeory: auth,log etc)
func (consumer *Consumer) Listen(topics []string) error {
	ch, err := consumer.conn.Channel()
	if err != nil {
		return err
	}

	defer ch.Close()

	//Get a random queue, currently no name to Queue
	queue, err := declareRandomQueue(ch)
	if err != nil {
		return err
	}

	for _, topic := range topics {
		//Bind channel to each topics
		ch.QueueBind(
			queue.Name,   // queue name
			topic,        // routing key
			"logs_topic", // exchange name
			false,        // wait ?
			nil,          // arguments
		)

		//QueueBind method is used to bind a queue to an exchange with a specific routing key.
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
			_ = json.Unmarshal(d.Body, &payload) //get the data from the queue and put in payload

			go handlePayload(payload)

		}
	}()

	fmt.Printf("Waiting for message on [Exchang,Queue] [logs_topic,%s]\n", queue.Name)
	<-forever

	//<-forever is a blocking operation that waits for a value to be sent on the forever channel. Since there's no explicit sender for the channel in this code, it essentially blocks the program from exiting, creating an infinite loop.

	return nil

}

//Based on payloadName ,take the payload and do the necessary
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

//Call the logger-service
func logEvent(entry Payload) error {
	log.Println("Inside Listener-logEvent function")
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
