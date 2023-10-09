package faustv2

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
)

func LoadCurrentJson(url string) (*CurrentJson, error) {
	resp, err := http.DefaultClient.Get(url)
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, errors.New("response was nil for some reason")
	}
	log.Printf("GET %s => %s", url, resp.Status)
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	unpacked := new(CurrentJson)
	err = json.Unmarshal(data, unpacked)
	if err != nil {
		return nil, err
	}

	return unpacked, nil
}

type CurrentJson struct {
	State            int64   `json:"state"`
	CurrentTick      int64   `json:"current_tick"`
	CurrentTickUntil float64 `json:"current_tick_until"`
	ScoreboardTick   int64   `json:"scoreboard_tick"`
}
