package decision

import "strings"

// JEVActionCandidateVerificationMetricInput contains exact disposition counts
// and the digest of the evidence source that produced them.
type JEVActionCandidateVerificationMetricInput struct {
	MetricName           string
	VerifiedCount        uint64
	ReviewCount          uint64
	CounterexampleCount  uint64
	UnknownCount         uint64
	EvidenceSourceDigest string
	NonAuthorizing       bool
}

// JEVActionCandidateVerificationMetric records counts without inferring
// improvement, regression, or execution success.
type JEVActionCandidateVerificationMetric struct {
	Status                string
	MetricName            string
	Total                 uint64
	VerifiedCount         uint64
	ReviewCount           uint64
	CounterexampleCount   uint64
	UnknownCount          uint64
	EvidenceSourceDigest  string
	EvidenceDigest        string
	MissingStage          string
	NonExecuting          bool
	NonAuthorizing        bool
}

// MeasureJEVActionCandidateVerificationMetric preserves every disposition and
// requires an evidence source before exposing a measured result.
func MeasureJEVActionCandidateVerificationMetric(input JEVActionCandidateVerificationMetricInput) JEVActionCandidateVerificationMetric {
	output := JEVActionCandidateVerificationMetric{
		Status:         "UNKNOWN",
		NonExecuting:   true,
		NonAuthorizing: true,
	}
	if !input.NonAuthorizing {
		output.NonAuthorizing = false
		output.MissingStage = "authorization-boundary"
		return output
	}
	if strings.TrimSpace(input.MetricName) == "" {
		output.MissingStage = "metric-name"
		return output
	}
	if strings.TrimSpace(input.EvidenceSourceDigest) == "" {
		output.MissingStage = "evidence-source"
		return output
	}
	counts := []uint64{
		input.VerifiedCount,
		input.ReviewCount,
		input.CounterexampleCount,
		input.UnknownCount,
	}
	var total uint64
	for _, count := range counts {
		if ^uint64(0)-total < count {
			output.MissingStage = "count-overflow"
			return output
		}
		total += count
	}
	if total == 0 {
		output.MissingStage = "observations"
		return output
	}
	evidenceDigest, err := Digest(struct {
		MetricName           string
		VerifiedCount        uint64
		ReviewCount          uint64
		CounterexampleCount  uint64
		UnknownCount         uint64
		EvidenceSourceDigest string
	}{
		MetricName:           input.MetricName,
		VerifiedCount:        input.VerifiedCount,
		ReviewCount:          input.ReviewCount,
		CounterexampleCount:  input.CounterexampleCount,
		UnknownCount:         input.UnknownCount,
		EvidenceSourceDigest: input.EvidenceSourceDigest,
	})
	if err != nil {
		output.MissingStage = "metric-evidence"
		return output
	}
	output.MetricName = input.MetricName
	output.Total = total
	output.VerifiedCount = input.VerifiedCount
	output.ReviewCount = input.ReviewCount
	output.CounterexampleCount = input.CounterexampleCount
	output.UnknownCount = input.UnknownCount
	output.EvidenceSourceDigest = input.EvidenceSourceDigest
	output.EvidenceDigest = evidenceDigest
	output.Status = "measured"
	if output.UnknownCount > 0 {
		output.Status = "measured-with-unknown"
	} else if output.CounterexampleCount > 0 {
		output.Status = "measured-with-counterexample"
	}
	return output
}
