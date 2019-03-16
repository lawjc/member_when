package response

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"log"
	"net/http"
)

// Wrapper used to marshal a response and add it to an APIGatewayProxyResponse
type apiResponse struct {
	Message []string `json:"message"`
}

type errorResponse struct {
	Message string `json:"message"`
}

// Takes an HTTP status code, message, and error as input, then returns
// an APIGatewayProxyResponse and error.
// Marshals the message into the APIGatewayProxyResponse so it can be consumed by the client as JSON.
func ApiResponse(statusCode int, message []string, err error) (events.APIGatewayProxyResponse, error) {

	js, err := json.Marshal(apiResponse{message})
	if err != nil {
		log.Fatal("marshalling API response: " + err.Error())
	}

	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Body:       string(js),
	}, err
}

func Error(statusCode int, message string) (events.APIGatewayProxyResponse, error) {

	js, err := json.Marshal(errorResponse{message})
	if err != nil {
		log.Fatal("marshalling API response: " + err.Error())
	}

	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Body:       string(js),
	}, err
}

func Success() (events.APIGatewayProxyResponse, error) {

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
	}, nil
}
