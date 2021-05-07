package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/ChimeraCoder/anaconda"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/kotaroooo0/kagura/similar_searcher"
	"github.com/kotaroooo0/kagura/snow_forecaster"
	"golang.org/x/exp/utf8string"
)

type PostTwitterWebhookRequest struct {
	UserID            string             `json:"for_user_id"`
	TweetCreateEvents []TweetCreateEvent `json:"tweet_create_events"`
}

type TweetCreateEvent struct {
	TweetID    int64  `json:"id"`
	TweetIDStr string `json:"id_str"`
	Text       string `json:"text"`
	User       struct {
		UserID     int64  `json:"id"`
		IDStr      string `json:"id_str"`
		ScreenName string `json:"screen_name"`
	} `json:"user"`
}

func post(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var twitterWebhookRequest PostTwitterWebhookRequest

	if err := json.Unmarshal([]byte(req.Body), &twitterWebhookRequest); err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	if len(twitterWebhookRequest.TweetCreateEvents) < 1 {
		return events.APIGatewayProxyResponse{}, errors.New("not found reply or invalid user")
	}

	tweetText := twitterWebhookRequest.TweetCreateEvents[0].Text

	// @snowfall_botを消す
	replyText := strings.Replace(tweetText, "@snowfall_bot ", "", -1)

	ss := similar_searcher.NewSimilarSearcher()
	similar, err := ss.Search(replyText)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	key := toKey(similar)
	sf := snow_forecaster.NewSnowForecaster()
	f, err := sf.Forecast(key)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	// TODO
	content, err := replyContent("hoge", f)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	params := url.Values{}
	params.Set("in_reply_to_status_id", twitterWebhookRequest.TweetCreateEvents[0].TweetIDStr)

	anaconda.SetConsumerKey(os.Getenv("CONSUMER_KEY"))
	anaconda.SetConsumerSecret(os.Getenv("CONSUMER_SECRET"))
	c := anaconda.NewTwitterApi(os.Getenv("ACCESS_TOKEN_KEY"), os.Getenv("ACCESS_TOKEN_SECRET"))

	if _, err := c.PostTweet(fmt.Sprintf("@%s %s", twitterWebhookRequest.TweetCreateEvents[0].User.ScreenName, content), params); err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
	}, nil
}

func toKey(str string) string {
	// TBD
	return ""
}

func replyContent(name string, sf snow_forecaster.Forecast) (string, error) {
	// TODO: 仮の文章
	content := name + "\n"
	content += "今日 | 明日 | 明後日\n"
	content += "3日後 | 4日後 | 5日後\n"
	content += strconv.Itoa(sf.Snows[0].Morning) + addRainyChar(sf.Rains[0].Morning) + ", " + strconv.Itoa(sf.Snows[0].Noon) + addRainyChar(sf.Rains[0].Noon) + ", " + strconv.Itoa(sf.Snows[0].Night) + addRainyChar(sf.Rains[0].Night) + "cm | "
	content += strconv.Itoa(sf.Snows[1].Morning) + addRainyChar(sf.Rains[1].Morning) + ", " + strconv.Itoa(sf.Snows[1].Noon) + addRainyChar(sf.Rains[1].Noon) + ", " + strconv.Itoa(sf.Snows[1].Night) + addRainyChar(sf.Rains[1].Night) + "cm | "
	content += strconv.Itoa(sf.Snows[2].Morning) + addRainyChar(sf.Rains[2].Morning) + ", " + strconv.Itoa(sf.Snows[2].Noon) + addRainyChar(sf.Rains[2].Noon) + ", " + strconv.Itoa(sf.Snows[2].Night) + addRainyChar(sf.Rains[2].Night) + "cm\n"
	content += strconv.Itoa(sf.Snows[3].Morning) + addRainyChar(sf.Rains[3].Morning) + ", " + strconv.Itoa(sf.Snows[3].Noon) + addRainyChar(sf.Rains[3].Noon) + ", " + strconv.Itoa(sf.Snows[3].Night) + addRainyChar(sf.Rains[3].Night) + "cm |"
	content += strconv.Itoa(sf.Snows[4].Morning) + addRainyChar(sf.Rains[4].Morning) + ", " + strconv.Itoa(sf.Snows[4].Noon) + addRainyChar(sf.Rains[4].Noon) + ", " + strconv.Itoa(sf.Snows[4].Night) + addRainyChar(sf.Rains[4].Night) + "cm |"
	content += strconv.Itoa(sf.Snows[5].Morning) + addRainyChar(sf.Rains[5].Morning) + ", " + strconv.Itoa(sf.Snows[5].Noon) + addRainyChar(sf.Rains[5].Noon) + ", " + strconv.Itoa(sf.Snows[5].Night) + addRainyChar(sf.Rains[5].Night) + "cm"

	// 140字までに切り詰めて返す
	if len([]rune(content)) > 140 {
		return utf8string.NewString(content).Slice(0, 140), nil
	}
	return content, nil
}

func addRainyChar(rainfall int) string {
	if rainfall > 5 {
		return "☔️"
	} else if rainfall > 0 {
		return "☂️"
	}
	return ""
}

func main() {
	lambda.Start(post)
}
