package gooo

import "testing"

func boundaryMetric(t *testing.T, confidence int64, kind string) RevisionSelfImprovementCycleJEVDecisionConfidenceMetricObservation {
	t.Helper()
	signal, err := ObserveRevisionSelfImprovementCycleJEVTypedDecisionSignal(
		RevisionSelfImprovementCycleJEVTypedDecisionSignalObservationInput{
			QuestionID: "question:route", QuestionKind: kind,
			DecisionDigest: digestString("decision:route"),
			EvidenceDigest: digestString("evidence:route"), ConfidenceMilli: confidence,
		},
	)
	if err != nil { t.Fatalf("source signal returned an error: %v", err) }
	metric, err := ObserveRevisionSelfImprovementCycleJEVDecisionConfidenceMetric(signal)
	if err != nil { t.Fatalf("confidence metric returned an error: %v", err) }
	return metric
}

func boundaryPolicy() RevisionSelfImprovementCycleJEVDecisionAbstainBoundaryPolicy {
	return RevisionSelfImprovementCycleJEVDecisionAbstainBoundaryPolicy{
		AcceptThresholdMilli: 900, AbstainLowerMilli: 500, AbstainUpperMilli: 899,
		GenerationDigest: digestString("generation:route"),
		ReverseObservationDigest: digestString("reverse:route"),
	}
}

func TestObserveRevisionSelfImprovementCycleJEVDecisionAbstainBoundaryRoutesHighConfidence(t *testing.T) {
	observation := ObserveRevisionSelfImprovementCycleJEVDecisionAbstainBoundary(boundaryMetric(t, 920, "choice"), boundaryPolicy())
	if observation.Status != "BOUND" || observation.Disposition != "direct_filter_candidate" ||
		observation.CalibrationStatus != "unverified" || !observation.CalibrationRequired {
		t.Fatalf("unexpected high-confidence boundary: %#v", observation)
	}
	if !observation.NonExecuting || !observation.NonAuthorizing {
		t.Fatalf("boundary crossed an authority boundary: %#v", observation)
	}
	if err := observation.Validate(); err != nil { t.Fatalf("boundary did not validate: %v", err) }
}

func TestObserveRevisionSelfImprovementCycleJEVDecisionAbstainBoundaryAbstainsInMiddleBand(t *testing.T) {
	observation := ObserveRevisionSelfImprovementCycleJEVDecisionAbstainBoundary(boundaryMetric(t, 650, "choice"), boundaryPolicy())
	if observation.Status != "BOUND" || observation.Disposition != "abstain_for_calibration" {
		t.Fatalf("unexpected abstain boundary: %#v", observation)
	}
}

func TestObserveRevisionSelfImprovementCycleJEVDecisionAbstainBoundaryPreservesUnknown(t *testing.T) {
	observation := ObserveRevisionSelfImprovementCycleJEVDecisionAbstainBoundary(
		RevisionSelfImprovementCycleJEVDecisionConfidenceMetricObservation{}, boundaryPolicy())
	if observation.Status != "UNKNOWN" || observation.MissingStage == "" ||
		observation.Disposition != "unknown" || observation.CalibrationStatus != "unverified" {
		t.Fatalf("unknown boundary must preserve unresolved state: %#v", observation)
	}
	if err := observation.Validate(); err != nil { t.Fatalf("unknown boundary should validate: %v", err) }
}

func TestRevisionSelfImprovementCycleJEVDecisionAbstainBoundaryRejectsTampering(t *testing.T) {
	observation := ObserveRevisionSelfImprovementCycleJEVDecisionAbstainBoundary(boundaryMetric(t, 920, "choice"), boundaryPolicy())
	observation.Disposition = "execute"
	if err := observation.Validate(); err == nil {
		t.Fatal("expected authority-crossing disposition to be rejected")
	}
}

func TestRevisionSelfImprovementCycleJEVDecisionAbstainBoundaryRejectsNoulConfidence(t *testing.T) {
	observation := ObserveRevisionSelfImprovementCycleJEVDecisionAbstainBoundary(boundaryMetric(t, -1, "noul"), boundaryPolicy())
	if observation.Status != "UNKNOWN" || observation.MissingStage == "" {
		t.Fatalf("noul boundary must remain unknown: %#v", observation)
	}
	if err := observation.Validate(); err != nil { t.Fatalf("unknown noul boundary should validate: %v", err) }
}