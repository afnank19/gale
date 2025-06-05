package requester

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

type Result struct {
	Reqs int // total number of requests made
	RespSize int64 // Collective size of all responses in bytes
	Latency []time.Duration
	StatusCodes map[int]int
	TestDuration time.Duration

	mu sync.Mutex
}

type ReqData struct {
	Host string
	Scheme string
	Target string
	Path string
	Duration time.Duration
}

func (r *Result) MakeRequest(wg *sync.WaitGroup, rd *ReqData) {
	defer wg.Done()
	target := rd.Target
	duration := rd.Duration

	// Pre-build the raw HTTP/1.1 GET request bytes (keep-alive is default)
	reqLines := []string{
		"GET " + rd.Path + " HTTP/1.1",
		fmt.Sprintf("Host: %s", rd.Host),
		"Connection: keep-alive", // explicit, though default
		"Accept-Encoding: gzip, deflate, br",
		"",                       // end headers
		"",
	}
	rawReq := []byte(strings.Join(reqLines, "\r\n"))

	// Open one TCP connection
	conn, err := net.Dial("tcp", target)
	if err != nil {
		log.Fatalln("Could not connect:", err)
	}
	defer conn.Close()

	var reader *bufio.Reader
	var writer *bufio.Writer

	if rd.Scheme == "http" {
		reader = bufio.NewReader(conn)
		writer = bufio.NewWriter(conn)
	}

	if (rd.Scheme == "https") {
		tlsConn := tls.Client(conn, &tls.Config{
			ServerName:         rd.Host, // for SNI + cert validation
			InsecureSkipVerify: false,  // set true only for self-signed dev servers
			NextProtos: []string{"http/1.1"},
		})
		if err := tlsConn.Handshake(); err != nil {
			fmt.Println(err)
			return
		}
		defer tlsConn.Close()

		reader = bufio.NewReader(tlsConn)
		writer = bufio.NewWriter(tlsConn)
	}

	deadline := time.Now().Add(duration)

	for time.Now().Before(deadline) {
		// Start Req
		// send request
		if _, err := writer.Write(rawReq); err != nil {
			break
		}
		start := time.Now()
		if err := writer.Flush(); err != nil {
			break
		}

		// read response
		// we only care about consuming the headers and body (if any)
		// http.ReadResponse needs a *http.Request for context; use minimal stub
		resp, err := http.ReadResponse(reader, &http.Request{Method: "GET"})
		end := time.Now()
		latency := end.Sub(start)

		if err != nil {
			fmt.Println(err)
			break
		}
		// End Req
		respSize := calculateResponseSize(resp)
		// consume and discard body
		if resp.ContentLength > 0 {
			io.CopyN(io.Discard, resp.Body, resp.ContentLength)
		}
		resp.Body.Close()

		r.mu.Lock()
		r.Reqs++
		r.RespSize += respSize
		r.Latency = append(r.Latency, latency)
		r.StatusCodes[resp.StatusCode]++
		r.mu.Unlock()
	}
}

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