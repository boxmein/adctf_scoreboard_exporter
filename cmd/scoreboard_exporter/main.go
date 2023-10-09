package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/boxmein/adctf_scoreboard_exporter/pkg/exporters/faustv1exporter"
	"github.com/boxmein/adctf_scoreboard_exporter/pkg/exporters/faustv2exporter"
	"github.com/boxmein/adctf_scoreboard_exporter/pkg/metrics"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	listenAddr = flag.String("listenAddr", ":5001", "address to listen on (e.g. localhost:5001)")
)

func main() {
	flag.Parse()
	rest := flag.Args()

	if cleanup, err := metrics.Setup(); err != nil {
		log.Fatalf("error setting up metrics: %v", err)
	} else {
		defer cleanup()
	}

	subcmd := rest[0]
	args := rest[1:]

	if subcmd == "faustv1" {
		exporter := faustv1exporter.New()
		if err := exporter.Init(args); err != nil {
			log.Fatalf("error: %v", err)
		}
	} else if subcmd == "faustv2" {
		exporter := faustv2exporter.New()
		if err := exporter.Init(args); err != nil {
			log.Fatalf("error: %v", err)
		}
	} else {
		log.Fatalf("subcmd is not accepted! only faustv1 and faustv2 are accepted")
	}

	http.Handle("/metrics", promhttp.Handler())

	log.Printf("Listening on http://%s/metrics", *listenAddr)
	if err := http.ListenAndServe(*listenAddr, nil); err != nil {
		log.Fatalf("error while listening: %v", err)
	}
}
