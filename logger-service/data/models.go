package data

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func New(mongo *mongo.Client) Models {
	client = mongo

	return Models{
		LogEntry: LogEntry{},
	}
}

type Models struct {
	LogEntry LogEntry
}

//Structure for all data stored in Mongo
type LogEntry struct {
	ID        string    `bson:"_id,omitempty" json:"id,omitempty"` // ?
	Name      string    `bson:"name" json:"name"`
	Data      string    `bson:"data" json:"data:`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

func (l *LogEntry) Insert(entry LogEntry) error {
	collection := client.Database("logs").Collection("logs") // Collection = Table

	_, err := collection.InsertOne(context.TODO(), LogEntry{
		Name:      entry.Name,
		Data:      entry.Data,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	if err != nil {
		log.Println("Error inserting into logs")
		return err
	}
	return nil
}

func (l *LogEntry) All() ([]*LogEntry, error) {
	// Make sure this doesnt execute forver
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)

	defer cancel()

	collection := client.Database("logs").Collection("logs")

	opts := options.Find()

	opts.SetSort(bson.D{{"created_at", -1}}) //sorts the results in descending order based on the "created_at" field.

	cursor, err := collection.Find(context.TODO(), bson.D{}, opts) //returs a cursor that can be iterated over to retrieve the results. The query specifies an empty filter (bson.D{}), meaning it fetches all documents from the collection.

	if err != nil {
		log.Println("Finding all docs error", err)
		return nil, err
	}

	defer cursor.Close(ctx)

	var logs []*LogEntry

	for cursor.Next(ctx) {
		var item LogEntry
		err := cursor.Decode(&item) //decodes each document into a LogEntry struct

		if err != nil {
			log.Println("Error decoding log into slice..", err)
			return nil, err
		} else {
			logs = append(logs, &item) // Append to the logs slice
		}
	}

	return logs, nil
}

func (l *LogEntry) getOne(id string) (*LogEntry, error) {
	// Make sure this doesnt execute forver
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)

	defer cancel()

	collection := client.Database("logs").Collection("logs")

	docID, err := primitive.ObjectIDFromHex(id) // convert the hexadecimal string id into a MongoDB ObjectID using the primitive.ObjectIDFromHex function.
	if err != nil {
		return nil, err
	}

	var entry LogEntry

	err = collection.FindOne(ctx, bson.M{"_id": docID}).Decode(&entry) // The result is then decoded into the entry

	if err != nil {
		return nil, err
	}

	return &entry, nil
}

func (l *LogEntry) DropCollection() error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)

	defer cancel()

	collection := client.Database("logs").Collection("logs")

	if err := collection.Drop(ctx); err != nil {
		return err
	}

	return nil
}

func (l *LogEntry) Update() (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)

	defer cancel()

	collection := client.Database("logs").Collection("logs")

	docID, err := primitive.ObjectIDFromHex(l.ID)

	if err != nil {
		return nil, err
	}

	result, err := collection.UpdateOne(
		ctx,
		bson.M{"_id": docID},
		bson.D{
			{"$set", bson.D{
				{"name", l.Name},
				{"data", l.Data},
				{"updated_at", time.Now()},
			}},
		},
	)

	if err != nil {
		return nil, err
	}

	return result, nil
}
