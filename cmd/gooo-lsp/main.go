package main

import (
    "bufio"
    "encoding/json"
    "fmt"
    "os"
    "strings"

    "github.com/kimjooyoon/gooo-jev/gooo"
)

type request struct {
    ID string `json:"id"`
    Source string `json:"source"`
    Prefix string `json:"prefix"`
    Query string `json:"query"`
}

type problem struct {
    Code string `json:"code"`
    Message string `json:"message"`
}

type response struct {
    ID string `json:"id"`
    Completion *gooo.SyntaxCompletionResponse `json:"completion,omitempty"`
    Discovery *gooo.UsageDiscoveryResponse `json:"discovery,omitempty"`
    CapabilityDiscovery *gooo.CapabilityQueryResponse `json:"capability_discovery,omitempty"`
    Assistance *gooo.UsageAssistance `json:"assistance,omitempty"`
    Error *problem `json:"error,omitempty"`
}

func main() {
    scanner := bufio.NewScanner(os.Stdin)
    scanner.Buffer(make([]byte, 1024), 1024*1024)
    encoder := json.NewEncoder(os.Stdout)
    for scanner.Scan() {
        var input request
        if err := json.Unmarshal(scanner.Bytes(), &input); err != nil {
            if encodeErr := encoder.Encode(response{Error: &problem{Code: "invalid_request", Message: err.Error()}}); encodeErr != nil {
                fmt.Fprintln(os.Stderr, encodeErr)
                os.Exit(1)
            }
            continue
        }
        completion := gooo.CompleteSyntax(input.Source, input.Prefix)
        if err := completion.Validate(); err != nil {
            if encodeErr := encoder.Encode(response{ID: input.ID, Error: &problem{Code: "invalid_completion", Message: err.Error()}}); encodeErr != nil {
                fmt.Fprintln(os.Stderr, encodeErr)
                os.Exit(1)
            }
            continue
        }
        discovery := gooo.DiscoverUsage(input.Source, input.Prefix)
        plan, err := gooo.PlanUsageActions(discovery)
        if err != nil {
            if encodeErr := encoder.Encode(response{ID: input.ID, Error: &problem{Code: "invalid_usage_plan", Message: err.Error()}}); encodeErr != nil {
                fmt.Fprintln(os.Stderr, encodeErr)
                os.Exit(1)
            }
            continue
        }
        assistance, err := gooo.ExplainUsage(discovery, plan)
        if err != nil {
            if encodeErr := encoder.Encode(response{ID: input.ID, Error: &problem{Code: "invalid_usage_assistance", Message: err.Error()}}); encodeErr != nil {
                fmt.Fprintln(os.Stderr, encodeErr)
                os.Exit(1)
            }
            continue
        }
        query := strings.TrimSpace(input.Query)
        if query == "" {
            query = "What can gooo do with this declaration?"
        }
        capabilityDiscovery := gooo.DiscoverCapabilityQueryWithDeclaration(query, input.Source)
        if err := capabilityDiscovery.Validate(); err != nil {
            if encodeErr := encoder.Encode(response{ID: input.ID, Error: &problem{Code: "invalid_capability_discovery", Message: err.Error()}}); encodeErr != nil {
                fmt.Fprintln(os.Stderr, encodeErr)
                os.Exit(1)
            }
            continue
        }
        if err := encoder.Encode(response{ID: input.ID, Completion: &completion, Discovery: &discovery, CapabilityDiscovery: &capabilityDiscovery, Assistance: &assistance}); err != nil {
            fmt.Fprintln(os.Stderr, err)
            os.Exit(1)
        }
    }
    if err := scanner.Err(); err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
}
