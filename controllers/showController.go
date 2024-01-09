package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"siencefestival/api/models"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/streadway/amqp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var showCollection *mongo.Collection

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	connectionString := os.Getenv("MONGODB_CONNECTION_STRING")
	clientOptions := options.Client().ApplyURI(connectionString)

	client, error := mongo.Connect(context.TODO(), clientOptions)

	if error != nil {
		log.Fatal(error)
	}

	fmt.Println("Mongodb connection success")

	dbName := os.Getenv("DBNAME")
    colName := os.Getenv("COLNAME")

	showCollection = client.Database(dbName).Collection(colName)

	fmt.Println("Collection istance is ready")
}

func AddShow(c *gin.Context) {
	var ctx, cancel = context.WithCancel(context.Background())
	defer cancel()

	var show models.Show

	if err := c.BindJSON(&show); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var existingShow models.Show
    err := showCollection.FindOne(ctx, bson.M{"release_date": show.ReleaseDate}).Decode(&existingShow)
    if err != mongo.ErrNoDocuments {
        c.JSON(http.StatusBadRequest, gin.H{"error": "A show with the same date and time already exists"})
        return
    }

	show.Accepted = false

	_, insertErr := showCollection.InsertOne(ctx, show)

	if insertErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while inserting new show"})
		return
	}

	c.JSON(http.StatusOK, show)

}

//get all shows
func GetAllShows(c *gin.Context) {
    var ctx, cancel = context.WithCancel(context.Background())
    defer cancel()

    var shows []models.Show

    cursor, err := showCollection.Find(ctx, bson.M{})
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while getting all shows"})
        return
    }
    defer cursor.Close(ctx)

    for cursor.Next(ctx) {
        var show models.Show
        cursor.Decode(&show)
        shows = append(shows, show)
    }

    if err := cursor.Err(); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while getting all shows"})
        return
    }

    c.JSON(http.StatusOK, shows)
}

//get show by id
func GetShowById(c *gin.Context) {
    var ctx, cancel = context.WithCancel(context.Background())
    defer cancel()

    id := c.Param("id")

    objectId, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid id format"})
        return
    }

    var show models.Show
    err = showCollection.FindOne(ctx, bson.M{"_id": objectId}).Decode(&show)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while getting show by id"})
        return
    }

    c.JSON(http.StatusOK, show)
}

func GetAcceptedShows(c *gin.Context) {
    var ctx, cancel = context.WithCancel(context.Background())
    defer cancel()

    var shows []models.Show

    cursor, err := showCollection.Find(ctx, bson.M{"accepted": true})
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while getting accepted shows"})
        return
    }
    defer cursor.Close(ctx)

    for cursor.Next(ctx) {
        var show models.Show
        cursor.Decode(&show)
        shows = append(shows, show)
    }

    if err := cursor.Err(); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while getting accepted shows"})
        return
    }

    c.JSON(http.StatusOK, shows)
}

func GetUnacceptedShows(c *gin.Context) {
	var ctx, cancel = context.WithCancel(context.Background())
	defer cancel()
	
	var shows []models.Show
	
	cursor, err := showCollection.Find(ctx, bson.M{"accepted": nil})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while getting accepted shows"})
		return
	}
	defer cursor.Close(ctx)
	
	for cursor.Next(ctx) {
		var show models.Show
		cursor.Decode(&show)
		shows = append(shows, show)
	}
	
	if err := cursor.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while getting accepted shows"})
		return
	}
	
	c.JSON(http.StatusOK, shows)
}

func StartMessageConsumer() {
    rabbitMqConnectionString := os.Getenv("RABBITMQ_CONNECTION_STRING")
    conn, err := amqp.Dial(rabbitMqConnectionString)
    if err != nil {
        log.Fatalf("Failed to connect to RabbitMQ: %v", err)
    }
    defer conn.Close()

    ch, err := conn.Channel()
    if err != nil {
        log.Fatalf("Failed to open a channel: %v", err)
    }
    defer ch.Close()

    queue, err := ch.QueueDeclare(
        "reviews",
        false,
        false,
        false,
        false,
        nil,
    )
    if err != nil {
        log.Fatalf("Failed to declare a queue: %v", err)
    }

    msgs, err := ch.Consume(
        queue.Name,
        "",
        true,
        false,
        false,
        false,
        nil,
    )
    if err != nil {
        log.Fatalf("Failed to register a consumer: %v", err)
    }

    for msg := range msgs {
        log.Printf("Received a message: %s", msg.Body)

        var message models.Message
        if err := json.Unmarshal(msg.Body, &message); err != nil {
            log.Printf("Error decoding message: %v", err)
            continue
        }

        showID := message.ShowId

        updateShowAccepted(showID)
    }
}

func updateShowAccepted(showID string) {
	objID, err := primitive.ObjectIDFromHex(showID)
    if err != nil {
        log.Printf("Could not convert showID to ObjectID: %v", err)
        return
    }

    filter := bson.M{"_id": objID, "accepted": bson.M{"$ne": true}}

    update := bson.M{
        "$set": bson.M{
            "accepted": true,
        },
    }

    updateResult, err := showCollection.UpdateOne(context.TODO(), filter, update)
    if err != nil {
        log.Printf("Could not update show: %v", err)
        return
    }

    if updateResult.MatchedCount == 0 {
        log.Println("No show found with ID:", showID)
        return
    }


    log.Println("Updated show with ID:", showID)
}