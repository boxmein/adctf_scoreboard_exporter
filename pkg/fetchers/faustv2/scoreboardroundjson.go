package faustv2

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
)

func LoadScoreboardRoundJson(urlPattern string, roundId int64) (*ScoreboardRoundJson, error) {
	url := fmt.Sprintf(urlPattern, roundId)
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
	unpacked := new(ScoreboardRoundJson)
	err = json.Unmarshal(data, unpacked)
	if err != nil {
		return nil, err
	}

	return unpacked, nil
}

type ScoreboardRoundJson struct {
	Tick               int64                 `json:"tick"`
	Scoreboard         []ScoreboardV2Team    `json:"scoreboard"`
	StatusDescriptions map[int64]string      `json:"status-descriptions"`
	Services           []ScoreboardV2Service `json:"services"`
}

type ScoreboardV2Team struct {
	Rank         int64                       `json:"rank"`
	ID           int64                       `json:"team_id"`
	Services     []*ScoreboardV2ServiceScore `json:"services"`
	Points       float64                     `json:"points"`
	Offense      float64                     `json:"o"`
	OffenseDelta float64                     `json:"do"`
	Defense      float64                     `json:"d"`
	DefenseDelta float64                     `json:"dd"`
	SLA          float64                     `json:"s"`
	SLADelta     float64                     `json:"ds"`
}

type ScoreboardV2ServiceScore struct {
	Message       string  `json:"m"`
	Status        int64   `json:"c"`
	StatusDelta   []int64 `json:"dc"`
	Offense       float64 `json:"o"`
	OffenseDelta  float64 `json:"do"`
	Defense       float64 `json:"d"`
	DefenseDelta  float64 `json:"dd"`
	SLA           float64 `json:"s"`
	SLADelta      float64 `json:"ds"`
	Captures      int64   `json:"cap"`
	CapturesDelta int64   `json:"dcap"`
	Stolen        int64   `json:"st"`
	StolenDelta   int64   `json:"dst"`
}

type ScoreboardV2Service struct {
	Name       string  `json:"name"`
	Attackers  int64   `json:"attackers"`
	Victims    int64   `json:"victims"`
	FirstBlood []int64 `json:"first_blood"`
}
