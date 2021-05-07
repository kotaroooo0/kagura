package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func get(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	ct, ok := req.QueryStringParameters["crc_token"]
	if !ok {
		return events.APIGatewayProxyResponse{}, errors.New("query paramater `crc_token` is not found")
	}

	mac := hmac.New(sha256.New, []byte(os.Getenv("CONSUMER_SECRET")))
	mac.Write([]byte(ct))
	token := "sha256=" + base64.StdEncoding.EncodeToString(mac.Sum(nil))

	return events.APIGatewayProxyResponse{
		Body:       fmt.Sprintf(`{"token": "%s"}`, token),
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(get)
}
