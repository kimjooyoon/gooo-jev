package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/kimjooyoon/gooo-jev/gooo"
)

func main() {
	jsonOutput := flag.Bool("json", false, "emit provenance-bound JSON instead of human-readable guidance")
	flag.Parse()
	if flag.NArg() < 1 || flag.NArg() > 2 {
		fmt.Fprintln(os.Stderr, "usage: go run ./cmd/gooo-assist [-json] <declaration.gooo> [prefix]")
		os.Exit(64)
	}
	path := filepath.Clean(flag.Arg(0))
	prefix := ""
	if flag.NArg() == 2 {
		prefix = flag.Arg(1)
	}
	source, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "read declaration: %v\n", err)
		os.Exit(1)
	}
	discovery := gooo.DiscoverUsage(string(source), prefix)
	plan, err := gooo.PlanUsageActions(discovery)
	if err != nil {
		fmt.Fprintf(os.Stderr, "plan usage actions: %v\n", err)
		os.Exit(1)
	}
	assistance, err := gooo.ExplainUsage(discovery, plan)
	if err != nil {
		fmt.Fprintf(os.Stderr, "explain usage: %v\n", err)
		os.Exit(1)
	}
	if *jsonOutput {
		if err := json.NewEncoder(os.Stdout).Encode(assistance); err != nil {
			fmt.Fprintf(os.Stderr, "write assistance: %v\n", err)
			os.Exit(1)
		}
		return
	}
	text, err := gooo.RenderUsageAssistance(assistance)
	if err != nil {
		fmt.Fprintf(os.Stderr, "render assistance: %v\n", err)
		os.Exit(1)
	}
	fmt.Print(text)
}
