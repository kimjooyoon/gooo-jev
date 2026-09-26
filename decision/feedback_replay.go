package decision

import "fmt"

func ObserveFeedbackFromReplay(
	ledger Ledger,
	inputs []ReplayInput,
	index int,
	kind FeedbackKind,
	metricName string,
	metricValue float64,
	evidenceDigest string,
	recordedAt time.Time,
) (FeedbackObservation, error) {
	if err := ledger.Replay(inputs); err != nil {
		return FeedbackObservation{}, err
	}
	if index < 0 || index >= len(ledger.Entries) {
		return FeedbackObservation{}, fmt.Errorf("feedback replay index %d is outside ledger", index)
	}
	return ObserveFeedback(
		ledger,
		ledger.Entries[index].Receipt,
		kind,
		metricName,
		metricValue,
		evidenceDigest,
		recordedAt,
	)
}
