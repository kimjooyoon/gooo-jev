package gooo

import "testing"

func selfImprovementReceiptWithProfile(t *testing.T, bytes, lines int, irChanged, exactSourceMatch, structureMatch bool) RevisionSelfImprovementReceipt {
	t.Helper()
	applicationObservation, metricsBinding, applicationAssessment, generationAssessment := selfImprovementReceiptInputs(t)
	receipt, err := ObserveRevisionSelfImprovementReceipt(applicationObservation, metricsBinding, applicationAssessment, generationAssessment)
	if err != nil {
		t.Fatalf("ObserveRevisionSelfImprovementReceipt() error = %v", err)
	}
	receipt.ChangedByteCount = bytes
	receipt.ChangedLineCount = lines
	receipt.IRChanged = irChanged
	receipt.ExactSourceMatch = exactSourceMatch
	receipt.StructureMatch = structureMatch
	receipt.ReceiptDigest = digestRevisionSelfImprovementReceipt(receipt)
	return receipt
}

func TestObserveRevisionSelfImprovementWindowStable(t *testing.T) {
	baseline := selfImprovementReceiptWithProfile(t, 4, 1, false, true, true)
	candidate := selfImprovementReceiptWithProfile(t, 4, 1, false, true, true)
	window, err := ObserveRevisionSelfImprovementWindow(baseline, candidate)
	if err != nil {
		t.Fatalf("ObserveRevisionSelfImprovementWindow() error = %v", err)
	}
	if window.Status != "BOUND" || window.ComparisonSignal != "stable" ||
		window.ByteDelta != 0 || window.LineDelta != 0 {
		t.Fatalf("unexpected stable window: %#v", window)
	}
	if err := window.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
}

func TestObserveRevisionSelfImprovementWindowReportsNarrowerSurface(t *testing.T) {
	baseline := selfImprovementReceiptWithProfile(t, 8, 3, false, false, true)
	candidate := selfImprovementReceiptWithProfile(t, 2, 1, false, false, true)
	window, err := ObserveRevisionSelfImprovementWindow(baseline, candidate)
	if err != nil {
		t.Fatalf("ObserveRevisionSelfImprovementWindow() error = %v", err)
	}
	if window.ComparisonSignal != "narrower" || window.ByteDelta != -6 || window.LineDelta != -2 {
		t.Fatalf("unexpected narrower window: %#v", window)
	}
}

func TestObserveRevisionSelfImprovementWindowReportsMixedEvidence(t *testing.T) {
	baseline := selfImprovementReceiptWithProfile(t, 8, 3, false, true, true)
	candidate := selfImprovementReceiptWithProfile(t, 2, 1, true, true, true)
	window, err := ObserveRevisionSelfImprovementWindow(baseline, candidate)
	if err != nil {
		t.Fatalf("ObserveRevisionSelfImprovementWindow() error = %v", err)
	}
	if window.ComparisonSignal != "mixed" || window.IRChangeStable {
		t.Fatalf("unexpected mixed window: %#v", window)
	}
}

func TestObserveRevisionSelfImprovementWindowRetainsUnknownBaseline(t *testing.T) {
	baseline := selfImprovementReceiptWithProfile(t, 4, 1, false, true, true)
	baseline.ReceiptDigest = digestString("tampered")
	candidate := selfImprovementReceiptWithProfile(t, 4, 1, false, true, true)
	window, err := ObserveRevisionSelfImprovementWindow(baseline, candidate)
	if err == nil {
		t.Fatal("ObserveRevisionSelfImprovementWindow() error = nil, want baseline failure")
	}
	if window.Status != "UNKNOWN" || window.MissingStage != "revision-self-improvement-window-baseline" {
		t.Fatalf("unexpected unknown baseline window: %#v", window)
	}
}

func TestObserveRevisionSelfImprovementWindowIsDeterministic(t *testing.T) {
	baseline := selfImprovementReceiptWithProfile(t, 4, 1, false, true, true)
	candidate := selfImprovementReceiptWithProfile(t, 4, 1, false, true, true)
	first, err := ObserveRevisionSelfImprovementWindow(baseline, candidate)
	if err != nil {
		t.Fatalf("first ObserveRevisionSelfImprovementWindow() error = %v", err)
	}
	second, err := ObserveRevisionSelfImprovementWindow(baseline, candidate)
	if err != nil {
		t.Fatalf("second ObserveRevisionSelfImprovementWindow() error = %v", err)
	}
	if first.WindowDigest != second.WindowDigest {
		t.Fatal("same receipt pair produced different window digest")
	}
}