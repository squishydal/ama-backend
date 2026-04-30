package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type ama struct {
	ID         bson.ObjectID `bson:"_id,omitempty" json:"id"`
	AskedBy    string        `bson:"askedBy"       json:"askedBy"`
	Question   string        `bson:"question"      json:"question"`
	Answer     string        `bson:"answer"        json:"answer"`
	AnsweredAt time.Time     `bson:"answeredAt"    json:"answeredAt"`
}

var collection *mongo.Collection

func main() {
	godotenv.Load()

	uri := os.Getenv("MONGODB_URI")
	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal("Failed to connect MongoDB:", err)
	}
	defer client.Disconnect(context.TODO())

	if err := client.Ping(context.TODO(), nil); err != nil {
		log.Fatal("MongoDB ping failed:", err)
	}
	log.Println("Connected to MongoDB!")

	collection = client.Database("portofolio").Collection("questions")

	router := gin.Default()
	router.Use(cors.Default())
	router.GET("/api/questions", getAmas)
	router.POST("/api/questions", postAmas)
	if err := router.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func getAmas(c *gin.Context) {
	filter := bson.M{"answer": bson.M{"$ne": ""}}
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	results := []ama{} // ← ini yang bikin [] bukan null
	if err := cursor.All(context.TODO(), &results); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, results)
}

func postAmas(c *gin.Context) {
	var newAma ama
	if err := c.BindJSON(&newAma); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newAma.ID = bson.NewObjectID()
	newAma.AnsweredAt = time.Now()

	_, err := collection.InsertOne(context.TODO(), newAma)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusCreated, newAma)
}
