package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/kimjooyoon/gooo-jev/gooo"
)

type report struct {
	Path                string                        `json:"path,omitempty"`
	Discovery            *gooo.UsageDiscoveryResponse  `json:"discovery,omitempty"`
	CapabilityDiscovery  *gooo.CapabilityQueryResponse `json:"capability_discovery,omitempty"`
	CapabilityTrail      *gooo.CapabilityQueryTrail    `json:"capability_trail,omitempty"`
}

func main() {
	query := flag.String("query", "", "discover capabilities from a natural-language query")
	flag.Parse()
	if flag.NArg() > 2 {
		fmt.Fprintln(os.Stderr, "usage: go run ./cmd/gooo-discover [--query <question>] [declaration.gooo] [prefix]")
		os.Exit(64)
	}
	var path string
	var source []byte
	if flag.NArg() > 0 {
		path = filepath.Clean(flag.Arg(0))
		var err error
		source, err = os.ReadFile(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "read declaration: %v\n", err)
			os.Exit(1)
		}
	}
	queryText := strings.TrimSpace(*query)
	if queryText == "" && len(source) > 0 {
		queryText = "What can gooo do with this declaration?"
	}
	var capabilityDiscovery *gooo.CapabilityQueryResponse
	var capabilityTrail *gooo.CapabilityQueryTrail
	if queryText != "" {
		trail := gooo.DiscoverCapabilityQueryTrail(queryText, string(source))
		if err := trail.Validate(); err != nil {
			fmt.Fprintf(os.Stderr, "invalid capability query trail: %v\n", err)
			os.Exit(1)
		}
		capabilityDiscovery = &trail.Response
		capabilityTrail = &trail
	}
	if flag.NArg() == 0 {
		if capabilityDiscovery == nil {
			fmt.Fprintln(os.Stderr, "usage: go run ./cmd/gooo-discover [--query <question>] [declaration.gooo] [prefix]")
			os.Exit(64)
		}
		if err := json.NewEncoder(os.Stdout).Encode(report{CapabilityDiscovery: capabilityDiscovery, CapabilityTrail: capabilityTrail}); err != nil {
			fmt.Fprintf(os.Stderr, "write capability discovery: %v\n", err)
			os.Exit(1)
		}
		return
	}
	prefix := ""
	if flag.NArg() == 2 {
		prefix = flag.Arg(1)
	}
	discovery := gooo.DiscoverUsage(string(source), prefix)
	if err := discovery.Validate(); err != nil {
		fmt.Fprintf(os.Stderr, "invalid usage discovery: %v\n", err)
		os.Exit(1)
	}
	if err := json.NewEncoder(os.Stdout).Encode(report{Path: path, Discovery: &discovery, CapabilityDiscovery: capabilityDiscovery, CapabilityTrail: capabilityTrail}); err != nil {
		fmt.Fprintf(os.Stderr, "write usage discovery: %v\n", err)
		os.Exit(1)
	}
}
