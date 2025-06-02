// Package to user results from request results
package requester

import (
	"time"
)

type Report struct {
	AvgLatency int
	Rps float64 // requests per second
	Tps float64 // transfer per second
}

func GenerateReport(reqResult Result) *Report{
	var report Report
	report.CalculateRPS(float64(reqResult.Reqs), reqResult.TestDuration)
	report.CalculateTPS(float64(reqResult.RespSize), reqResult.TestDuration)

	return &report
}

func (r *Report) CalculateRPS(totalReqs float64, duration time.Duration) {
	r.Rps = totalReqs / duration.Seconds()
}

func (r *Report) CalculateTPS(totalBytes float64, duration time.Duration) {
	r.Tps = totalBytes / duration.Seconds()
}

func (r *Report) CalculateAvgLatency() {

}