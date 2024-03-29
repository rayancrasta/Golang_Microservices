package main

import (
	"context"
	"fmt"
	"log"
	"log-service/data"
	"net"
	"net/http"
	"net/rpc"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	webPort  = "80"
	rpcPort  = "5001"
	mongoURL = "mongodb://mongo:27017"
	gRpcPort = "50001"
)

var client *mongo.Client

type Config struct {
	Models data.Models
}

func main() {
	// Connect to Mongo
	mongoClient, err := connectToMongo()
	if err != nil {
		log.Panic(err)
	}

	client = mongoClient

	// Create a context in order to disconnect
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)

	defer cancel() // cancel the context and release resources

	// Close connection. This will happen at end of the main fun execution
	defer func() {
		if err = client.Disconnect(ctx); err != nil { // enusres that disconnection is also bound by the 15sec timeout
			log.Panic(err)
		}
	}()

	app := Config{
		Models: data.New(client),
	}

	// Register the RPC Server

	err = rpc.Register(new(RPCServer))
	go app.rpcListen()

	go app.gRPCListen()
	// Start the web server

	// go app.serve()
	log.Println("Starting service on port: %s", webPort)
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err = srv.ListenAndServe()

	if err != nil {
		log.Panic(err)
	}

}

func (app *Config) rpcListen() error {
	log.Println("Starting RPC server on port", rpcPort)
	listen, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", rpcPort)) // Listen announces on the local network address

	if err != nil {
		return err
	}

	defer listen.Close()

	for {
		rpcConn, err := listen.Accept() // Accept waits for and returns the next connection to the listener.
		if err != nil {
			continue
		}

		go rpc.ServeConn(rpcConn) // use connection from above
	}

}

func connectToMongo() (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(mongoURL) // IP , Port of MongoDB
	clientOptions.SetAuth(options.Credential{
		Username: "admin",
		Password: "password",
	})

	// Connect
	connection, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Println("Error Connecting.. ", err)
		return nil, err
	}

	log.Println("Connected to Mongo: ")
	return connection, nil
	// The context.TODO() is used as the context parameter when establishing a connection to MongoDB. In this case, it indicates that there are no special context requirements for the MongoDB connection, and a basic background context is sufficient.

}
