package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/kimjooyoon/gooo-jev/gooo"
)

type planDocument struct {
	Plan gooo.UsageActionPlan `json:"plan"`
}

type observationDocument struct {
	Observation gooo.UsageObservation `json:"observation"`
}

func main() {
	if len(os.Args) != 8 {
		fmt.Fprintln(os.Stderr, "usage: gooo-observe PLAN ACTION_ID KIND METRIC_NAME METRIC_VALUE EVIDENCE_DIGEST RECORDED_AT_RFC3339")
		os.Exit(2)
	}

	data, err := os.ReadFile(os.Args[1])
	if err != nil {
		fail(err)
	}
	var document planDocument
	if err := json.Unmarshal(data, &document); err != nil {
		fail(err)
	}
	metricValue, err := strconv.ParseFloat(os.Args[5], 64)
	if err != nil {
		fail(fmt.Errorf("metric value: %w", err))
	}
	recordedAt, err := time.Parse(time.RFC3339, os.Args[7])
	if err != nil {
		fail(fmt.Errorf("recorded at: %w", err))
	}

	observation, err := gooo.ObserveUsageAction(
		document.Plan,
		os.Args[2],
		gooo.UsageObservationKind(os.Args[3]),
		os.Args[4],
		metricValue,
		os.Args[6],
		recordedAt,
	)
	if err != nil {
		fail(err)
	}
	if err := observation.ValidateAgainst(document.Plan); err != nil {
		fail(err)
	}
	output := observationDocument{Observation: observation}
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(output); err != nil {
		fail(err)
	}
}

func fail(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(2)
}

