package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"runtime"
	"strings"
	"sync"

	"github.com/afnank19/gale/requester"
)

// type Result struct {
// 	reqs int
// 	respSize int64 // Collective size of all responses in bytes
// 	latency []time.Duration
// 	statusCodes []int
// 	testDuration time.Duration
// }

// type ReqData struct {
// 	host string
// 	scheme string
// 	target string
// 	path string
// 	duration time.Duration
// }

func setMaxProcs(procs int) {
	runtime.GOMAXPROCS(procs)
}

func main() {
	// For a 4 core 8 thread cpu, the default by Go will be 8, but best perf comes from 4 so we default to that
	setMaxProcs(runtime.NumCPU() / 2)
	fmt.Println("GALE! Hammer your server, Threads:", runtime.NumCPU())
	// ParseArgs()

	argsState := ParseArgs()

	if argsState.threads != -1 {
		setMaxProcs(argsState.threads)
	}

	host, target := resolveUrlPort(argsState.url)
	fmt.Println("Host: "+host+" | Target: ", target+" | Path: ", argsState.url.Path)
	var rd requester.ReqData = requester.ReqData{
		Host: host,
		Scheme: argsState.url.Scheme,
		Target: target,
		Path: argsState.url.Path,
		Duration: argsState.duration,
	}
	var wg sync.WaitGroup
	var r requester.Result
	r.TestDuration = argsState.duration

	// range is connections
	for range argsState.connections {
		wg.Add(1)
		go r.MakeRequest(&wg, &rd)
	}

	wg.Wait()
	fmt.Println("Total Reqs: ", r.Reqs)
	fmt.Println("Resp size: ", r.RespSize / 1000, "KB")
	fmt.Println("Report:", requester.GenerateReport(r))
	// fmt.Println("Latency Slice :", e.latency)
}

func resolveUrlPort(url *url.URL) (host string, target string) {
	if url.Scheme == "http" {
		// handleHTTPUrl()
		host = url.Host

		// host already has a port specified
		if strings.Contains(host, ":") {
			target = host
		} else {
			target = host + ":80"
		}
	}

	if url.Scheme == "https" {
		// hanldeHTTPSUrl()
		host = url.Host
		target = host + ":443"
	}

	return
}

// func (r *Result) makeRequest(wg *sync.WaitGroup, rd *ReqData) {
// 	defer wg.Done()
// 	target := rd.target
// 	duration := rd.duration

// 	// Pre-build the raw HTTP/1.1 GET request bytes (keep-alive is default)
// 	reqLines := []string{
// 		"GET " + rd.path + " HTTP/1.1",
// 		fmt.Sprintf("Host: %s", rd.host),
// 		"Connection: keep-alive", // explicit, though default
// 		"Accept-Encoding: gzip, deflate, br",
// 		"",                       // end headers
// 		"",
// 	}
// 	rawReq := []byte(strings.Join(reqLines, "\r\n"))

// 	// Open one TCP connection
// 	conn, err := net.Dial("tcp", target)
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer conn.Close()

// 	var reader *bufio.Reader
// 	var writer *bufio.Writer

// 	if rd.scheme == "http" {
// 		reader = bufio.NewReader(conn)
// 		writer = bufio.NewWriter(conn)
// 	}

// 	if (rd.scheme == "https") {
// 		tlsConn := tls.Client(conn, &tls.Config{
// 			ServerName:         rd.host, // for SNI + cert validation
// 			InsecureSkipVerify: false,  // set true only for self-signed dev servers
// 			NextProtos: []string{"http/1.1"},
// 		})
// 		if err := tlsConn.Handshake(); err != nil {
// 			panic(err)
// 		}
// 		defer tlsConn.Close()

// 		reader = bufio.NewReader(tlsConn)
// 		writer = bufio.NewWriter(tlsConn)
// 	}

// 	var count uint64
// 	deadline := time.Now().Add(duration)

// 	for time.Now().Before(deadline) {
// 		// Start Req
// 		start := time.Now()
// 		// send request
// 		if _, err := writer.Write(rawReq); err != nil {
// 			break
// 		}
// 		if err := writer.Flush(); err != nil {
// 			break
// 		}

// 		// read response
// 		// we only care about consuming the headers and body (if any)
// 		// http.ReadResponse needs a *http.Request for context; use minimal stub
// 		resp, err := http.ReadResponse(reader, &http.Request{Method: "GET"})
// 		if err != nil {
// 			fmt.Println(err)
// 			break
// 		}
// 		// End Req
// 		end := time.Now()

// 		latency := end.Sub(start)
// 		// fmt.Println("TIme for req: ", elapsed)
// 		respSize := calculateResponseSize(resp)
// 		// consume and discard body
// 		if resp.ContentLength > 0 {
// 			io.CopyN(io.Discard, resp.Body, resp.ContentLength)
// 		}
// 		resp.Body.Close()

// 		atomic.AddUint64(&count, 1)
// 		r.reqs++
// 		r.respSize += respSize
// 		r.latency = append(r.latency, latency)
// 		r.statusCodes = append(r.statusCodes, resp.StatusCode)
// 	}

// 	fmt.Printf("Requests over one connection in %v: %d\n", duration, count)
// }

func calculateResponseSize(resp *http.Response) int64 {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}	
	// Adding to to body size accounting for \r\n
	bSize := int64(len(body)) + 2
	var hSize int = 0

	for key, values := range resp.Header {
		// Headers have colon ':' but the key variable doesnt so we add 1
		// There is a space between the key and value of headers, that is omitted here so we add 1, so a total of 2 bytes added
		hSize += len(key) + 2
        for _, value := range values {
			hSize += len(value)
    	}
		// adding 2 again to account for \r\n
		hSize += 2
	}

	// fmt.Println(hSize)
	// fmt.Println(bSize)

	// The final addition here is adding the Status Line bytes
	return bSize + int64(hSize) + int64(len(resp.Proto + " " + resp.Status + "\r\n"))
}