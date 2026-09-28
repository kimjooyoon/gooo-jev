# capability action planning

The plan command converts capability discovery into explicit next steps.
READY actions describe already-bound, non-executing operations. DEFERRED
actions preserve UNKNOWN reasons and the declaration operation needed next.

~~~text
go run ./cmd/gooo-discover declaration.gooo pro > discovery.json
go run ./cmd/gooo-plan discovery.json
~~~

The plan is provenance-bound and non-authorizing. It proposes work; it does
not execute a declaration or claim that the domain is complete.