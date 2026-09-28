package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/kimjooyoon/gooo-jev/decision"
)

func main() {
	policyDigest := flag.String("policy-digest", "", "digest identifying the deterministic policy")
	minimumConfidence := flag.Float64("min-confidence", -1, "minimum confidence required for an eligible disposition")
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "usage: go run ./cmd/gooo-decision-policy -policy-digest <digest> -min-confidence <0..1> [observation.json|-]")
	}
	flag.Parse()
	if flag.NArg() > 1 || *minimumConfidence < 0 || *minimumConfidence > 1 {
		flag.Usage()
		os.Exit(64)
	}

	var data []byte
	var err error
	if flag.NArg() == 0 || flag.Arg(0) == "-" {
		data, err = io.ReadAll(os.Stdin)
	} else {
		data, err = os.ReadFile(flag.Arg(0))
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "read decision observation: %v\n", err)
		os.Exit(1)
	}

	policy, err := decision.NewDecisionPolicy(*policyDigest, *minimumConfidence)
	if err != nil {
		fmt.Fprintf(os.Stderr, "create decision policy: %v\n", err)
		os.Exit(1)
	}
	outcome, err := decision.EvaluateJSONPolicy(data, policy)
	if err != nil {
		fmt.Fprintf(os.Stderr, "evaluate decision policy: %v\n", err)
		os.Exit(1)
	}
	if err := outcome.Validate(); err != nil {
		fmt.Fprintf(os.Stderr, "invalid policy outcome: %v\n", err)
		os.Exit(1)
	}
	if err := json.NewEncoder(os.Stdout).Encode(outcome); err != nil {
		fmt.Fprintf(os.Stderr, "write policy outcome: %v\n", err)
		os.Exit(1)
	}
}
