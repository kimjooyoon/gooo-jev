package decision

import "testing"

const sampleGoooDeclaration = "package jevdecision\n" +
	"namespace jevdecision\n\n" +
	"entity DecisionSpec id \"gooo://jev/decision/spec\"\n" +
	"  property Name string\n\n" +
	"entity DecisionReceipt id \"gooo://jev/decision/receipt\"\n\n" +
	"activity ObserveDecision(DecisionSpec) -> DecisionReceipt\n"

func TestDeriveGoooDeclarationIRGenerationRoundTrip(t *testing.T) {
	derived := DeriveGoooDeclarationIRGeneration(GoooDeclarationIRGenerationInput{
		SourceText:     sampleGoooDeclaration,
		NonAuthorizing: true,
	})
	if derived.Status != "ready" || derived.MissingStage != "" || derived.SourceDigest == "" || derived.IRDigest == "" || derived.GenerationDigest == "" || derived.GeneratedSource == "" ||
		derived.ReverseObservationStatus != "ready" || derived.RoundTripIRDigest != derived.IRDigest || derived.ReverseObservationDigest == "" {
		t.Fatalf("derived = %#v", derived)
	}
	if err := derived.Validate(); err != nil {
		t.Fatalf("derived validation error = %v", err)
	}
	roundTrip := DeriveGoooDeclarationIRGeneration(GoooDeclarationIRGenerationInput{
		SourceText:     derived.GeneratedSource,
		NonAuthorizing: true,
	})
	if roundTrip.Status != "ready" || roundTrip.IRDigest != derived.IRDigest || roundTrip.GenerationDigest != derived.GenerationDigest {
		t.Fatalf("round trip = %#v, first = %#v", roundTrip, derived)
	}
	if len(derived.IR.Entities) != 2 || len(derived.IR.Activities) != 1 || derived.IR.Activities[0].Parameters[0].Type != "DecisionSpec" {
		t.Fatalf("unexpected IR = %#v", derived.IR)
	}
}

func TestDeriveGoooDeclarationIRGenerationRejectsTamperedReverseObservation(t *testing.T) {
	derived := DeriveGoooDeclarationIRGeneration(GoooDeclarationIRGenerationInput{
		SourceText:     sampleGoooDeclaration,
		NonAuthorizing: true,
	})
	derived.RoundTripIRDigest = "tampered-round-trip-ir"
	if err := derived.Validate(); err == nil {
		t.Fatal("tampered round-trip IR must fail validation")
	}
	derived = DeriveGoooDeclarationIRGeneration(GoooDeclarationIRGenerationInput{
		SourceText:     sampleGoooDeclaration,
		NonAuthorizing: true,
	})
	derived.ReverseObservationDigest = "tampered-reverse-observation"
	if err := derived.Validate(); err == nil {
		t.Fatal("tampered reverse observation must fail validation")
	}
}

func TestDeriveGoooDeclarationIRGenerationPreservesUnknownStage(t *testing.T) {
	invalid := DeriveGoooDeclarationIRGeneration(GoooDeclarationIRGenerationInput{
		SourceText:     "package broken\nnamespace broken\nproperty Missing string\n",
		NonAuthorizing: true,
	})
	if invalid.Status != "UNKNOWN" || invalid.MissingStage != "syntax" || invalid.GeneratedSource != "" {
		t.Fatalf("invalid = %#v", invalid)
	}
	missing := DeriveGoooDeclarationIRGeneration(GoooDeclarationIRGenerationInput{
		SourceText:     "package broken\nnamespace broken\n",
		NonAuthorizing: true,
	})
	if missing.Status != "UNKNOWN" || missing.MissingStage != "syntax" {
		t.Fatalf("missing declaration = %#v", missing)
	}
}

func TestDeriveGoooDeclarationIRGenerationRejectsAuthorization(t *testing.T) {
	unknown := DeriveGoooDeclarationIRGeneration(GoooDeclarationIRGenerationInput{
		SourceText:     sampleGoooDeclaration,
		NonAuthorizing: false,
	})
	if unknown.Status != "UNKNOWN" || unknown.MissingStage != "authorization-boundary" || unknown.NonAuthorizing {
		t.Fatalf("unknown = %#v", unknown)
	}
}
