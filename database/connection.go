package database

import (
	// "context"
	"fmt"
	"os"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var DB *mongo.Database

func InitDb() {
	var err error
	client, err := mongo.Connect(options.Client().
		ApplyURI(os.Getenv("MONGO_URI")))
	if err != nil {
		panic(err)
	}
	DB = client.Database(os.Getenv("MONGO_DB_NAME"))
	// defer func() {
	// 	if err := client.Disconnect(context.TODO()); err != nil {
	// 		panic(err)
	// 	}
	// }()
	fmt.Println("Connect Database")
}
