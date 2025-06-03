// Package to user results from request results
package requester

import (
	"math"
	"slices"
	"time"
)

type Report struct {
	// Latency Fields
	AvgLatency time.Duration
	Max time.Duration // maximum latency 
	Min time.Duration // minimum latency
	StdDev time.Duration

	Rps float64 // requests per second
	Tps float64 // transfer per second in MB

	TotalBytes int64
	TotalReqs int

	StatusCodes map[int]int
}

func GenerateReport(reqResult *Result) *Report{
	var report Report
	report.CalculateRPS(float64(reqResult.Reqs), reqResult.TestDuration)
	// Calculating transfer per second in MB, where MB is decimal (1,000,000)
	report.CalculateTPS(float64(reqResult.RespSize) / 1_000_000, reqResult.TestDuration)
	report.CalculateAvgLatency(reqResult.Latency, reqResult.Reqs)
	report.CalculateStandardDeviation(reqResult.Latency, reqResult.Reqs)

	// Maximum and minimum latency
	report.Max = slices.Max(reqResult.Latency)
	report.Min = slices.Min(reqResult.Latency)

	report.TotalReqs = reqResult.Reqs
	report.TotalBytes = reqResult.RespSize
	
	report.StatusCodes = reqResult.StatusCodes

	return &report
}

func (r *Report) CalculateRPS(totalReqs float64, duration time.Duration) {
	r.Rps = totalReqs / duration.Seconds()
}

func (r *Report) CalculateTPS(totalBytes float64, duration time.Duration) {
	r.Tps = totalBytes / duration.Seconds()
}

func (r *Report) CalculateAvgLatency(latencies []time.Duration, totalReqs int) {
	var totalLatency float64

	for _, l := range latencies {
		totalLatency += float64(l)
	}

	r.AvgLatency = time.Duration(totalLatency / float64(len(latencies)))
}

func (r *Report) CalculateStandardDeviation(latencies []time.Duration, totalReqs int) {
	var deviationSum float64
	for _, l := range latencies {
		deviation := (float64(l) - float64(r.AvgLatency))
		deviationSum += deviation * deviation
	}

	variance := deviationSum / float64(totalReqs)
	stdDev := math.Sqrt(variance)
	r.StdDev = time.Duration(stdDev)
}