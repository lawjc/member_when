package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/go-redis/redis"
	"log"
	"member-when/internal/response"
	"net/http"
	"os"
)

func main() {
	// Bit clear the date and time flags so they don't show up in CloudWatch (it logs this information already).
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
	log.Println("Request received")
	lambda.Start(listMemories)
}

// Returns a sorted list of the top 10 memories from Redis.
// If there are no entries, then an empty list is returned.
//
// The elements are considered to be ordered from the highest to the lowest score.
// Descending lexicographical order is used for elements with equal score.
func listMemories(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	log.Printf("Creating Redis client with: %s:%s", os.Getenv("REDIS_ENDPOINT"), os.Getenv("REDIS_PORT"))

	// Create a Redis client.
	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ENDPOINT") + ":" + os.Getenv("REDIS_PORT"),
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	resp, err := client.ZRevRange("memories", 0, 9).Result()
	if err != nil {
		log.Println("Redis ZREVRANGE failed: " + err.Error())
		return response.Error(http.StatusInternalServerError, "Redis ZREVRANGE failed: "+err.Error())
	}

	log.Println("Returning list")

	return response.ApiResponse(http.StatusOK, resp, nil)
}
