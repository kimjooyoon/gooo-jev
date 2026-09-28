package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/kimjooyoon/gooo-jev/gooo"
)

type inputDocument struct {
	Query          string   `json:"query"`
	Declaration    string   `json:"declaration,omitempty"`
	VerifiedStages []string `json:"verified_stages"`
}

type outputDocument struct {
	Response gooo.CapabilityQueryResponse `json:"response"`
	Feedback gooo.CapabilityQueryFeedback `json:"feedback"`
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run ./cmd/gooo-capability-feedback <feedback.json>")
		os.Exit(64)
	}
	path := filepath.Clean(os.Args[1])
	file, err := os.Open(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "open capability feedback: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	var input inputDocument
	if err := json.NewDecoder(file).Decode(&input); err != nil {
		fmt.Fprintf(os.Stderr, "decode capability feedback: %v\n", err)
		os.Exit(1)
	}

	var response gooo.CapabilityQueryResponse
	if strings.TrimSpace(input.Declaration) == "" {
		response = gooo.DiscoverCapabilityQuery(input.Query)
	} else {
		response = gooo.DiscoverCapabilityQueryWithDeclaration(input.Query, input.Declaration)
	}
	feedback, err := gooo.ObserveCapabilityQueryFeedback(response, input.VerifiedStages)
	if err != nil {
		fmt.Fprintf(os.Stderr, "observe capability feedback: %v\n", err)
		os.Exit(1)
	}
	if err := feedback.Validate(); err != nil {
		fmt.Fprintf(os.Stderr, "invalid capability feedback: %v\n", err)
		os.Exit(1)
	}
	if err := json.NewEncoder(os.Stdout).Encode(outputDocument{Response: response, Feedback: feedback}); err != nil {
		fmt.Fprintf(os.Stderr, "write capability feedback: %v\n", err)
		os.Exit(1)
	}
}
