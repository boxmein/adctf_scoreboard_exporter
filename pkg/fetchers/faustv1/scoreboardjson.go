package faustv1

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
)

func LoadScoreboardJson(url string) (*ScoreboardJson, error) {
	resp, err := http.DefaultClient.Get(url)
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, errors.New("resp was nil for some reason")
	}
	log.Printf("GET %s => %s", url, resp.Status)
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	unpacked := new(ScoreboardJson)
	err = json.Unmarshal(data, unpacked)
	if err != nil {
		return nil, err
	}

	return unpacked, nil
}

type ScoreboardJson struct {
	Tick               int64                `json:"tick"`
	Teams              []ScoreboardJsonTeam `json:"teams"`
	StatusDescriptions map[int64]string     `json:"status-descriptions"`
}

type ScoreboardJsonTeam struct {
	Rank      int64      `json:"rank"`
	ID        int64      `json:"id"`
	Name      string     `json:"name"`
	Services  []*Service `json:"services"`
	Offense   float64    `json:"offense"`
	Defense   float64    `json:"defense"`
	SLA       float64    `json:"sla"`
	Total     float64    `json:"total"`
	Image     *string    `json:"image"`
	Thumbnail *string    `json:"thumbnail"`
}

type Service struct {
	Status  int64   `json:"status"`
	Offense float64 `json:"offense"`
	Defense float64 `json:"defense"`
	SLA     float64 `json:"sla"`
}
