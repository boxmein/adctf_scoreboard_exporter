package faustv2

import (
	"encoding/json"
	"io"
	"log"

	"github.com/boxmein/adctf_scoreboard_exporter/pkg/httpclient"
)

func LoadTeamsJson(url string) (ScoreboardTeamsJson, error) {
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
