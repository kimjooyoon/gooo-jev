package decision

import "testing"

func TestProjectJEVActionCandidateLSPClearAndReview(t *testing.T) {
	verified := JEVActionCandidateVerification{
		Status: "verified", CandidateID: "candidate-1", EvidenceDigest: "verification-evidence",
		NonExecuting: true, NonAuthorizing: true,
	}
	output := ProjectJEVActionCandidateLSP(JEVActionCandidateLSPProjectionInput{
		Verification:         verified,
		EvidencePrefixDigest: "prefix-digest",
		NonAuthorizing:       true,
	})
	if output.Status != "clear" || output.Publishable || output.Severity != "info" || output.Code != "jev-candidate-verified" {
		t.Fatalf("unexpected clear output: %+v", output)
	}

	review := JEVActionCandidateVerification{
		Status: "review", CandidateID: "candidate-1", EvidenceDigest: "verification-evidence",
		MissingStage: "reverse-observation", FirstMismatch: "evidence_digest",
		NonExecuting: true, NonAuthorizing: true,
	}
	output = ProjectJEVActionCandidateLSP(JEVActionCandidateLSPProjectionInput{
		Verification:         review,
		MissingStageIndex:    3,
		EvidencePrefixDigest: "prefix-digest",
		NonAuthorizing:       true,
	})
	if output.Status != "publishable" || !output.Publishable || output.Severity != "error" || output.Code != "jev-reverse-counterexample" || output.MissingStageIndex != 3 {
		t.Fatalf("unexpected counterexample output: %+v", output)
	}
}

func TestProjectJEVActionCandidateLSPKeepsIncompleteEvidenceUnknown(t *testing.T) {
	review := JEVActionCandidateVerification{
		Status: "review", CandidateID: "candidate-1", EvidenceDigest: "verification-evidence",
		MissingStage: "candidate-guard", NonExecuting: true, NonAuthorizing: true,
	}
	output := ProjectJEVActionCandidateLSP(JEVActionCandidateLSPProjectionInput{
		Verification:         review,
		EvidencePrefixDigest: "prefix-digest",
		NonAuthorizing:       true,
	})
	if output.Status != "UNKNOWN" || output.Publishable || output.Code != "lsp-diagnostic-evidence" {
		t.Fatalf("unexpected incomplete output: %+v", output)
	}

	unknown := review
	unknown.Status = "UNKNOWN"
	unknown.MissingStage = "reverse-evidence"
	output = ProjectJEVActionCandidateLSP(JEVActionCandidateLSPProjectionInput{
		Verification:         unknown,
		MissingStageIndex:    1,
		EvidencePrefixDigest: "prefix-digest",
		NonAuthorizing:       true,
	})
	if output.Status != "UNKNOWN" || output.Publishable {
		t.Fatalf("unexpected unknown output: %+v", output)
	}
}

func TestProjectJEVActionCandidateLSPFailsClosedOnAuthorization(t *testing.T) {
	output := ProjectJEVActionCandidateLSP(JEVActionCandidateLSPProjectionInput{
		Verification: JEVActionCandidateVerification{
			Status: "verified", CandidateID: "candidate-1", EvidenceDigest: "verification-evidence",
			NonExecuting: true, NonAuthorizing: true,
		},
		EvidencePrefixDigest: "prefix-digest",
		NonAuthorizing:       false,
	})
	if output.Status != "UNKNOWN" || output.NonAuthorizing || output.Code != "authorization-boundary" {
		t.Fatalf("unexpected authorization output: %+v", output)
	}
}
