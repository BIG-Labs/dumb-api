package actions

import (
	"errors"
	"net/http"
	"time"

	"dumb-api/models"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop/v6"
)

type Candlestick struct {
	Open  float64
	Close float64
	High  float64
	Low   float64
	Time  time.Time
}

type Candlesticks struct {
	OneMinute   []Candlestick
	FiveMinutes []Candlestick
	OneHour     []Candlestick
	OneDay      []Candlestick
}

type Group struct {
	Time  time.Time
	Ticks []models.PriceTick
}

func GetPriceData(c buffalo.Context) error {

	tokenIn := c.Param("tokenIn")
	tokenOut := c.Param("tokenOut")

	var candlesticks Candlesticks

	tx := c.Value("tx").(*pop.Connection)
	var results []models.PriceTick
	query := tx.Where("token_in = ? AND token_out = ?", tokenIn, tokenOut).Order("created_at asc")
	err := query.All(&results)
	if err != nil {
		return c.Error(http.StatusInternalServerError, err)
	}
	candlesticks, err = generateAllCandlesticks(results)
	if err != nil {
		return c.Error(http.StatusInternalServerError, err)
	}

	return c.Render(http.StatusOK, r.JSON(candlesticks))

}

func generateAllCandlesticks(results []models.PriceTick) (Candlesticks, error) {
	var candlesticks Candlesticks

	ticks := []string{"1m", "5m", "1h", "1d"}
	for _, tick := range ticks {
		cs, err := generateCandlesticks(results, tick)
		if err != nil {
			return Candlesticks{}, err
		}
		switch tick {
		case "1m":
			candlesticks.OneMinute = cs
		case "5m":
			candlesticks.FiveMinutes = cs
		case "1h":
			candlesticks.OneHour = cs
		case "1d":
			candlesticks.OneDay = cs
		}
	}

	return candlesticks, nil
}

func generateCandlesticks(results []models.PriceTick, tick string) ([]Candlestick, error) {
	var candlesticks []Candlestick

	tickDurations := map[string]time.Duration{
		"1m": 1 * time.Minute,
		"5m": 5 * time.Minute,
		"1h": 1 * time.Hour,
		"1d": 24 * time.Hour,
	}

	tickDuration, ok := tickDurations[tick]
	if !ok {
		return nil, errors.New("invalid tick")
	}

	var groups []Group
	var currentGroup Group

	for _, result := range results {
		groupTime := result.CreatedAt.Truncate(tickDuration)
		if currentGroup.Time != groupTime {
			if currentGroup.Time != (time.Time{}) {
				groups = append(groups, currentGroup)
			}
			currentGroup = Group{
				Time:  groupTime,
				Ticks: []models.PriceTick{result},
			}
		} else {
			currentGroup.Ticks = append(currentGroup.Ticks, result)
		}
	}

	if len(currentGroup.Ticks) > 0 {
		groups = append(groups, currentGroup)
	}

	start := len(groups) - 365
	if start < 0 {
		start = 0
	}

	for j, group := range groups[start:] {
		var open, close, high, low float64
		for i, result := range group.Ticks {
			price := result.Price
			if i == 0 {
				if j == 0 {
					open = price
				} else {
					open = candlesticks[len(candlesticks)-1].Close
				}
			}
			if i == len(group.Ticks)-1 {
				close = price
			}
			if price > high || i == 0 {
				high = price
			}
			if price < low || i == 0 {
				low = price
			}
		}
		candlesticks = append(candlesticks, Candlestick{Open: open, Close: close, High: high, Low: low, Time: group.Time})
	}

	if candlesticks == nil {
		candlesticks = []Candlestick{}
	}

	return candlesticks, nil
}
