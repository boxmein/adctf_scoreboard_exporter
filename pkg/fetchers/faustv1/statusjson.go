package faustv1

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
)

func LoadStatusJson(url string) (*StatusJson, error) {
	resp, err := http.DefaultClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("while loading status.json: %w", err)
	}
	if resp == nil {
		return nil, errors.New("resp was nil for some reason")
	}
	log.Printf("GET %s => %s", url, resp.Status)
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("while reading body: %w", err)
	}
	unpacked := new(StatusJson)
	err = json.Unmarshal(data, unpacked)
	if err != nil {
		return nil, fmt.Errorf("while unmarshaling into StatusJson: %w", err)
	}

	return unpacked, nil
}

type StatusJson struct {
	// Mapping of array index to tick number
	Ticks []int64          `json:"ticks"`
	Teams []StatusJsonTeam `json:"teams"`
	// Mapping of status code to status meaning
	StatusDescriptions map[int64]string `json:"status-descriptions"`
	// Mapping of service index to service name
	Services []string `json:"services"`
}

type StatusJsonTeam struct {
	ID   int64  `json:"id"`
	Nop  bool   `json:"nop"`
	Name string `json:"name"`
	// mapping of service index to last 5 status codes
	Ticks [][]int64 `json:"ticks"`
}
