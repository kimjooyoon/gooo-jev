package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/kimjooyoon/gooo-jev/gooo"
)

type report struct {
	Path     string
	Response gooo.SyntaxCompletionResponse
}

func main() {
	if len(os.Args) < 2 || len(os.Args) > 3 {
		fmt.Fprintln(os.Stderr, "usage: go run ./cmd/gooo-complete <declaration.gooo> [prefix]")
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

	response := gooo.CompleteSyntax(string(source), prefix)
	if err := response.Validate(); err != nil {
		fmt.Fprintf(os.Stderr, "invalid completion evidence: %v\n", err)
		os.Exit(1)
	}

	output := report{
		Path:     path,
		Response: response,
	}
	if err := json.NewEncoder(os.Stdout).Encode(output); err != nil {
		fmt.Fprintf(os.Stderr, "write completion: %v\n", err)
		os.Exit(1)
	}
}