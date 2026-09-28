# Capability Discovery Workflow

This workflow answers a bounded question such as `What can gooo do with this declaration?` without executing the declaration or granting access.

## Flow

1. The client sends a query and an optional `.gooo` source document.
2. The discovery layer binds the query, source, declaration, catalog, schema, evidence, and toolchain digests.
3. LSP hover presents the capabilities supported by the bound declaration and identifies the first missing stage for any incomplete result.
4. The action projection maps the result to one non-executing next operation:
   - `AVAILABLE`: inspect the next supported operation.
   - `DEFERRED`: declare or measure the named missing stage.
   - `UNKNOWN`: repair or re-bind provenance before reasoning about capability.
5. The JEV observation envelope records the same state, digests, missing stage, and boundary flags for later replay.
6. A later observation may confirm or refute the usefulness of the suggested operation. It must carry the prior evidence digest.

## Invariants

- Capability discovery never executes a `.gooo` activity.
- Capability discovery never authorizes a user, network, credential, or catalog mutation.
- `UNKNOWN` never becomes `AVAILABLE` through cache presence or an inferred default.
- The first missing stage is retained until the missing evidence is supplied or the receipt is explicitly invalidated.
- A capability count is not a completeness claim.
- A successful LSP presentation is not proof that generated code is semantically correct.

## Example interpretation

If syntax completion is `AVAILABLE` but canonical generation is `DEFERRED` at `generation_digest`, the useful answer is not “gooo can generate code.” The answer is: syntax completion is supported, generation evidence is missing, and the next safe operation is to bind generation evidence. That answer is actionable and preserves the boundary between discovery and execution.
