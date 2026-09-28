# declaration evaluation with trend

This command joins declaration provenance with two validated feedback
windows. It emits the trend digest and an observation readiness boundary.
It does not authorize execution or treat a rising mean as completion.

~~~text
go run ./cmd/gooo-evaluate-trend declaration.gooo previous-window.json current-window.json 2026-09-28T00:00:00Z
~~~

Coverage and unknown counts remain visible beside the mean delta, so a
higher known mean with lower coverage is not silently treated as progress.