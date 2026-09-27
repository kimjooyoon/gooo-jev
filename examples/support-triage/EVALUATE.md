# declaration and feedback evaluation

The evaluate command joins a verified declaration pipeline with a verified
FeedbackWindow. It reports whether the pair is ready for observation; it does
not claim semantic completion or authorize execution.

~~~text
go run ./cmd/gooo-feedback feedback.json > /tmp/feedback-window.json
go run ./cmd/gooo-evaluate declaration.gooo /tmp/feedback-window.json
~~~

Declaration UNKNOWN, missing generation stages, invalid feedback evidence,
and zero capability boundaries remain explicit. A non-ready result exits with
status 2 after emitting its structured report.