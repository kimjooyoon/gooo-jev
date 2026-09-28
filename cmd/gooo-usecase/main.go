package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/kimjooyoon/gooo-jev/decision"
)

type report struct {
	Path                     string
	Status                   string
	MissingStage             string
	SourceDigest             string
	IRDigest                 string
	GenerationDigest         string
	RoundTripIRDigest        string
	ReverseObservationDigest string
	GeneratedSource          string
	NonExecuting             bool
	NonAuthorizing           bool
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run ./cmd/gooo-usecase <declaration.gooo>")
		os.Exit(64)
	}

	path := filepath.Clean(os.Args[1])
	source, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "read declaration: %v\n", err)
		os.Exit(1)
	}

	result := decision.DeriveGoooDeclarationIRGeneration(
		decision.GoooDeclarationIRGenerationInput{
			SourceText:     string(source),
			NonAuthorizing: true,
		},
	)
	output := report{
		Path:                     path,
		Status:                   result.Status,
		MissingStage:             result.MissingStage,
		SourceDigest:             result.SourceDigest,
		IRDigest:                 result.IRDigest,
		GenerationDigest:         result.GenerationDigest,
		RoundTripIRDigest:        result.RoundTripIRDigest,
		ReverseObservationDigest: result.ReverseObservationDigest,
		GeneratedSource:          result.GeneratedSource,
		NonExecuting:             result.NonExecuting,
		NonAuthorizing:           result.NonAuthorizing,
	}
	if err := json.NewEncoder(os.Stdout).Encode(output); err != nil {
		fmt.Fprintf(os.Stderr, "write report: %v\n", err)
		os.Exit(1)
	}
	if err := result.Validate(); err != nil {
		fmt.Fprintf(os.Stderr, "invalid declaration evidence: %v\n", err)
		os.Exit(1)
	}
	if result.Status == "UNKNOWN" {
		os.Exit(2)
	}
}
