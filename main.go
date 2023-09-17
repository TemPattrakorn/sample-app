package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		log.Fatal("You must set your 'MONGODB_URI' environment variable. See\n\t https://www.mongodb.com/docs/drivers/go/current/usage-examples/#environment-variable")
	}
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	fmt.Println("Connected to MongoDB")
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	// HTTP REQUEST LISTENER
	http.HandleFunc("/post", formSubmitHandler(client))
	err2 := http.ListenAndServe(":8080", nil)
	if err2 != nil {
		panic(err2)
	}

}

type MyDocument struct {
	ItemName    string `bson:"item-name"`
	Faculty     string `bson:"faculty"`
	TeacherName string `bson:"teacher-name"`
}

// insert
func insertDocument(client *mongo.Client, document MyDocument) error {
	// Get a handle to the target database and collection
	collection := client.Database("senior_project").Collection("durable_goods")

	// Insert the document
	_, err := collection.InsertOne(context.Background(), document)
	return err
}

/*
// update
func UpdateItem(client *mongo.Client, databaseName, collectionName string, itemID int, updatedItem Item) error {
	// Get a handle to the MongoDB collection.
	collection := client.Database(databaseName).Collection(collectionName)

	// Define the filter to find the item you want to update.
	filter := bson.M{"_id": itemID}

	// Define the update operation.
	update := bson.M{
		"$set": updatedItem, // Use $set to update specific fields.
	}

	// Specify options for the update operation (optional).
	options := options.Update().SetUpsert(false)

	// Perform the update operation.
	_, err := collection.UpdateOne(context.TODO(), filter, update, options)
	if err != nil {
		return err
	}

	return nil
}
*/

//find
/*
func
*/

//delete

func formSubmitHandler(client *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			var doc MyDocument
			err := r.ParseForm()
			// Read the request body
			body, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "Error reading request body", http.StatusInternalServerError)
				return
			}

			doc.ItemName = r.FormValue("item-name")
			doc.Faculty = r.FormValue("faculty")
			doc.TeacherName = r.FormValue("teacher-name")
			// INSERT THE DOCUMENT TO MONGODB
			err = insertDocument(client, doc)
			// Print the received data to the server console
			fmt.Println("Received data:", string(body))

			// Respond with a confirmation message
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Data received successfully"))
		} else {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		}
	}
}
