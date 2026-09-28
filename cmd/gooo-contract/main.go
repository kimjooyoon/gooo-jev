package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/kimjooyoon/gooo-jev/decision"
)

type contractReport struct {
	Path                    string
	Status                  string
	MissingStage            string
	ContractDigest          string
	DeclarationSource       string
	DeclarationIRDigest     string
	GenerationDigest        string
	ReverseObservationDigest string
	NonExecuting            bool
	NonAuthorizing          bool
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run ./cmd/gooo-contract <contract.gooo>")
		os.Exit(64)
	}
	path := filepath.Clean(os.Args[1])
	source, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "read contract: %v\n", err)
		os.Exit(1)
	}
	projection := decision.DeriveGoooContractProjection(decision.GoooContractInput{
		SourceText:     string(source),
		NonAuthorizing: true,
	})
	report := contractReport{
		Path:                    path,
		Status:                  projection.Status,
		MissingStage:            projection.MissingStage,
		ContractDigest:          projection.ContractDigest,
		DeclarationSource:       projection.DeclarationSource,
		DeclarationIRDigest:     projection.Evidence.IRDigest,
		GenerationDigest:        projection.Evidence.GenerationDigest,
		ReverseObservationDigest: projection.Evidence.ReverseObservationDigest,
		NonExecuting:            projection.NonExecuting,
		NonAuthorizing:          projection.NonAuthorizing,
	}
	if err := json.NewEncoder(os.Stdout).Encode(report); err != nil {
		fmt.Fprintf(os.Stderr, "write report: %v\n", err)
		os.Exit(1)
	}
	if err := projection.Validate(); err != nil {
		fmt.Fprintf(os.Stderr, "invalid contract projection: %v\n", err)
		os.Exit(1)
	}
	if projection.Status == "UNKNOWN" {
		os.Exit(2)
	}
}
