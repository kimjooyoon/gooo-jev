package gooo

import "fmt"

// ProvenanceEdge records one immutable relationship between two evidence nodes.
type ProvenanceEdge struct {
	FromKind   string
	FromDigest string
	ToKind     string
	ToDigest   string
	Relation   string
}

// LineageReceipt binds source, decision IR, assessment, evidence, and revision candidate.
type LineageReceipt struct {
	Status           string
	MissingStage     string
	SourceDigest     string
	DecisionIRDigest string
	DecisionDigest   string
	AssessmentDigest string
	EvidenceDigest   string
	CandidateDigest  string
	ChainDigest      string
	Edges            []ProvenanceEdge
	NonExecuting     bool
	NonAuthorizing   bool
}

// ObserveLineage constructs a verifiable provenance chain without mutating source.
func ObserveLineage(document DecisionDocumentIR, assessment DecisionAssessment, candidate RevisionCandidate) (LineageReceipt, error) {
	receipt := LineageReceipt{
		Status:           "UNKNOWN",
		MissingStage:     "lineage-observation",
		SourceDigest:     document.SourceDigest,
		DecisionIRDigest: document.IRDigest,
		AssessmentDigest: assessment.AssessmentDigest,
		EvidenceDigest:   assessment.EvidenceDigest,
		CandidateDigest:  candidate.CandidateDigest,
		NonExecuting:     true,
		NonAuthorizing:   true,
	}
	if err := document.Validate(); err != nil {
		receipt.MissingStage = "lineage-document"
		return receipt, fmt.Errorf("gooo lineage: document: %w", err)
	}
	if err := assessment.Validate(); err != nil {
		receipt.MissingStage = "lineage-assessment"
		return receipt, fmt.Errorf("gooo lineage: assessment: %w", err)
	}
	if err := candidate.Validate(); err != nil {
		receipt.MissingStage = "lineage-candidate"
		return receipt, fmt.Errorf("gooo lineage: candidate: %w", err)
	}

	var declaration *DecisionDecl
	for index := range document.Decisions {
		if document.Decisions[index].Name == assessment.DecisionName {
			declaration = &document.Decisions[index]
			break
		}
	}
	if declaration == nil ||
		declaration.ID != assessment.DecisionID ||
		declaration.Digest != assessment.DecisionDigest ||
		candidate.DecisionName != assessment.DecisionName ||
		candidate.DecisionID != assessment.DecisionID ||
		candidate.AssessmentDigest != assessment.AssessmentDigest ||
		candidate.EvidenceDigest != assessment.EvidenceDigest {
		receipt.MissingStage = "lineage-binding"
		return receipt, fmt.Errorf("gooo lineage: assessment and candidate are not bound to the same declaration")
	}
	receipt.DecisionDigest = declaration.Digest
	receipt.Edges = []ProvenanceEdge{
		{FromKind: "gooo.source", FromDigest: document.SourceDigest, ToKind: "decision.ir", ToDigest: document.IRDigest, Relation: "parsed-into"},
		{FromKind: "decision.declaration", FromDigest: declaration.Digest, ToKind: "decision.assessment", ToDigest: assessment.AssessmentDigest, Relation: "observed-as"},
		{FromKind: "decision.evidence", FromDigest: assessment.EvidenceDigest, ToKind: "decision.assessment", ToDigest: assessment.AssessmentDigest, Relation: "supports"},
		{FromKind: "decision.assessment", FromDigest: assessment.AssessmentDigest, ToKind: "revision.candidate", ToDigest: candidate.CandidateDigest, Relation: "proposes"},
	}
	receipt.Status = "BOUND"
	receipt.MissingStage = ""
	receipt.ChainDigest = digestLineage(receipt)
	return receipt, nil
}

// Validate checks every edge and the chain digest.
func (r LineageReceipt) Validate() error {
	if r.Status != "BOUND" {
		return fmt.Errorf("lineage status must be BOUND")
	}
	if r.MissingStage != "" {
		return fmt.Errorf("lineage missing stage must be empty")
	}
	if !r.NonExecuting || !r.NonAuthorizing {
		return fmt.Errorf("lineage must remain non-executing and non-authorizing")
	}
	for name, digest := range map[string]string{
		"source": r.SourceDigest, "decision IR": r.DecisionIRDigest, "decision": r.DecisionDigest,
		"assessment": r.AssessmentDigest, "evidence": r.EvidenceDigest, "candidate": r.CandidateDigest,
		"chain": r.ChainDigest,
	} {
		if !validDigest(digest) {
			return fmt.Errorf("lineage %s digest is invalid", name)
		}
	}
	if len(r.Edges) != 4 {
		return fmt.Errorf("lineage edge count = %d, want 4", len(r.Edges))
	}
	for _, edge := range r.Edges {
		if edge.FromKind == "" || edge.FromDigest == "" || edge.ToKind == "" || edge.ToDigest == "" || edge.Relation == "" {
			return fmt.Errorf("lineage edge is incomplete")
		}
		if !validDigest(edge.FromDigest) || !validDigest(edge.ToDigest) {
			return fmt.Errorf("lineage edge digest is invalid")
		}
	}
	if expected := digestLineage(r); expected != r.ChainDigest {
		return fmt.Errorf("lineage chain digest does not match its edges")
	}
	return nil
}

func digestLineage(receipt LineageReceipt) string {
	digest := fmt.Sprintf("%s|%s|%s|%s|%s|%s", receipt.SourceDigest, receipt.DecisionIRDigest, receipt.DecisionDigest, receipt.AssessmentDigest, receipt.EvidenceDigest, receipt.CandidateDigest)
	for _, edge := range receipt.Edges {
		digest += fmt.Sprintf("|%s|%s|%s|%s|%s", edge.FromKind, edge.FromDigest, edge.ToKind, edge.ToDigest, edge.Relation)
	}
	return digestString(digest)
}
