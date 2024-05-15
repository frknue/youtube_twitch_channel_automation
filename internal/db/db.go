package db

import (
	"context"
	"fmt"
	badger "github.com/dgraph-io/badger/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

func SaveClipID(clipID string) error {
	// Set client options and connect to MongoDB
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.TODO())

	// Connect to the database and collection
	collection := client.Database("youtube_twitch_channel_automation").Collection("clips")

	// Check if the clipID already exists
	var result bson.M
	err = collection.FindOne(context.TODO(), bson.M{"clipID": clipID}).Decode(&result)
	if err == mongo.ErrNoDocuments {
		// Document not found, insert new document
		_, err := collection.InsertOne(context.TODO(), bson.M{
			"clipID": clipID,
			"status": "processed",
			"time":   time.Now(),
		})
		if err != nil {
			return err
		}
		log.Println("Clip ID processed and stored")
	} else if err != nil {
		return err
	} else {
		log.Println("Clip ID already processed")
	}

	return nil
}

// Check if the clip ID is already processed
func CheckClipID(clipID string) bool {
	// Set client options and connect to MongoDB
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Printf("Error connecting to MongoDB: %v", err)
		return false
	}
	defer client.Disconnect(context.TODO())

	// Connect to the database and collection
	collection := client.Database("youtube_twitch_channel_automation").Collection("clips")

	// Check if the clipID already exists
	var result bson.M
	err = collection.FindOne(context.TODO(), bson.M{"clipID": clipID}).Decode(&result)
	if err == mongo.ErrNoDocuments {
		log.Println("Clip ID not found")
		return false // Clip ID not found
	} else if err != nil {
		log.Printf("Error checking clip ID: %v", err)
		return false // An error occurred during the operation
	}

	log.Println("Clip ID found")
	return true // Clip ID exists
}

// SaveVideoData saves the Run data to the "video" collection in MongoDB
func SaveVideoData(video interface{}) error {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatalf("Failed to get MongoDB client: %v", err)
		return err
	}
	defer client.Disconnect(context.Background())

	// Get a handle for your collection
	collection := client.Database("youtube_twitch_channel_automation").Collection("video")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Inserting the video data into the video collection
	result, err := collection.InsertOne(ctx, video)
	if err != nil {
		log.Printf("Could not insert video into video collection: %v", err)
		return err
	}

	fmt.Printf("Inserted document with ID: %v\n", result.InsertedID)
	return nil
}

func GetLatestEpisodeByGameID(gameID string) (int, error) {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatalf("Failed to get MongoDB client: %v", err)
		return 0, err
	}
	defer client.Disconnect(context.Background())

	collection := client.Database("youtube_twitch_channel_automation").Collection("video")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Define filter query for fetching the specific document
	filter := bson.D{{"gameid", gameID}} // Make sure the field name is exactly as in MongoDB

	// Define the sorting order - descending by 'videoepisode'
	options := options.FindOne().SetSort(bson.D{{"videoepisode", -1}})

	// Find a single document from the video collection with the highest 'videoepisode'
	var result bson.M
	err = collection.FindOne(ctx, filter, options).Decode(&result)
	if err != nil {
		log.Printf("Could not find document: %v", err)
		return 0, err
	}

	episodeNumber, ok := result["videoepisode"].(int32)
	if !ok {
		log.Println("Failed to assert episode as int32")
		return 0, fmt.Errorf("episode field is not of type int32")
	}
	episode := int(episodeNumber)

	return episode, nil
}

// GetLatestVideoByGameID fetches the latest video document by game ID from the MongoDB collection
func GetLatestVideoByGameID(gameID string) (bson.M, error) {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatalf("Failed to get MongoDB client: %v", err)
		return nil, err
	}
	defer client.Disconnect(context.Background())

	collection := client.Database("youtube_twitch_channel_automation").Collection("video")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Define filter query for fetching the specific document
	filter := bson.D{{"gameid", gameID}}

	// Find a single document from the video collection
	var result bson.M
	err = collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		log.Printf("Could not find document: %v", err)
		return nil, err
	}

	fmt.Println("Found document:", result)
	return result, nil
}

// Create a lock to prevent multiple instances of the application from running concurrently
func CreateLock() error {
	// Open the Badger database located in the specified directory
	db, err := badger.Open(badger.DefaultOptions("/tmp/badger"))
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// Use Update function for transactional operation
	err = db.Update(func(txn *badger.Txn) error {
		_, err := txn.Get([]byte("lock"))
		if err == badger.ErrKeyNotFound {
			// Key not found, create the lock
			e := badger.NewEntry([]byte("lock"), []byte("locked"))
			err = txn.SetEntry(e)
			if err != nil {
				return err
			}
			log.Println("Lock created")
		} else if err != nil {
			return err
		} else {
			log.Println("Lock already exists")
		}
		return nil
	})

	if err != nil {
		log.Fatalf("Failed to create lock: %v", err)
	}

	return nil
}

// Remove the lock to allow other instances of the application to run
func RemoveLock() error {
	// Open the Badger database located in the specified directory
	db, err := badger.Open(badger.DefaultOptions("/tmp/badger"))
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// Use Update function for transactional operation
	err = db.Update(func(txn *badger.Txn) error {
		err := txn.Delete([]byte("lock"))
		if err != nil {
			return err
		}
		log.Println("Lock removed")
		return nil
	})

	if err != nil {
		log.Fatalf("Failed to remove lock: %v", err)
	}

	return nil
}
