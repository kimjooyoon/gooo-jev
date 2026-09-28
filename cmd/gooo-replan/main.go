package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/kimjooyoon/gooo-jev/gooo"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "usage: gooo-replan [window.json|-]")
		fmt.Fprintln(os.Stderr, "derive a non-executing usage replan from an observed window")
	}
	flag.Parse()

	reader := os.Stdin
	if flag.NArg() > 0 && flag.Arg(0) != "-" {
		file, err := os.Open(flag.Arg(0))
		if err != nil {
			fmt.Fprintf(os.Stderr, "open usage window: %v\n", err)
			os.Exit(1)
		}
		defer file.Close()
		reader = file
	}

	var window gooo.UsageObservationWindow
	if err := json.NewDecoder(reader).Decode(&window); err != nil {
		fmt.Fprintf(os.Stderr, "decode usage window: %v\n", err)
		os.Exit(1)
	}
	proposal := gooo.ProposeUsageReplan(window)
	if err := proposal.Validate(); err != nil {
		fmt.Fprintf(os.Stderr, "validate usage replan: %v\n", err)
		os.Exit(1)
	}
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(proposal); err != nil {
		fmt.Fprintf(os.Stderr, "encode usage replan: %v\n", err)
		os.Exit(1)
	}
}