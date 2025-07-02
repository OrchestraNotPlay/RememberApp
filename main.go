package main

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"github.com/gofiber/fiber/v2"
)

type Remember struct {
	ID          bson.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Text        string        `json:"text"`
	Date        string        `json:"date"`
	ExpiredTime string        `json:"expiredtime"`
}

var collection *mongo.Collection

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal(err)
	}

	MONGODB_URI := os.Getenv("MONGODB_URI")

	clientOptions := options.Client().ApplyURI(MONGODB_URI)
	client, err := mongo.Connect(clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	collection = client.Database("RememberApp").Collection("Products")

	app := fiber.New()

	app.Get("/api/getRemembers", getRemembers)
	app.Post("/api/addRemember", addRemember)
	app.Delete("/api/deleteRemember/:id", deleteRemember)
	app.Patch("/api/updateRemember/:id", updateRemembers)

	PORT := os.Getenv("PORT")
	log.Fatal(app.Listen("0.0.0.0:" + PORT))
}

func getRemembers(c *fiber.Ctx) error {
	var remembers []Remember

	cursor, err := collection.Find(context.Background(), bson.M{})

	if err != nil {
		log.Fatal(err)
	}

	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var remember Remember

		if err := cursor.Decode(&remember); err != nil {
			log.Fatal(err)
		}

		remembers = append(remembers, remember)
	}

	return c.JSON(remembers)
}

func addRemember(c *fiber.Ctx) error {
	remember := new(Remember)

	if err := c.BodyParser(remember); err != nil {
		return err
	}

	insertResult, err := collection.InsertOne(context.Background(), remember)
	if err != nil {
		return err
	}

	remember.ID = insertResult.InsertedID.(bson.ObjectID)

	return c.Status(201).JSON(remember)
}

func deleteRemember(c *fiber.Ctx) error {
	id := c.Params("id")

	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		log.Fatal(err)
	}

	collection.DeleteOne(context.Background(), bson.M{"_id": objectID})

	return c.Status(200).JSON(fiber.Map{"success": true})
}

func updateRemembers(c *fiber.Ctx) error {
	id := c.Params("id")
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		log.Fatal(err)
	}

	remember := new(Remember)

	if err := c.BodyParser(remember); err != nil {
		log.Fatal(err)
	}

	filter := bson.M{"_id": objectID}
	update := bson.M{"$set": bson.M{"text": remember.Text, "date": remember.Date, "expiredtime": remember.ExpiredTime}}

	_, err = collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Fatal(err)
	}

	return c.Status(200).JSON(fiber.Map{"success": true})

}
