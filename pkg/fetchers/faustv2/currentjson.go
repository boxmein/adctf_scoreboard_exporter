package faustv2

import (
	"encoding/json"
	"io"
	"log"

	"github.com/boxmein/adctf_scoreboard_exporter/pkg/httpclient"
)

func LoadCurrentJson(url string) (*CurrentJson, error) {
	resp, err := httpclient.HttpClient.Get(url)
	log.Printf("GET %s => %s", url, resp.Status)
	if err != nil {
		return nil, err
	}
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
