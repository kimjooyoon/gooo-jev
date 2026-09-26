package decision

import (
	"strings"
	"testing"
)

const sampleGoooEvidenceDeclaration = "package jevdecision\n" +
	"namespace jevdecision\n\n" +
	"entity DecisionSpec id \"gooo://jev/decision/spec\"\n" +
	"  property Name string\n\n" +
	"activity ObserveDecision(DecisionSpec) -> DecisionReceipt\n" +
	"evidence FullProvenance from ObserveDecision\n"

func TestDeriveGoooEvidenceDeclarationIRGenerationRoundTrip(t *testing.T) {
	derived := DeriveGoooEvidenceDeclarationIRGeneration(GoooEvidenceDeclarationIRGenerationInput{
		SourceText:     sampleGoooEvidenceDeclaration,
		NonAuthorizing: true,
	})
	if derived.Status != "ready" || len(derived.Evidence) != 1 ||
		derived.Evidence[0].Name != "FullProvenance" ||
		derived.Evidence[0].Source != "ObserveDecision" ||
		derived.BaseIRDigest == "" || derived.EvidenceDigest == "" ||
		derived.GenerationDigest == "" {
		t.Fatalf("derived = %#v", derived)
	}
	if !strings.Contains(derived.GeneratedSource, "evidence FullProvenance from ObserveDecision") {
		t.Fatalf("generated source lost evidence declaration: %q", derived.GeneratedSource)
	}
	roundTrip := DeriveGoooEvidenceDeclarationIRGeneration(GoooEvidenceDeclarationIRGenerationInput{
		SourceText:     derived.GeneratedSource,
		NonAuthorizing: true,
	})
	if roundTrip.Status != "ready" || roundTrip.BaseIRDigest != derived.BaseIRDigest ||
		roundTrip.EvidenceDigest != derived.EvidenceDigest ||
		roundTrip.GenerationDigest != derived.GenerationDigest {
		t.Fatalf("round trip = %#v, first = %#v", roundTrip, derived)
	}
}

func TestDeriveGoooEvidenceDeclarationIRGenerationPreservesSyntaxUnknown(t *testing.T) {
	invalid := DeriveGoooEvidenceDeclarationIRGeneration(GoooEvidenceDeclarationIRGenerationInput{
		SourceText: "package broken\nnamespace broken\nevidence Missing\n",
		NonAuthorizing: true,
	})
	if invalid.Status != "UNKNOWN" || invalid.MissingStage != "syntax" ||
		invalid.GeneratedSource != "" {
		t.Fatalf("invalid = %#v, want syntax UNKNOWN", invalid)
	}
	duplicate := DeriveGoooEvidenceDeclarationIRGeneration(GoooEvidenceDeclarationIRGenerationInput{
		SourceText: "package broken\nnamespace broken\nentity Example id \"gooo://example\"\nevidence One from A\nevidence One from B\n",
		NonAuthorizing: true,
	})
	if duplicate.Status != "UNKNOWN" || duplicate.MissingStage != "syntax" {
		t.Fatalf("duplicate = %#v, want syntax UNKNOWN", duplicate)
	}
}

func TestDeriveGoooEvidenceDeclarationIRGenerationRejectsAuthorization(t *testing.T) {
	unknown := DeriveGoooEvidenceDeclarationIRGeneration(GoooEvidenceDeclarationIRGenerationInput{
		SourceText:     sampleGoooEvidenceDeclaration,
		NonAuthorizing: false,
	})
	if unknown.Status != "UNKNOWN" || unknown.MissingStage != "authorization-boundary" ||
		unknown.NonAuthorizing {
		t.Fatalf("unknown = %#v, want authorization UNKNOWN", unknown)
	}
}