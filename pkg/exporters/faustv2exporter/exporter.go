package faustv2exporter

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/boxmein/adctf_scoreboard_exporter/pkg/fetchers/faustv2"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

const SCOPE_NAME string = "faustv2_exporter"

type FaustV2Exporter struct {
	fs                 *flag.FlagSet
	baseURL            *string
	currentURL         *string
	scoreboardRoundUrl *string
	teamsURL           *string
	offense            metric.Float64ObservableGauge
	defense            metric.Float64ObservableGauge
	sla                metric.Float64ObservableGauge
	captures           metric.Int64ObservableGauge
	stolen             metric.Int64ObservableGauge
	tick               metric.Int64ObservableGauge
}

func New() FaustV2Exporter {
	f := FaustV2Exporter{
		fs: flag.NewFlagSet("faustv2", flag.ContinueOnError),
	}

	f.baseURL = f.fs.String("base-url", "", "where is the ctf-gameserver hosted? example: http://localhost:5101")
	f.currentURL = f.fs.String("current-url", "", "current.json URL, defaults to baseUrl + /competition/scoreboard-v2/scoreboard_current.json")
	f.scoreboardRoundUrl = f.fs.String("round-url", "", "round URL, falls back to baseUrl + /competition/scoreboard-v2/scoreboard_round_%d.json")
	f.teamsURL = f.fs.String("teams-url", "", "teams.json URL, falls back to baseUrl + /competition/scoreboard-v2/scoreboard_teams.json")

	return f
}

func (f *FaustV2Exporter) Init(args []string) error {
	err := f.fs.Parse(args)

	if err == flag.ErrHelp {
		os.Exit(0)
	}

	if *f.baseURL == "" && *f.scoreboardRoundUrl == "" && *f.currentURL == "" && *f.teamsURL == "" {
		log.Fatalf("set --base-url, or set --round-url and --current-url")
	}

	if *f.baseURL != "" && *f.currentURL == "" {
		*f.currentURL = *f.baseURL + "/competition/scoreboard-v2/scoreboard_current.json"
	}

	if *f.baseURL != "" && *f.scoreboardRoundUrl == "" {
		*f.scoreboardRoundUrl = *f.baseURL + "/competition/scoreboard-v2/scoreboard_round_%d.json"
	}

	if *f.baseURL != "" && *f.teamsURL == "" {
		*f.teamsURL = *f.baseURL + "/competition/scoreboard-v2/scoreboard_teams.json"
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

	captures, err := meter.Int64ObservableGauge("scoreboard_captures", metric.WithDescription("Flags gained. Faceted by service and team."), metric.WithInt64Callback(f.GetCaptureMetrics))

	if err != nil {
		return fmt.Errorf("while setting up captures gauge: %w", err)
	}

	f.captures = captures

	stolen, err := meter.Int64ObservableGauge("scoreboard_stolen", metric.WithDescription("Flags lost. Faceted by service and team."), metric.WithInt64Callback(f.GetCaptureMetrics))

	if err != nil {
		return fmt.Errorf("while setting up stolen gauge: %w", err)
	}

	f.stolen = stolen

	return nil
}

func (f *FaustV2Exporter) GetRound() (*faustv2.ScoreboardRoundJson, error) {
	tick, err := f.GetTick()
	if err != nil {
		return nil, err
	}
	return faustv2.LoadScoreboardRoundJson(*f.scoreboardRoundUrl, tick)
}

func (f *FaustV2Exporter) GetTeams() (faustv2.ScoreboardTeamsJson, error) {
	return faustv2.LoadTeamsJson(*f.teamsURL)
}

func (f *FaustV2Exporter) GetTick() (int64, error) {
	data, err := faustv2.LoadCurrentJson(*f.currentURL)
	if err != nil {
		return -1, err
	}
	return data.ScoreboardTick, nil
}

func (f *FaustV2Exporter) GetOffenseMetrics(ctx context.Context, observer metric.Float64Observer) error {
	data, err := f.GetRound()
	if err != nil {
		return fmt.Errorf("while loading scoreboard: %w", err)
	}

	teams, err := f.GetTeams()
	if err != nil {
		return fmt.Errorf("while loading teams: %w", err)
	}

	for _, team := range data.Scoreboard {
		svc := team.Services
		teamName := teams[team.ID].Name
		for idx, service := range svc {
			observer.Observe(
				service.Offense,
				metric.WithAttributes(
					attribute.String("team", teamName),
					attribute.String("service", data.Services[idx].Name),
				),
			)
		}
	}
	return nil
}

func (f *FaustV2Exporter) GetDefenseMetrics(ctx context.Context, observer metric.Float64Observer) error {
	data, err := f.GetRound()
	if err != nil {
		return fmt.Errorf("while loading scoreboard: %w", err)
	}

	teams, err := f.GetTeams()
	if err != nil {
		return fmt.Errorf("while loading teams: %w", err)
	}

	for _, team := range data.Scoreboard {
		svc := team.Services
		teamName := teams[team.ID].Name
		for idx, service := range svc {
			observer.Observe(
				service.Defense,
				metric.WithAttributes(
					attribute.String("team", teamName),
					attribute.String("service", data.Services[idx].Name),
				),
			)
		}
	}
	return nil
}

func (f *FaustV2Exporter) GetSLAMetrics(ctx context.Context, observer metric.Float64Observer) error {
	data, err := f.GetRound()
	if err != nil {
		return fmt.Errorf("while loading scoreboard: %w", err)
	}

	teams, err := f.GetTeams()
	if err != nil {
		return fmt.Errorf("while loading teams: %w", err)
	}

	for _, team := range data.Scoreboard {
		svc := team.Services
		teamName := teams[team.ID].Name
		for idx, service := range svc {
			observer.Observe(
				service.SLA,
				metric.WithAttributes(
					attribute.String("team", teamName),
					attribute.String("service", data.Services[idx].Name),
				),
			)
		}
	}
	return nil
}

func (f *FaustV2Exporter) GetCaptureMetrics(ctx context.Context, observer metric.Int64Observer) error {
	data, err := f.GetRound()
	if err != nil {
		return fmt.Errorf("while loading scoreboard: %w", err)
	}

	teams, err := f.GetTeams()
	if err != nil {
		return fmt.Errorf("while loading teams: %w", err)
	}

	for _, team := range data.Scoreboard {
		svc := team.Services
		teamName := teams[team.ID].Name
		for idx, service := range svc {
			observer.Observe(
				service.Captures,
				metric.WithAttributes(
					attribute.String("team", teamName),
					attribute.String("service", data.Services[idx].Name),
				),
			)
		}
	}
	return nil
}

func (f *FaustV2Exporter) GetStolenMetrics(ctx context.Context, observer metric.Int64Observer) error {
	data, err := f.GetRound()
	if err != nil {
		return fmt.Errorf("while loading scoreboard: %w", err)
	}

	teams, err := f.GetTeams()
	if err != nil {
		return fmt.Errorf("while loading teams: %w", err)
	}
	for _, team := range data.Scoreboard {
		svc := team.Services
		teamName := teams[team.ID].Name
		for idx, service := range svc {
			observer.Observe(
				service.Stolen,
				metric.WithAttributes(
					attribute.String("team", teamName),
					attribute.String("service", data.Services[idx].Name),
				),
			)
		}
	}
	return nil
}

func (f *FaustV2Exporter) GetTickMetrics(ctx context.Context, observer metric.Int64Observer) error {
	tick, err := f.GetTick()
	if err != nil {
		return fmt.Errorf("while getting tick: %w", err)
	}

	observer.Observe(tick)
	return nil
}
