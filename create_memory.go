package main

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/go-redis/redis"
	"log"
	"member-when/internal/response"
	"net/http"
	"os"
)

type Memory struct {
	Content string `json:"memory"`
}

func main() {
	// Bit clear the date and time flags so they don't show up in CloudWatch (it logs the timestamp anyway).
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
	lambda.Start(createMemory)
}

// Accepts a message as input and adds it to a Redis sorted set.
// If the memory does not exist, then it is added to the set with a score of 1. Otherwise it is incremented by 1.
// Returns an APIGatewayProxyResponse with a JSON formatted message in the body and an HTTP Status Code header
func createMemory(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	// Unmarshal the request.
	memory := new(Memory)
	err := json.Unmarshal([]byte(request.Body), memory)
	if err != nil {
		log.Println("json unmarshal failure: " + err.Error())
		return response.Error(http.StatusInternalServerError, "json unmarshal failure: "+err.Error())
	}

	if memory.Content == "" {
		log.Println("Memory was empty.")
		return response.Error(http.StatusBadRequest, "Empty memory")
	}

	// Create a Redis client.
	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ENDPOINT") + ":" + os.Getenv("REDIS_PORT"),
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	// Add memory to the sorted set.
	resp, err := client.ZIncrBy("memories", 1, memory.Content).Result()

	if err != nil {
		log.Println("Redis ZADD failed: " + err.Error())
		return response.Error(http.StatusInternalServerError, "Redis ZADD failed: "+err.Error())
	}

	log.Printf("Memory added: "+memory.Content+", Score: %.f", resp)

	return response.Success()
}
