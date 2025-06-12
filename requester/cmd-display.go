package requester

import (
	"fmt"
	"os"
	"runtime"

	"github.com/charmbracelet/lipgloss"
)

/* STYLES */
var title = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.AdaptiveColor{Light: "#3c3836", Dark: "#fbf1c7"})
var titleNB = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "#3c3836", Dark: "#fbf1c7"}).Italic(true) // No bold
var text = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "#3c3836", Dark: "#fbf1c7"})
var red = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "#9d0006", Dark: "#fb4934"})
var yellow = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "#b57614", Dark: "#fabd2f"})
var green = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "#79740e", Dark: "#b8bb26"})
var blue = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "#076678", Dark: "#83a598"})

/* STYLES */

func DisplayReport(report *Report) {
	fmt.Println(title.Render("Stats:"))
	output := fmt.Sprintf(
		"  Total Requests: %d\n  Total Data Read: %s\n",
		report.TotalReqs,
		convertBytes(report.TotalBytes),
	)

	avg := blue.Render(report.AvgLatency.String())
	stddev := blue.Render(report.StdDev.String())
	max := yellow.Render(report.Max.String())
	min := green.Render(report.Min.String())

	output += fmt.Sprintf("  Latency → Avg: %s | Std Dev: %s | Max: %s | Min: %s",
		avg,
		stddev,
		max,
		min,
	)
	fmt.Println(text.Render(output))

	printLatencyPercentiles(report.Percentiles)
	printStatusCodeAnalysis(report.StatusCodes)

	output = fmt.Sprintf(
		"\nRequests/sec: %v\n"+
			"Transfer/sec: %v MB",
		report.Rps,
		report.Tps,
	)
	fmt.Println(text.Render(output))
}

func DisplayTestParameters(duration, url string, connections int) {
	output := fmt.Sprintf(
		"Running %s test on %s\n"+
			"  %d threads with %d connections\n",
		duration,
		url,
		runtime.GOMAXPROCS(0),
		connections,
	)
	fmt.Println(text.Render(output))
}

func ShowUsage() {
	output := fmt.Sprintf(
		"-USAGE-\n" +
			"Flag Structure: -[flagletter]=[value] OR --[flagname]=[value]\n" +
			"  --threads OR -t\t\tNumber of maximum threads to use. Default is No. of physical cores you have\n" +
			"  --connections OR -c\t\tNumber of concurrent connections. Ex: -c=10\n" +
			"  --duration OR -d\t\tTime to run the test. Ex: -d=10s (Units can be: s,m,h)\n" +
			"  --url OR -u\t\t\tThe url of the server. Ex: http://localhost:3000 (REQUIRED)",
	)
	fmt.Println(text.Render(output))

	os.Exit(1)
}

func convertBytes(bytes int64) string {
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

	fmt.Println(title.Render("\nStatus Codes [Code: Count]"))
	for key, value := range statusCodes {
		if key < 400 {
			fmt.Println(green.Render(fmt.Sprintf("  %d : %d", key, value)))
		} else if key < 500 {
			fmt.Println(yellow.Render(fmt.Sprintf("  %d : %d", key, value)))
		} else {
			fmt.Println(red.Render(fmt.Sprintf("  %d : %d", key, value)))
		}

		if key > 399 {
			non200Count += value
		}
	}

	if non200Count > 0 {
		fmt.Println(titleNB.Render(fmt.Sprintf("\nNon 2XX or 3XX response count: %d", non200Count)))
	}
}

func printLatencyPercentiles(p Percentiles) {
	fmt.Println(title.Render("\nLatency Distribution"))
	fmt.Println(green.Render("  50th:", p.P50.String()))
	fmt.Println(yellow.Render("  75th:", p.P75.String()))
	fmt.Println(yellow.Render("  90th:", p.P90.String()))
	fmt.Println(red.Render("  99th:", p.P99.String()))
}
