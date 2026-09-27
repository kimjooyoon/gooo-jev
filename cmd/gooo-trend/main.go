package main

import (
  "encoding/json"
  "fmt"
  "os"
  "path/filepath"
  "time"

  "github.com/kimjooyoon/gooo-jev/decision"
)

type windowDocument struct {
  Window decision.FeedbackWindow `json:"window"`
}

type report struct {
  Previous decision.FeedbackWindow `json:"previous"`
  Current decision.FeedbackWindow `json:"current"`
  Trend decision.FeedbackWindowTrend `json:"trend"`
}

func readWindow(path string) (decision.FeedbackWindow, error) {
  data, err := os.ReadFile(filepath.Clean(path))
  if err != nil {
    return decision.FeedbackWindow{}, err
  }
  var document windowDocument
  if err := json.Unmarshal(data, &document); err != nil {
    return decision.FeedbackWindow{}, err
  }
  if err := document.Window.Validate(); err != nil {
    return decision.FeedbackWindow{}, err
  }
  return document.Window, nil
}

func main() {
  if len(os.Args) != 4 {
    fmt.Fprintln(os.Stderr, "usage: go run ./cmd/gooo-trend <previous-window.json> <current-window.json> <compared-at-rfc3339>")
    os.Exit(64)
  }
  previous, err := readWindow(os.Args[1])
  if err != nil {
    fmt.Fprintf(os.Stderr, "read previous window: %v\n", err)
    os.Exit(1)
  }
  current, err := readWindow(os.Args[2])
  if err != nil {
    fmt.Fprintf(os.Stderr, "read current window: %v\n", err)
    os.Exit(1)
  }
  comparedAt, err := time.Parse(time.RFC3339Nano, os.Args[3])
  if err != nil {
    fmt.Fprintf(os.Stderr, "parse comparison time: %v\n", err)
    os.Exit(1)
  }
  trend, err := decision.CompareFeedbackWindows(previous, current, comparedAt)
  if err != nil {
    fmt.Fprintf(os.Stderr, "compare feedback windows: %v\n", err)
    os.Exit(1)
  }
  if err := trend.Validate(); err != nil {
    fmt.Fprintf(os.Stderr, "invalid feedback window trend: %v\n", err)
    os.Exit(1)
  }
  if err := json.NewEncoder(os.Stdout).Encode(report{Previous: previous, Current: current, Trend: trend}); err != nil {
    fmt.Fprintf(os.Stderr, "write feedback trend: %v\n", err)
    os.Exit(1)
  }
}