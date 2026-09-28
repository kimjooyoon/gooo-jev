# gooo-jev
Provider-neutral typed decision receipts for Go and .gooo execution plans

## Natural capability discovery

Ask what a .gooo declaration can currently expose without invoking a
provider or granting execution authority. The answer is catalog-bound and
keeps AVAILABLE, DEFERRED, and UNKNOWN distinct.

    go run ./cmd/gooo-capabilities -q "What can gooo do?" declaration.gooo
    cat declaration.gooo | go run ./cmd/gooo-capabilities -q "이 선언으로 무엇을 할 수 있어?" -
    go run ./cmd/gooo-capabilities -json -q "Can gooo inspect provenance?" declaration.gooo

The readable output includes matched capabilities, a next operation, example
questions, declaration signals, constraints, and query/trail/guide digests.
Execution and authorization requests remain deferred until an explicit
external boundary is supplied.

## Typed decision JSON boundary

Use the provider-neutral boundary when a JEV-like provider returns a bounded
decision. The input contains a declaration of the question, a state digest,
and a typed result. The result is validated before a non-authorizing receipt
is emitted.

    go run ./cmd/gooo-decision-receipt observation.json
    cat observation.json | go run ./cmd/gooo-decision-receipt -

The boundary accepts choice, score, and noul observations. Unknown fields,
trailing JSON, invalid choices, and incomplete provenance are rejected. A
receipt never grants execution, changes source, or replaces the later
capability envelope and reverse-observation checks.
