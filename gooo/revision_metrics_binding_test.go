package gooo

import "testing"

func metricsBindingInputs(t *testing.T) (RevisionApplicationReceipt, RevisionMetrics) {
	t.Helper()
	plan, application := appliedRevisionReceiptInputs(t)
	receipt, err := ObserveRevisionApplicationReceipt(plan, application)
	if err != nil {
		t.Fatalf("ObserveRevisionApplicationReceipt() error = %v", err)
	}
	metrics, err := MeasureRevision(validContract, application)
	if err != nil {
		t.Fatalf("MeasureRevision() error = %v", err)
	}
	return receipt, metrics
}

func TestObserveRevisionMetricsBindingBindsReceiptAndMetrics(t *testing.T) {
	receipt, metrics := metricsBindingInputs(t)
	binding, err := ObserveRevisionMetricsBinding(receipt, metrics)
	if err != nil {
		t.Fatalf("ObserveRevisionMetricsBinding() error = %v", err)
	}
	if binding.Status != "BOUND" ||
		binding.ApplicationDigest != receipt.ApplicationDigest ||
		binding.MetricsDigest != metrics.MetricsDigest ||
		binding.ChangedByteCount != metrics.ChangedByteCount ||
		binding.IRChanged != metrics.IRChanged {
		t.Fatalf("unexpected metrics binding: %#v", binding)
	}
	if err := binding.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
}

func TestObserveRevisionMetricsBindingRetainsReceiptFailure(t *testing.T) {
	receipt, metrics := metricsBindingInputs(t)
	receipt.ReceiptDigest = digestString("tampered")
	binding, err := ObserveRevisionMetricsBinding(receipt, metrics)
	if err == nil {
		t.Fatal("ObserveRevisionMetricsBinding() error = nil, want receipt failure")
	}
	if binding.Status != "UNKNOWN" || binding.MissingStage != "revision-metrics-binding-receipt" {
		t.Fatalf("unexpected unknown metrics binding: %#v", binding)
	}
}

func TestObserveRevisionMetricsBindingRetainsLinkFailure(t *testing.T) {
	receipt, metrics := metricsBindingInputs(t)
	metrics.ApplicationDigest = digestString("other-application")
	metrics.MetricsDigest = digestRevisionMetrics(metrics)
	binding, err := ObserveRevisionMetricsBinding(receipt, metrics)
	if err == nil {
		t.Fatal("ObserveRevisionMetricsBinding() error = nil, want link failure")
	}
	if binding.Status != "UNKNOWN" || binding.MissingStage != "revision-metrics-binding-link" {
		t.Fatalf("unexpected unknown metrics binding: %#v", binding)
	}
}

func TestObserveRevisionMetricsBindingIsDeterministic(t *testing.T) {
	receipt, metrics := metricsBindingInputs(t)
	first, err := ObserveRevisionMetricsBinding(receipt, metrics)
	if err != nil {
		t.Fatalf("first ObserveRevisionMetricsBinding() error = %v", err)
	}
	second, err := ObserveRevisionMetricsBinding(receipt, metrics)
	if err != nil {
		t.Fatalf("second ObserveRevisionMetricsBinding() error = %v", err)
	}
	if first.BindingDigest != second.BindingDigest {
		t.Fatal("same receipt and metrics produced different binding digest")
	}
}
