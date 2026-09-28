package gooo

import (
	"encoding/json"
	"testing"
)

func TestSyntaxCompletionJSONRoundTripPreservesProvenance(t *testing.T) {
	source := `package support
namespace triage
entity ticket id "ticket"
activity assign(ticket) -> ticket
`
	original := CompleteSyntax(source, "ti")
	if err := original.Validate(); err != nil {
		t.Fatalf("validate original completion: %v", err)
	}

	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("marshal completion: %v", err)
	}
	var restored SyntaxCompletionResponse
	if err := json.Unmarshal(data, &restored); err != nil {
		t.Fatalf("unmarshal completion: %v", err)
	}
	if restored.SourceDigest != original.SourceDigest ||
		restored.IRDigest != original.IRDigest ||
		restored.ItemsDigest != original.ItemsDigest ||
		len(restored.Items) != len(original.Items) {
		t.Fatalf("completion provenance changed across JSON round trip: original=%+v restored=%+v", original, restored)
	}
	if err := restored.Validate(); err != nil {
		t.Fatalf("validate restored completion: %v", err)
	}
}