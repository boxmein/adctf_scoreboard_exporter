package faustv1exporter

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/boxmein/adctf_scoreboard_exporter/pkg/fetchers/faustv1"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

const SCOPE_NAME string = "faustv1_exporter"

type FaustV1Exporter struct {
	fs            *flag.FlagSet
	baseURL       *string
	scoreboardURL *string
	statusURL     *string
	offense       metric.Float64ObservableGauge
	defense       metric.Float64ObservableGauge
	sla           metric.Float64ObservableGauge
	tick          metric.Int64ObservableGauge
}

func New() FaustV1Exporter {
	f := FaustV1Exporter{
		fs: flag.NewFlagSet("faustv1", flag.ContinueOnError),
	}

	f.baseURL = f.fs.String("base-url", "", "where is the ctf-gameserver hosted? example: http://localhost:5101")
	f.scoreboardURL = f.fs.String("scoreboard-url", "", "scoreboard.json URL, falls back to baseUrl + /competition/scoreboard.json")
	f.statusURL = f.fs.String("status-url", "", "status.json URL, falls back to baseUrl + /competition/status.json")

	return f
}

func (f *FaustV1Exporter) Init(args []string) error {
	err := f.fs.Parse(args)

	if err == flag.ErrHelp {
		os.Exit(0)
	}

	if *f.baseURL == "" && *f.scoreboardURL == "" && *f.statusURL == "" {
		log.Fatalf("set baseUrl, or set scoreboardUrl and statusUrl")
	}

	if *f.baseURL != "" && *f.scoreboardURL == "" {
		*f.scoreboardURL = *f.baseURL + "/competition/scoreboard.json"
	}

	if *f.baseURL != "" && *f.statusURL == "" {
		*f.statusURL = *f.baseURL + "/competition/status.json"
	}

	meter := otel.Meter(SCOPE_NAME)

	offense, err := meter.Float64ObservableGauge("scoreboard_offense", metric.WithDescription("Offense points. Faceted by service and team."), metric.WithFloat64Callback(f.GetOffenseMetrics))

	if err != nil {
		return fmt.Errorf("while setting up offense gauge: %w", err)
	}

	f.offense = offense

	defense, err := meter.Float64ObservableGauge("scoreboard_defense", metric.WithDescription("Defense points. Faceted by service and team."), metric.WithFloat64Callback(f.GetDefenseMetrics))

	if err != nil {
		return fmt.Errorf("while setting up defense gauge: %w", err)
	}

	f.defense = defense

	sla, err := meter.Float64ObservableGauge("scoreboard_sla", metric.WithDescription("SLA points. Faceted by service and team."), metric.WithFloat64Callback(f.GetSLAMetrics))

	if err != nil {
		return fmt.Errorf("while setting up sla gauge: %w", err)
	}

	f.sla = sla

	tick, err := meter.Int64ObservableGauge("scoreboard_tick", metric.WithDescription("Current tick."), metric.WithInt64Callback(f.GetTickMetrics))

	if err != nil {
		return fmt.Errorf("while setting up tick gauge: %w", err)
	}

	f.tick = tick

	return nil
}

func (f *FaustV1Exporter) GetTeams() (*faustv1.ScoreboardJson, error) {
	return faustv1.LoadScoreboardJson(*f.scoreboardURL)
}

func (f *FaustV1Exporter) GetServiceNamesByIndex() ([]string, error) {
	data, err := faustv1.LoadStatusJson(*f.statusURL)
	if err != nil {
		return nil, err
	}
	return data.Services, nil
}

func (f *FaustV1Exporter) GetTick() (int64, error) {
	data, err := faustv1.LoadScoreboardJson(*f.scoreboardURL)
	if err != nil {
		return -1, err
	}
	return data.Tick, nil
}

func (f *FaustV1Exporter) GetOffenseMetrics(ctx context.Context, observer metric.Float64Observer) error {
	data, err := f.GetTeams()
	if err != nil {
		return fmt.Errorf("while loading teams: %w", err)
	}
	svcNames, err := f.GetServiceNamesByIndex()
	if err != nil {
		return fmt.Errorf("while loading service names: %w", err)
	}

	for _, team := range data.Teams {
		svc := team.Services
		for idx, service := range svc {
			observer.Observe(
				service.Offense,
				metric.WithAttributes(
					attribute.String("team", team.Name),
					attribute.String("service", svcNames[idx]),
				),
			)
		}
	}
	return nil
}

func (f *FaustV1Exporter) GetDefenseMetrics(ctx context.Context, observer metric.Float64Observer) error {
	data, err := f.GetTeams()
	if err != nil {
		return fmt.Errorf("while loading teams: %w", err)
	}

	svcNames, err := f.GetServiceNamesByIndex()
	if err != nil {
		return fmt.Errorf("while loading service names: %w", err)
	}

	for _, team := range data.Teams {
		svc := team.Services
		for idx, service := range svc {
			observer.Observe(
				service.Defense,
				metric.WithAttributes(
					attribute.String("team", team.Name),
					attribute.String("service", svcNames[idx]),
				),
			)
		}
	}
	return nil
}

func (f *FaustV1Exporter) GetSLAMetrics(ctx context.Context, observer metric.Float64Observer) error {
	data, err := f.GetTeams()
	if err != nil {
		return fmt.Errorf("while loading teams: %w", err)
	}

	svcNames, err := f.GetServiceNamesByIndex()
	if err != nil {
		return fmt.Errorf("while loading service names: %w", err)
	}

	for _, team := range data.Teams {
		svc := team.Services
		for idx, service := range svc {
			observer.Observe(
				service.SLA,
				metric.WithAttributes(
					attribute.String("team", team.Name),
					attribute.String("service", svcNames[idx]),
				),
			)
		}
	}
	return nil
}

func (f *FaustV1Exporter) GetTickMetrics(ctx context.Context, observer metric.Int64Observer) error {
	tick, err := f.GetTick()
	if err != nil {
		return fmt.Errorf("while getting tick: %w", err)
	}

	observer.Observe(tick)
	return nil
}
