package snow_forecaster

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

type SnowForecaster struct{}

func NewSnowForecaster() *SnowForecaster {
	return &SnowForecaster{}
}

// 以下の3パターンの予報が取得できる
// 1.本日の朝からの予報が見れる時
// 2.本日の昼からの予報が見れる時
// 3.本日の夜からの予報が見れる時
func (sf *SnowForecaster) Forecast(key string) (Forecast, error) {
	doc, err := goquery.NewDocument(fmt.Sprintf("https://ja.snow-forecast.com/resorts/%s/6day/top", key))
	if err != nil {
		return Forecast{}, err
	}

	var convertErrFlag bool
	snows := make([]Snow, 0)
	forecastTableSnow := doc.Find("td.forecast-table-snow__cell")
	forecastTableSnow.Each(func(index int, s *goquery.Selection) {
		if s.HasClass("day-end") {
			morning, noon := -1, -1
			// 夜の情報は必ず取得できる
			night := convertInt(s.Text(), &convertErrFlag)
			// 昼の情報も取得できる時
			if index > 0 {
				noon = convertInt(forecastTableSnow.Eq(index-1).Text(), &convertErrFlag)
			}
			// 朝の情報も取得できる時
			if index > 1 {
				morning = convertInt(forecastTableSnow.Eq(index-2).Text(), &convertErrFlag)
			}
			snows = append(snows, NewSnow(morning, noon, night))
		}
	})
	if convertErrFlag {
		return Forecast{}, errors.New("convert error occur")
	}

	rains := make([]Rain, 0)
	forecastTableRain := doc.Find("td.forecast-table-rain__cell")
	forecastTableRain.Each(func(index int, s *goquery.Selection) {
		if s.HasClass("day-end") {
			morning, noon := -1, -1
			// 夜の情報は必ず取得できる
			night := convertInt(s.Text(), &convertErrFlag)
			// 昼の情報も取得できる時
			if index > 0 {
				noon = convertInt(forecastTableRain.Eq(index-1).Text(), &convertErrFlag)
			}
			// 朝の情報も取得できる時
			if index > 1 {
				morning = convertInt(forecastTableRain.Eq(index-2).Text(), &convertErrFlag)
			}
			rains = append(rains, NewRain(morning, noon, night))
		}
	})
	if convertErrFlag {
		return Forecast{}, errors.New("convert error occur")
	}

	if len(snows) == 0 || len(rains) == 0 {
		return Forecast{}, errors.New("can not get forecasts")
	}
	return NewForecast(snows, rains, key), nil
}

// "-"の場合は特別に0を返すキャスト
func convertInt(s string, convertErrFlag *bool) int {
	if s == "-" {
		return 0
	}
	ret, err := strconv.Atoi(s)
	if err != nil {
		*convertErrFlag = true
	}
	return ret
}

type Forecast struct {
	Snows     []Snow
	Rains     []Rain
	SkiResort string
}

func NewForecast(snows []Snow, rains []Rain, skiResort string) Forecast {
	return Forecast{
		Snows:     snows,
		Rains:     rains,
		SkiResort: skiResort,
	}
}

type Snow struct {
	Morning int
	Noon    int
	Night   int
}

func NewSnow(morning, noon, night int) Snow {
	return Snow{
		Morning: morning,
		Noon:    noon,
		Night:   night,
	}
}

type Rain struct {
	Morning int
	Noon    int
	Night   int
}

func NewRain(morning, noon, night int) Rain {
	return Rain{
		Morning: morning,
		Noon:    noon,
		Night:   night,
	}
}
