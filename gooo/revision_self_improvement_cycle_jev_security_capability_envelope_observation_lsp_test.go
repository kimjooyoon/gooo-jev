package gooo

import "testing"

func TestProjectJEVCapabilityEnvelopeObservationLSPBound(t *testing.T) {
	projection := ProjectJEVCapabilityEnvelopeObservationLSP(
		JEVCapabilityEnvelopeObservationLSPInput{
			Status:                   "BOUND",
			WorkloadIdentityStatus:   "BOUND",
			EvidencePrefixDigest:     "sha256:prefix",
			PlanDigest:               "sha256:plan",
			NetworkAllowlistDigest:   "sha256:allowlist",
			RequestedCapabilities:    []string{"fs.read", "net.read"},
			ObservedCapabilities:     []string{"fs.read", "net.read"},
			TargetStage:              "capability_envelope_observation",
			ReverseObservationDigest: "sha256:reverse",
			ObservationDigest:        "sha256:observation",
		},
	)

	if projection.Status != JEVCapabilityEnvelopeObservationLSPInformation {
		t.Fatalf("status = %q, want Information", projection.Status)
	}
	if projection.Code != "jev.security.capability_envelope.bound" {
		t.Fatalf("code = %q, want jev.security.capability_envelope.bound", projection.Code)
	}
	if len(projection.RequestedCapabilities) != 2 || projection.ObservedCapabilities[1] != "net.read" {
		t.Fatal("capability evidence was not preserved")
	}
	if projection.ReverseObservationDigest != "sha256:reverse" {
		t.Fatalf("reverse digest = %q, want sha256:reverse", projection.ReverseObservationDigest)
	}
	if !projection.IsReadOnly || projection.CanEdit || projection.CanExecute || projection.CanAuthorize {
		t.Fatal("capability envelope LSP projection must remain read-only and non-authorizing")
	}
}

func TestProjectJEVCapabilityEnvelopeObservationLSPDeferred(t *testing.T) {
	projection := ProjectJEVCapabilityEnvelopeObservationLSP(
		JEVCapabilityEnvelopeObservationLSPInput{
			Status:                 "DEFERRED",
			TargetStage:            "reverse_observation",
			ObservationDigest:      "sha256:observation",
		},
	)

	if projection.Status != JEVCapabilityEnvelopeObservationLSPWarning {
		t.Fatalf("status = %q, want Warning", projection.Status)
	}
	if projection.Code != "jev.security.capability_envelope.deferred" {
		t.Fatalf("code = %q, want jev.security.capability_envelope.deferred", projection.Code)
	}
}

func TestProjectJEVCapabilityEnvelopeObservationLSPUnknown(t *testing.T) {
	projection := ProjectJEVCapabilityEnvelopeObservationLSP(
		JEVCapabilityEnvelopeObservationLSPInput{
			Status:      "UNKNOWN",
			TargetStage: "capability_set",
		},
	)

	if projection.Status != JEVCapabilityEnvelopeObservationLSPError {
		t.Fatalf("status = %q, want Error", projection.Status)
	}
	if projection.Code != "jev.security.capability_envelope.unknown" {
		t.Fatalf("code = %q, want jev.security.capability_envelope.unknown", projection.Code)
	}
}

func TestProjectJEVCapabilityEnvelopeObservationLSPRejectsIncompleteBoundEvidence(t *testing.T) {
	projection := ProjectJEVCapabilityEnvelopeObservationLSP(
		JEVCapabilityEnvelopeObservationLSPInput{
			Status:                 "BOUND",
			WorkloadIdentityStatus: "BOUND",
			EvidencePrefixDigest:   "sha256:prefix",
			PlanDigest:             "sha256:plan",
			NetworkAllowlistDigest: "sha256:allowlist",
			TargetStage:            "reverse_observation",
			ObservationDigest:      "sha256:observation",
		},
	)

	if projection.Status != JEVCapabilityEnvelopeObservationLSPError {
		t.Fatalf("status = %q, want Error", projection.Status)
	}
	if projection.Code != "jev.security.capability_envelope.unknown" {
		t.Fatalf("code = %q, want unknown code", projection.Code)
	}
}
