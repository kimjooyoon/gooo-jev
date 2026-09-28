package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/kimjooyoon/gooo-jev/decision"
)

type inputDocument struct {
	MetricName  string
	Observations []decision.FeedbackObservation
}

type outputDocument struct {
	MetricName string
	Window     decision.FeedbackWindow
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run ./cmd/gooo-feedback <feedback.json>")
		os.Exit(64)
	}
	path := filepath.Clean(os.Args[1])
	file, err := os.Open(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "open feedback: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	var input inputDocument
	if err := json.NewDecoder(file).Decode(&input); err != nil {
		fmt.Fprintf(os.Stderr, "decode feedback: %v\n", err)
		os.Exit(1)
	}
	window, err := decision.MeasureFeedbackWindow(input.Observations, input.MetricName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "measure feedback: %v\n", err)
		os.Exit(1)
	}
	if err := window.Validate(); err != nil {
		fmt.Fprintf(os.Stderr, "invalid feedback window: %v\n", err)
		os.Exit(1)
	}
	if err := json.NewEncoder(os.Stdout).Encode(outputDocument{
		MetricName: input.MetricName,
		Window:     window,
	}); err != nil {
		fmt.Fprintf(os.Stderr, "write feedback window: %v\n", err)
		os.Exit(1)
	}
}
