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
    CapabilityTrail *gooo.CapabilityQueryTrail `json:"capability_trail,omitempty"`
    CapabilityOverview *gooo.CapabilityQueryOverview `json:"capability_overview,omitempty"`
    CapabilityGuide *gooo.CapabilityQueryGuide `json:"capability_guide,omitempty"`
    CapabilityHover *gooo.CapabilityQueryHover `json:"capability_hover,omitempty"`
    CapabilityActions []gooo.CapabilityQueryAction `json:"capability_actions,omitempty"`
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
        capabilityTrail := gooo.DiscoverCapabilityQueryTrail(query, input.Source)
        if err := capabilityTrail.Validate(); err != nil {
            if encodeErr := encoder.Encode(response{ID: input.ID, Error: &problem{Code: "invalid_capability_trail", Message: err.Error()}}); encodeErr != nil {
                fmt.Fprintln(os.Stderr, encodeErr)
                os.Exit(1)
            }
            continue
        }
        capabilityOverview := gooo.DiscoverCapabilityQueryOverview(query, input.Source)
        if err := capabilityOverview.Validate(); err != nil {
            if encodeErr := encoder.Encode(response{ID: input.ID, Error: &problem{Code: "invalid_capability_overview", Message: err.Error()}}); encodeErr != nil {
                fmt.Fprintln(os.Stderr, encodeErr)
                os.Exit(1)
            }
            continue
        }
        capabilityGuide := gooo.DiscoverCapabilityQueryGuide(query, input.Source)
        if err := capabilityGuide.Validate(); err != nil {
            if encodeErr := encoder.Encode(response{ID: input.ID, Error: &problem{Code: "invalid_capability_guide", Message: err.Error()}}); encodeErr != nil {
                fmt.Fprintln(os.Stderr, encodeErr)
                os.Exit(1)
            }
            continue
        }
        capabilityHover := gooo.DiscoverCapabilityQueryHover(query, input.Source)
        if err := capabilityHover.Validate(); err != nil {
            if encodeErr := encoder.Encode(response{ID: input.ID, Error: &problem{Code: "invalid_capability_hover", Message: err.Error()}}); encodeErr != nil {
                fmt.Fprintln(os.Stderr, err)
                os.Exit(1)
            }
            continue
        }
        capabilityActions, err := gooo.SuggestCapabilityQueryActions(capabilityHover)
        if err != nil {
            if encodeErr := encoder.Encode(response{ID: input.ID, Error: &problem{Code: "invalid_capability_actions", Message: err.Error()}}); encodeErr != nil {
                fmt.Fprintln(os.Stderr, encodeErr)
                os.Exit(1)
            }
            continue
        }
        if err := gooo.ValidateCapabilityQueryActions(capabilityHover, capabilityActions); err != nil {
            if encodeErr := encoder.Encode(response{ID: input.ID, Error: &problem{Code: "invalid_capability_actions", Message: err.Error()}}); encodeErr != nil {
                fmt.Fprintln(os.Stderr, err)
                os.Exit(1)
            }
            continue
        }
        capabilityDiscovery := capabilityTrail.Response
        if err := encoder.Encode(response{ID: input.ID, Completion: &completion, Discovery: &discovery, CapabilityDiscovery: &capabilityDiscovery, CapabilityTrail: &capabilityTrail, CapabilityOverview: &capabilityOverview, CapabilityGuide: &capabilityGuide, CapabilityHover: &capabilityHover, CapabilityActions: capabilityActions, Assistance: &assistance}); err != nil {
            fmt.Fprintln(os.Stderr, err)
            os.Exit(1)
        }
    }
    if err := scanner.Err(); err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
}
