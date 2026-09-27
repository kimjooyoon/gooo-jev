package main

import (
  "encoding/json"
  "fmt"
  "os"
  "path/filepath"

  "github.com/kimjooyoon/gooo-jev/gooo"
)

type report struct {
  Path string `json:"path"`
  Discovery gooo.UsageDiscoveryResponse `json:"discovery"`
}

func main() {
  if len(os.Args) < 2 || len(os.Args) > 3 {
    fmt.Fprintln(os.Stderr, "usage: go run ./cmd/gooo-discover <declaration.gooo> [prefix]")
    os.Exit(64)
  }
  path := filepath.Clean(os.Args[1])
  prefix := ""
  if len(os.Args) == 3 {
    prefix = os.Args[2]
  }
  source, err := os.ReadFile(path)
  if err != nil {
    fmt.Fprintf(os.Stderr, "read declaration: %v\n", err)
    os.Exit(1)
  }
  discovery := gooo.DiscoverUsage(string(source), prefix)
  if err := discovery.Validate(); err != nil {
    fmt.Fprintf(os.Stderr, "invalid usage discovery: %v\n", err)
    os.Exit(1)
  }
  if err := json.NewEncoder(os.Stdout).Encode(report{Path: path, Discovery: discovery}); err != nil {
    fmt.Fprintf(os.Stderr, "write usage discovery: %v\n", err)
    os.Exit(1)
  }
}