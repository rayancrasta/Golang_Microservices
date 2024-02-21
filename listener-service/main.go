package main

import (
	"listener/event"
	"log"
	"math"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	//Connect to rabbitmq
	rabbitConn, err := connect()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	defer rabbitConn.Close()

	//Start listening for messages
	log.Println("Listening for and consuming RabbitMQ messages...")

	//Create a consumer
	consumer, err := event.NewConsumer(rabbitConn)
	// consumer is aa struct of amqp connection and queue name

	if err != nil {
		panic(err)
	}

	//Watch the queue and consume events
	err = consumer.Listen([]string{"log.INFO", "log.WARNING", "log.ERROR"}) // topic names
	// because we call getRandomQueue func once in Listen, we use the same queue for all
	if err != nil {
		log.Println(err)
	}

}

// Connect to RabbitMQ, dial to it, retries and backoff
func connect() (*amqp.Connection, error) {

	//Back off in N times
	var count int64
	var backoff = 1 * time.Second
	var connection *amqp.Connection

	// Dont continue till rabbitmq is ready
	for {
		c, err := amqp.Dial("amqp://guest:guest@rabbitmq")
		if err != nil {
			log.Println("Rabbitmq not ready yet", err)
			count++
		} else {
			log.Println("Connected to RabbitMQ")
			connection = c
			break
		}

		if count > 5 {
			log.Println(err)
			return nil, err
		}

		backoff = time.Duration(math.Pow(float64(count), 2)) * time.Second
		log.Println("Backing off")
		time.Sleep(backoff)
		continue
	}

	return connection, nil
}
