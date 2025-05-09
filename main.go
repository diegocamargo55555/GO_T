package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Todo struct {
	ID        int    `json:"id" bson:"id"`
	Completed bool   `json:"complete"`
	Body      string `json:"body"`
}

var collection *mongo.Collection

func main() {
	println("hello world")

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error Loading .env file")
	}

	MONGODB_URI := os.Getenv("MONGODB_URI")
	clientOptions := options.Client().ApplyURI(MONGODB_URI)
	client, err := mongo.Connect(context.Background(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	defer client.Disconnect(context.Background())

	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MONGODB ATLAS")
	collection = client.Database("golang_db").Collection("todos")

	app := fiber.New()
	app.Get("/api/todos", getTodos)
	//app.Post("/api/todos", createTodo)
	//app.Patch("/api/todos/:id", updateTodo)
	//app.Delete("/api/todos/:id", deleteTodo)

	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	log.Fatal(app.Listen("0.0.0.0:" + port))
}

func getTodos(c *fiber.Ctx) error {
	var todos []Todo
	cursor, err := collection.Find(context.Background(), bson.M{})

	if err != nil {
		return err
	}

	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var todo Todo

		if err := cursor.Decode(&todo); err != nil {
			return err
		}
		todos = append(todos, todo)

	}

	return c.JSON(todos)
}

/*
func createTodo(c *fiber.Ctx) error {

}
func updateTodo(c *fiber.Ctx) error {

}
func deleteTodo(c *fiber.Ctx) error {

}
*/
