package decision

import "testing"

func TestMetricDeltaDirectionsAndUnknownState(t *testing.T) {
    cases := []struct {
        name      string
        before    float64
        after     float64
        status    MetricDeltaStatus
        direction MetricDeltaDirection
    }{
        {name: "improved", before: 1, after: 2, status: MetricObserved, direction: MetricImproved},
        {name: "regressed", before: 2, after: 1, status: MetricObserved, direction: MetricRegressed},
        {name: "unchanged", before: 1, after: 1, status: MetricObserved, direction: MetricUnchanged},
        {name: "unknown", before: 0, after: 0, status: MetricUnknown, direction: MetricDirectionUnknown},
    }

    for _, tc := range cases {
        t.Run(tc.name, func(t *testing.T) {
            delta, err := NewMetricDelta(
                "latency",
                tc.before,
                tc.after,
                "evidence-"+tc.name,
                tc.status,
            )
            if err != nil {
                t.Fatalf("NewMetricDelta() error = %v", err)
            }
            if delta.Direction != tc.direction {
                t.Fatalf("direction = %q, want %q", delta.Direction, tc.direction)
            }
            if err := delta.Validate(); err != nil {
                t.Fatalf("Validate() error = %v", err)
            }
        })
    }

    tampered, err := NewMetricDelta("quality", 0.5, 0.75, "evidence-quality", MetricObserved)
    if err != nil {
        t.Fatalf("tampered fixture construction error = %v", err)
    }
    tampered.After = 0.25
    if err := tampered.Validate(); err == nil {
        t.Fatal("tampered metric delta unexpectedly validated")
    }

    if _, err := NewMetricDelta("quality", 0, 0, "evidence-invalid", MetricObserved); err != nil {
        t.Fatalf("zero-valued observed metric should be valid: %v", err)
    }
}
