package main

import (
  "encoding/json"
  "fmt"
  "os"

  "github.com/kimjooyoon/gooo-jev/gooo"
)

type discoveryDocument struct {
  Discovery gooo.UsageDiscoveryResponse `json:"discovery"`
}

func main() {
  if len(os.Args) != 2 {
    fmt.Fprintln(os.Stderr, "usage: go run ./cmd/gooo-plan <discovery.json>")
    os.Exit(64)
  }
  data, err := os.ReadFile(os.Args[1])
  if err != nil {
    fmt.Fprintf(os.Stderr, "read discovery: %v\n", err)
    os.Exit(1)
  }
  var input discoveryDocument
  if err := json.Unmarshal(data, &input); err != nil {
    fmt.Fprintf(os.Stderr, "decode discovery: %v\n", err)
    os.Exit(1)
  }
  plan, err := gooo.PlanUsageActions(input.Discovery)
  if err != nil {
    fmt.Fprintf(os.Stderr, "plan usage actions: %v\n", err)
    os.Exit(1)
  }
  if err := plan.Validate(); err != nil {
    fmt.Fprintf(os.Stderr, "invalid usage action plan: %v\n", err)
    os.Exit(1)
  }
  if err := json.NewEncoder(os.Stdout).Encode(plan); err != nil {
    fmt.Fprintf(os.Stderr, "write usage action plan: %v\n", err)
    os.Exit(1)
  }
}