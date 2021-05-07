package main

import (
	"encoding/json"
	"errors"
	"os"
	"strconv"

	"github.com/ChimeraCoder/anaconda"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/kotaroooo0/kagura/snow_forecaster"
)

type Pair struct {
	First  string `json:"first"`
	Second string `json:"second"`
}

func schedule(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var p Pair
	if err := json.Unmarshal([]byte(req.Body), &p); err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	if p.First == "" || p.Second == "" {
		return events.APIGatewayProxyResponse{}, errors.New("two elements are needed")
	}

	c, err := tweetContent(p)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	anaconda.SetConsumerKey(os.Getenv("CONSUMER_KEY"))
	anaconda.SetConsumerSecret(os.Getenv("CONSUMER_SECRET"))
	t := anaconda.NewTwitterApi(os.Getenv("ACCESS_TOKEN_KEY"), os.Getenv("ACCESS_TOKEN_SECRET"))

	if _, err := t.PostTweet(c, nil); err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
	}, nil
}

func tweetContent(p Pair) (string, error) {
	sf := snow_forecaster.NewSnowForecaster()
	firstData, err := sf.Forecast(p.First)
	if err != nil {
		return "", err
	}
	secondData, err := sf.Forecast(p.Second)
	if err != nil {
		return "", err
	}
	content := "今日 | 明日 | 明後日 (朝,昼,夜)\n"
	content += p.First + "\n"
	content += areaLineString(firstData) + "\n"
	content += p.Second + "\n"
	content += areaLineString(secondData) + "\n"
	return content, nil
}

func areaLineString(snowfallForecast snow_forecaster.Forecast) string {
	content := strconv.Itoa(snowfallForecast.Snows[0].Morning) + addRainyChar(snowfallForecast.Rains[0].Morning) + ", " + strconv.Itoa(snowfallForecast.Snows[0].Noon) + addRainyChar(snowfallForecast.Rains[0].Noon) + ", " + strconv.Itoa(snowfallForecast.Snows[0].Night) + addRainyChar(snowfallForecast.Rains[0].Night) + "cm | "
	content += strconv.Itoa(snowfallForecast.Snows[1].Morning) + addRainyChar(snowfallForecast.Rains[1].Morning) + ", " + strconv.Itoa(snowfallForecast.Snows[1].Noon) + addRainyChar(snowfallForecast.Rains[1].Noon) + ", " + strconv.Itoa(snowfallForecast.Snows[1].Night) + addRainyChar(snowfallForecast.Rains[1].Night) + "cm | "
	content += strconv.Itoa(snowfallForecast.Snows[2].Morning) + addRainyChar(snowfallForecast.Rains[2].Morning) + ", " + strconv.Itoa(snowfallForecast.Snows[2].Noon) + addRainyChar(snowfallForecast.Rains[2].Noon) + ", " + strconv.Itoa(snowfallForecast.Snows[2].Night) + addRainyChar(snowfallForecast.Rains[2].Night) + "cm "
	return content
}

func addRainyChar(rainfall int) string {
	if rainfall > 5 {
		return "☔️"
	}
	if rainfall > 0 {
		return "☂️"
	}
	return ""
}

func main() {
	lambda.Start(schedule)
}
