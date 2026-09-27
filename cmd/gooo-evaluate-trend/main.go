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
  Declaration any `json:"declaration"`
  PreviousWindow decision.FeedbackWindow `json:"previous_window"`
  CurrentWindow decision.FeedbackWindow `json:"current_window"`
  Trend decision.FeedbackWindowTrend `json:"trend"`
  ReadyForObservation bool `json:"ready_for_observation"`
  NonExecuting bool `json:"non_executing"`
  NonAuthorizing bool `json:"non_authorizing"`
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
  if len(os.Args) != 5 {
    fmt.Fprintln(os.Stderr, "usage: go run ./cmd/gooo-evaluate-trend <declaration.gooo> <previous-window.json> <current-window.json> <compared-at-rfc3339>")
    os.Exit(64)
  }
  source, err := os.ReadFile(filepath.Clean(os.Args[1]))
  if err != nil {
    fmt.Fprintf(os.Stderr, "read declaration: %v\n", err)
    os.Exit(1)
  }
  previous, err := readWindow(os.Args[2])
  if err != nil {
    fmt.Fprintf(os.Stderr, "read previous window: %v\n", err)
    os.Exit(1)
  }
  current, err := readWindow(os.Args[3])
  if err != nil {
    fmt.Fprintf(os.Stderr, "read current window: %v\n", err)
    os.Exit(1)
  }
  comparedAt, err := time.Parse(time.RFC3339Nano, os.Args[4])
  if err != nil {
    fmt.Fprintf(os.Stderr, "parse comparison time: %v\n", err)
    os.Exit(1)
  }
  declaration := decision.DeriveGoooDeclarationIRGeneration(decision.GoooDeclarationIRGenerationInput{
    SourceText: string(source),
    NonAuthorizing: true,
  })
  if err := declaration.Validate(); err != nil {
    fmt.Fprintf(os.Stderr, "invalid declaration evidence: %v\n", err)
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
  ready := declaration.Status == "BOUND" && declaration.NonExecuting && declaration.NonAuthorizing && trend.NonAuthorizing
  output := report{
    Declaration: declaration,
    PreviousWindow: previous,
    CurrentWindow: current,
    Trend: trend,
    ReadyForObservation: ready,
    NonExecuting: declaration.NonExecuting,
    NonAuthorizing: declaration.NonAuthorizing && trend.NonAuthorizing,
  }
  if err := json.NewEncoder(os.Stdout).Encode(output); err != nil {
    fmt.Fprintf(os.Stderr, "write evaluation trend: %v\n", err)
    os.Exit(1)
  }
  if !ready {
    fmt.Fprintf(os.Stderr, "observation not ready: declaration_status=%s missing_stage=%s\n", declaration.Status, declaration.MissingStage)
    os.Exit(2)
  }
}