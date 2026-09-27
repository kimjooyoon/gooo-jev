# Feedback window

The feedback window CLI keeps known quality separate from missing observations.

Prepare a JSON document with the fields MetricName and Observations, using the
serialized FeedbackObservation values produced by the Go API. Then run:

    go run ./cmd/gooo-feedback feedback.json

The output contains KnownMean for confirmed and refuted observations,
ObservedCoverage for known observations divided by all observations, and
UnknownCount for missing feedback. Unknown observations never contribute to
KnownMean.
