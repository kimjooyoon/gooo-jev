# gooo editor completion stream

The gooo-lsp command provides a small editor adapter without requiring a full
LSP transport. It reads one JSON object per line and emits one JSON object per
line, preserving the completion response and its provenance digests.

From the repository root:

~~~text
printf "%s\\n" "{\"id\":\"1\",\"source\":\"package support\\nnamespace triage\\nentity ticket\\nproperty\",\"prefix\":\"pro\"}" | go run ./cmd/gooo-lsp
~~~

A response can remain UNKNOWN while still returning useful candidates. The
adapter is non-executing and non-authorizing; malformed requests are returned as
structured errors so an editor can keep the stream alive.