package requester

import "fmt"

func DisplayReport(report *Report) {
	fmt.Println("Stats:")
	fmt.Println("  Total Requests:", report.TotalReqs)
	fmt.Println("  Total KB Transferred:", convertBytes(report.TotalBytes))
	fmt.Println("  Latency-> Avg:", report.AvgLatency,"| Std Dev:", report.StdDev, "| Max:", report.Max, "| Min:", report.Min)

	printStatusCodeAnalysis(report.StatusCodes)
	
	fmt.Println("\nRequests/sec:", report.Rps)
	fmt.Println("Transfer/sec:", report.Tps, "MB")
}

func convertBytes(bytes int64) (string) {
	const KbThreshold int64 = 1000
	const MbThreshold int64 = 1_000_000

	var data int64
	var unit string
	
	if bytes > KbThreshold {
		data = bytes / 1000
		unit = "KB"
	}

	if bytes > MbThreshold {
		data = bytes / 1_000_000
		unit = "MB"
	} 

	return fmt.Sprintf("%d%s", data, unit)
}

func printStatusCodeAnalysis(statusCodes map[int]int) {
	var non200Count int = 0

	fmt.Println("\nStatus Codes [Code: Count]")
	for key, value := range statusCodes {
		fmt.Println(" ",key,":",value)
		if key > 399 {
			non200Count += value
		}
	}

	if non200Count > 0 {
		fmt.Println("Non 2XX or 3XX response count:", non200Count)
	}
}