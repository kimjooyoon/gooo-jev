package main

import (
    "bufio"
    "encoding/json"
    "fmt"
    "os"

    "github.com/kimjooyoon/gooo-jev/gooo"
)

type request struct {
    ID string `json:"id"`
    Source string `json:"source"`
    Prefix string `json:"prefix"`
}

type problem struct {
    Code string `json:"code"`
    Message string `json:"message"`
}

type response struct {
    ID string `json:"id"`
    Completion *gooo.SyntaxCompletionResponse `json:"completion,omitempty"`
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
        if err := encoder.Encode(response{ID: input.ID, Completion: &completion}); err != nil {
            fmt.Fprintln(os.Stderr, err)
            os.Exit(1)
        }
    }
    if err := scanner.Err(); err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
}