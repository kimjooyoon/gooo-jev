package decision

import (
  "errors"
  "math"
  "strings"
  "time"
)

type FeedbackWindowTrendDirection string

const (
  FeedbackWindowTrendRising FeedbackWindowTrendDirection = "rising"
  FeedbackWindowTrendFalling FeedbackWindowTrendDirection = "falling"
  FeedbackWindowTrendFlat FeedbackWindowTrendDirection = "flat"
)

type FeedbackWindowTrend struct {
  MetricName string
  PreviousEvidenceDigest string
  CurrentEvidenceDigest string
  PreviousKnownMean float64
  CurrentKnownMean float64
  KnownMeanDelta float64
  PreviousCoverage float64
  CurrentCoverage float64
  CoverageDelta float64
  PreviousUnknownCount int
  CurrentUnknownCount int
  Direction FeedbackWindowTrendDirection
  ComparedAt time.Time
  NonAuthorizing bool
  TrendDigest string
}

func CompareFeedbackWindows(previous, current FeedbackWindow, comparedAt time.Time) (FeedbackWindowTrend, error) {
  if err := previous.Validate(); err != nil {
    return FeedbackWindowTrend{}, err
  }
  if err := current.Validate(); err != nil {
    return FeedbackWindowTrend{}, err
  }
  if previous.MetricName != current.MetricName {
    return FeedbackWindowTrend{}, errors.New("feedback window trend metric names do not match")
  }
  if comparedAt.IsZero() {
    return FeedbackWindowTrend{}, errors.New("feedback window trend comparison time is required")
  }
  meanDelta := current.KnownMean - previous.KnownMean
  direction := FeedbackWindowTrendFlat
  switch {
  case meanDelta > 0:
    direction = FeedbackWindowTrendRising
  case meanDelta < 0:
    direction = FeedbackWindowTrendFalling
  }
  trend := FeedbackWindowTrend{
    MetricName: previous.MetricName,
    PreviousEvidenceDigest: previous.EvidenceDigest,
    CurrentEvidenceDigest: current.EvidenceDigest,
    PreviousKnownMean: previous.KnownMean,
    CurrentKnownMean: current.KnownMean,
    KnownMeanDelta: meanDelta,
    PreviousCoverage: previous.ObservedCoverage,
    CurrentCoverage: current.ObservedCoverage,
    CoverageDelta: current.ObservedCoverage - previous.ObservedCoverage,
    PreviousUnknownCount: previous.UnknownCount,
    CurrentUnknownCount: current.UnknownCount,
    Direction: direction,
    ComparedAt: comparedAt.UTC(),
    NonAuthorizing: true,
  }
  trend.TrendDigest, _ = Digest(trend)
  return trend, nil
}

func (trend FeedbackWindowTrend) Validate() error {
  if strings.TrimSpace(trend.MetricName) == "" ||
    strings.TrimSpace(trend.PreviousEvidenceDigest) == "" ||
    strings.TrimSpace(trend.CurrentEvidenceDigest) == "" ||
    strings.TrimSpace(trend.TrendDigest) == "" {
    return errors.New("feedback window trend is incomplete")
  }
  values := []float64{trend.PreviousKnownMean, trend.CurrentKnownMean, trend.KnownMeanDelta, trend.PreviousCoverage, trend.CurrentCoverage, trend.CoverageDelta}
  for _, value := range values {
    if math.IsNaN(value) || math.IsInf(value, 0) {
      return errors.New("feedback window trend metrics must be finite")
    }
  }
  if trend.ComparedAt.IsZero() || !trend.NonAuthorizing {
    return errors.New("feedback window trend boundary is invalid")
  }
  if trend.PreviousUnknownCount < 0 || trend.CurrentUnknownCount < 0 ||
    trend.PreviousCoverage < 0 || trend.PreviousCoverage > 1 ||
    trend.CurrentCoverage < 0 || trend.CurrentCoverage > 1 {
    return errors.New("feedback window trend coverage is invalid")
  }
  if trend.KnownMeanDelta != trend.CurrentKnownMean-trend.PreviousKnownMean {
    return errors.New("feedback window trend mean delta is invalid")
  }
  if trend.CoverageDelta != trend.CurrentCoverage-trend.PreviousCoverage {
    return errors.New("feedback window trend coverage delta is invalid")
  }
  expectedDirection := FeedbackWindowTrendFlat
  switch {
  case trend.KnownMeanDelta > 0:
    expectedDirection = FeedbackWindowTrendRising
  case trend.KnownMeanDelta < 0:
    expectedDirection = FeedbackWindowTrendFalling
  }
  if trend.Direction != expectedDirection {
    return errors.New("feedback window trend direction is invalid")
  }
  copy := trend
  copy.TrendDigest = ""
  digest, err := Digest(copy)
  if err != nil {
    return err
  }
  if digest != trend.TrendDigest {
    return errors.New("feedback window trend digest does not match its evidence")
  }
  return nil
}