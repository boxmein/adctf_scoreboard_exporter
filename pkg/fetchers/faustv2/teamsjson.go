package faustv2

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
)

func LoadTeamsJson(url string) (ScoreboardTeamsJson, error) {
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
	unpacked := new(ScoreboardTeamsJson)
	err = json.Unmarshal(data, unpacked)
	if err != nil {
		return nil, err
	}

	return *unpacked, nil
}

type ScoreboardTeamsJson map[int64]TeamsJsonTeam

type TeamsJsonTeam struct {
	Name        string `json:"name"`
	Affiliation string `json:"aff"`
	Vulnbox     string `json:"vulnbox"`
	Logo        string `json:"logo"`
}
