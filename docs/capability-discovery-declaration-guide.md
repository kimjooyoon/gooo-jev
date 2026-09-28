# Declaration-aware capability discovery

The capability discovery path answers a natural question about a .gooo declaration without pretending that the language can execute or authorize the requested work.

## User flow

1. Start with a declaration file and ask a broad question such as:
   "What can gooo do with this declaration?"
2. The discovery layer extracts only declaration signals that are present, such as module, input, output, status, field, constraint, observe, transform, contract, policy, or workflow.
3. It returns focused follow-up questions and a capability trail. The trail preserves the original question, declaration binding, evidence digest, and first missing stage.
4. The user can inspect the suggested next operation before any execution or authorization boundary is considered.

## CLI shape

    gooo-capabilities -q "What can gooo do with this declaration?" path/to/program.gooo

Use -json when another tool needs the structured trail. The structured response is intended to feed provenance-aware consumers such as jev-gooo; it is not an execution request.

## Safety and provenance rules

- A declaration signal narrows discovery; it does not prove that an implementation exists.
- AVAILABLE means that the described capability is available for inspection, not that it was executed.
- DEFERRED and UNKNOWN remain visible with their first missing stage.
- The original user query is preserved even when declaration-aware suggestions are added.
- Source, declaration, and evidence digests stay bound to the response.
- Discovery never authorizes an operation, grants a capability, or hides an unresolved stage.

## Example questions

- What can gooo inspect for this declaration?
- What input and output contracts are visible here?
- Which constraints can be checked next?
- What provenance evidence is still missing?

This makes the language useful before a full implementation exists: the declaration itself becomes a safe, inspectable starting point for choosing the next operation.
