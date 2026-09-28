package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/kimjooyoon/gooo-jev/decision"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "usage: go run ./cmd/gooo-decision-receipt [observation.json|-]")
	}
	flag.Parse()
	if flag.NArg() > 1 {
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

	receipt, err := decision.ObserveJSON(data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "observe decision: %v\n", err)
		os.Exit(1)
	}
	encoded, err := decision.MarshalReceiptJSON(receipt)
	if err != nil {
		fmt.Fprintf(os.Stderr, "marshal decision receipt: %v\n", err)
		os.Exit(1)
	}
	encoded = append(encoded, '\n')
	if _, err := os.Stdout.Write(encoded); err != nil {
		fmt.Fprintf(os.Stderr, "write decision receipt: %v\n", err)
		os.Exit(1)
	}
}
