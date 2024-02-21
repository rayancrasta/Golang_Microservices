package main

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const webPort = "80"

type Config struct {
	Rabbit *amqp.Connection
}

func main() {
	//Connect to rabbitmq
	rabbitConn, err := connect()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	defer rabbitConn.Close()

	app := Config{
		Rabbit: rabbitConn,
	}

	log.Printf("Starting broker service on port %s", webPort)

	// Define an HTTP server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	// Start the server
	err = srv.ListenAndServe()

	if err != nil {
		log.Panic(err)
	}

}

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
