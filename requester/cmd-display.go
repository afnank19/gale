package requester

import (
	"fmt"
	"os"
	"runtime"
)

func DisplayReport(report *Report) {
	fmt.Println("Stats:")
	fmt.Println("  Total Requests:", report.TotalReqs)
	fmt.Println("  Total KB Transferred:", convertBytes(report.TotalBytes))
	fmt.Println("  Latency-> Avg:", report.AvgLatency,"| Std Dev:", report.StdDev, "| Max:", report.Max, "| Min:", report.Min)

	printLatencyPercentiles(report.Percentiles)
	printStatusCodeAnalysis(report.StatusCodes)

	fmt.Println("\nRequests/sec:", report.Rps)
	fmt.Println("Transfer/sec:", report.Tps, "MB")
}

func DisplayTestParameters(duration, url string, connections int ) {
	fmt.Printf("Running %s test on %s\n", duration, url)
	fmt.Printf("  %d threads with %d connections\n", runtime.GOMAXPROCS(0), connections)
}

func ShowUsage() {
	fmt.Printf("-USAGE-\n\n")
	fmt.Printf("Flag Structure: -[flagletter]=[value] OR --[flagname]=[value]\n")
	fmt.Printf("  --threads OR -t\t\tNumber of maximum threads to use. Default is No. of physical cores you have\n")
	fmt.Printf("  --connections OR -c\t\tNumber of concurrent connections. Ex: -c=10\n")
	fmt.Printf("  --duration OR -d\t\tTime to run the test. Ex: -d=10s (Units can be: s,m,h)\n")
	fmt.Printf("  --url OR -u\t\t\tThe url of the server. Ex: http://localhost:3000\n (REQUIRED)")
	os.Exit(1)
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

func printLatencyPercentiles(p Percentiles) {
	fmt.Println("\nLatency Distribution")
	fmt.Println("  50th:", p.P50)
	fmt.Println("  75th:", p.P75)
	fmt.Println("  90th:", p.P90)
	fmt.Println("  99th:", p.P99)
}