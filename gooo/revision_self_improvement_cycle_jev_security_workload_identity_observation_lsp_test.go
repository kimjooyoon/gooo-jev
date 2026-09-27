package gooo

import "testing"

func TestProjectJEVWorkloadIdentityObservationLSPBound(t *testing.T) {
	projection := ProjectJEVWorkloadIdentityObservationLSP(
		JEVWorkloadIdentityObservationLSPInput{
			Status:               "BOUND",
			SpiffeID:             "spiffe://example.test/ns/gooo/sa/jev",
			Audience:             "jev-runtime",
			Capabilities:         []string{"provenance.read", "lsp.observe"},
			EvidencePrefixDigest: "prefix-digest",
		},
	)

	if projection.Status != JEVWorkloadIdentityObservationLSPInformation {
		t.Fatalf("status = %q, want Information", projection.Status)
	}
	if projection.Code != "jev.security.workload_identity.bound" {
		t.Fatalf("code = %q, want jev.security.workload_identity.bound", projection.Code)
	}
	if len(projection.Capabilities) != 2 || projection.Capabilities[1] != "lsp.observe" {
		t.Fatal("capabilities were not preserved")
	}
	if !projection.IsReadOnly || projection.CanEdit || projection.CanExecute || projection.CanAuthorize {
		t.Fatal("identity LSP projection must remain read-only and non-authorizing")
	}
}

func TestProjectJEVWorkloadIdentityObservationLSPUnknown(t *testing.T) {
	projection := ProjectJEVWorkloadIdentityObservationLSP(
		JEVWorkloadIdentityObservationLSPInput{
			Status: "UNKNOWN",
		},
	)

	if projection.Status != JEVWorkloadIdentityObservationLSPError {
		t.Fatalf("status = %q, want Error", projection.Status)
	}
}