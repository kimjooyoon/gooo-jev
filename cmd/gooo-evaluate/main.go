package main

import (
  "encoding/json"
  "fmt"
  "os"
  "path/filepath"

  "github.com/kimjooyoon/gooo-jev/decision"
)

type feedbackDocument struct {
  Window decision.FeedbackWindow `json:"window"`
}

type report struct {
  Declaration any `json:"declaration"`
  FeedbackWindow decision.FeedbackWindow `json:"feedback_window"`
  ReadyForObservation bool `json:"ready_for_observation"`
  NonExecuting bool `json:"non_executing"`
  NonAuthorizing bool `json:"non_authorizing"`
}

func main() {
  if len(os.Args) != 3 {
    fmt.Fprintln(os.Stderr, "usage: go run ./cmd/gooo-evaluate <declaration.gooo> <feedback-window.json>")
    os.Exit(64)
  }

  declarationPath := filepath.Clean(os.Args[1])
  feedbackPath := filepath.Clean(os.Args[2])
  source, err := os.ReadFile(declarationPath)
  if err != nil {
    fmt.Fprintf(os.Stderr, "read declaration: %v\n", err)
    os.Exit(1)
  }
  feedbackBytes, err := os.ReadFile(feedbackPath)
  if err != nil {
    fmt.Fprintf(os.Stderr, "read feedback window: %v\n", err)
    os.Exit(1)
  }
  var feedback feedbackDocument
  if err := json.Unmarshal(feedbackBytes, &feedback); err != nil {
    fmt.Fprintf(os.Stderr, "decode feedback window: %v\n", err)
    os.Exit(1)
  }
  if err := feedback.Window.Validate(); err != nil {
    fmt.Fprintf(os.Stderr, "invalid feedback window: %v\n", err)
    os.Exit(1)
  }

  declaration := decision.DeriveGoooDeclarationIRGeneration(
    decision.GoooDeclarationIRGenerationInput{
      SourceText: string(source),
      NonAuthorizing: true,
    },
  )
  if err := declaration.Validate(); err != nil {
    fmt.Fprintf(os.Stderr, "invalid declaration evidence: %v\n", err)
    os.Exit(1)
  }
  ready := declaration.Status == "BOUND" && declaration.NonExecuting && declaration.NonAuthorizing && feedback.Window.NonAuthorizing
  output := report{
    Declaration: declaration,
    FeedbackWindow: feedback.Window,
    ReadyForObservation: ready,
    NonExecuting: declaration.NonExecuting,
    NonAuthorizing: declaration.NonAuthorizing && feedback.Window.NonAuthorizing,
  }
  if err := json.NewEncoder(os.Stdout).Encode(output); err != nil {
    fmt.Fprintf(os.Stderr, "write evaluation: %v\n", err)
    os.Exit(1)
  }
  if !ready {
    fmt.Fprintf(os.Stderr, "observation not ready: declaration_status=%s missing_stage=%s coverage=%g\n", declaration.Status, declaration.MissingStage, feedback.Window.ObservedCoverage)
    os.Exit(2)
  }
}