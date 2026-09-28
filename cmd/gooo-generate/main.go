package main

import (
  "fmt"
  "os"
  "path/filepath"

  "github.com/kimjooyoon/gooo-jev/decision"
)

func main() {
  if len(os.Args) != 3 {
    fmt.Fprintln(os.Stderr, "usage: go run ./cmd/gooo-generate <input.gooo> <output.gooo>")
    os.Exit(64)
  }

  inputPath := filepath.Clean(os.Args[1])
  outputPath := filepath.Clean(os.Args[2])
  source, err := os.ReadFile(inputPath)
  if err != nil {
    fmt.Fprintf(os.Stderr, "read declaration: %v\n", err)
    os.Exit(1)
  }

  result := decision.DeriveGoooDeclarationIRGeneration(
    decision.GoooDeclarationIRGenerationInput{
      SourceText: string(source),
      NonAuthorizing: true,
    },
  )
  if err := result.Validate(); err != nil {
    fmt.Fprintf(os.Stderr, "invalid declaration evidence: %v\n", err)
    os.Exit(1)
  }
  if result.Status == "UNKNOWN" {
    fmt.Fprintf(os.Stderr, "generation unavailable: missing_stage=%s\n", result.MissingStage)
    os.Exit(2)
  }
  if !result.NonExecuting || !result.NonAuthorizing {
    fmt.Fprintln(os.Stderr, "generation crossed a capability boundary")
    os.Exit(1)
  }
  if err := os.WriteFile(outputPath, []byte(result.GeneratedSource), 0o644); err != nil {
    fmt.Fprintf(os.Stderr, "write generated declaration: %v\n", err)
    os.Exit(1)
  }
  fmt.Fprintf(os.Stderr, "source=%s ir=%s generation=%s round_trip=%s reverse=%s\n", result.SourceDigest, result.IRDigest, result.GenerationDigest, result.RoundTripIRDigest, result.ReverseObservationDigest)
}