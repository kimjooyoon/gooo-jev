package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/kimjooyoon/gooo-jev/gooo"
)

type discoveryDocument struct {
	Trail gooo.CapabilityQueryTrail
	Guide gooo.CapabilityQueryGuide
}

func main() {
	var question string
	jsonOutput := flag.Bool("json", false, "emit provenance-bound JSON instead of readable guidance")
	flag.StringVar(&question, "question", "", "natural-language question about what gooo can do")
	flag.StringVar(&question, "q", "", "short form of -question")
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "usage: go run ./cmd/gooo-capabilities [-q question] [declaration.gooo|-]")
	}
	flag.Parse()
	if flag.NArg() > 1 {
		flag.Usage()
		os.Exit(64)
	}

	var source []byte
	if flag.NArg() == 1 {
		var err error
		if flag.Arg(0) == "-" {
			source, err = io.ReadAll(os.Stdin)
		} else {
			source, err = os.ReadFile(filepath.Clean(flag.Arg(0)))
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "read declaration: %v\n", err)
			os.Exit(1)
		}
	}
	if question == "" && len(source) > 0 {
		question = "What can gooo do with this declaration?"
	}
	if question == "" {
		flag.Usage()
		os.Exit(64)
	}

	trail := gooo.DiscoverCapabilityQueryTrail(question, string(source))
	if err := trail.Validate(); err != nil {
		fmt.Fprintf(os.Stderr, "validate capability query trail: %v\n", err)
		os.Exit(1)
	}
	guide := gooo.DiscoverCapabilityQueryGuide(question, string(source))
	if err := guide.Validate(); err != nil {
		fmt.Fprintf(os.Stderr, "validate capability query guide: %v\n", err)
		os.Exit(1)
	}

	if *jsonOutput {
		if err := json.NewEncoder(os.Stdout).Encode(discoveryDocument{Trail: trail, Guide: guide}); err != nil {
			fmt.Fprintf(os.Stderr, "write capability discovery: %v\n", err)
			os.Exit(1)
		}
		return
	}
	rendered, err := gooo.RenderCapabilityQuery(trail, guide)
	if err != nil {
		fmt.Fprintf(os.Stderr, "render capability discovery: %v\n", err)
		os.Exit(1)
	}
	fmt.Print(rendered)
}
