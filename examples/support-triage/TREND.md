# feedback window trend

Compare two validated FeedbackWindow files without treating unknown
observations as quality. The report keeps known mean, coverage, unknown
counts, source evidence digests, and an explicit comparison time.

~~~text
go run ./cmd/gooo-trend previous-window.json current-window.json 2026-09-28T00:00:00Z
~~~

A rising known mean with falling coverage is not an unqualified improvement.
The trend is observation-only and non-authorizing.