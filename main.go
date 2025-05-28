package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type State struct {
	workers     int
	concurrency int
	duration    time.Duration
}

type Example struct {
	reqs int
}

func makeHTTPSRequest(e *Example, wg *sync.WaitGroup) {
	defer wg.Done()
	const target = "smiling-kerstin-afnan-we-2f62af17.koyeb.app:443"
	const duration = 5 * time.Second

	// Pre-build the raw HTTP/1.1 GET request bytes (keep-alive is default)
	reqLines := []string{
		"GET /api/blogs/utnmGBLv2oIOquzyXQxu HTTP/1.1",
		"Host: smiling-kerstin-afnan-we-2f62af17.koyeb.app",
		"Connection: keep-alive", // explicit, though default
		"",                       // end headers
		"",
	}
	rawReq := []byte(strings.Join(reqLines, "\r\n"))

	// Open one TCP connection
	conn, err := net.Dial("tcp", target)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	tlsConn := tls.Client(conn, &tls.Config{
		ServerName:         "smiling-kerstin-afnan-we-2f62af17.koyeb.app", // for SNI + cert validation
		InsecureSkipVerify: false,                                         // set true only for self-signed dev servers
	})
	if err := tlsConn.Handshake(); err != nil {
		panic(err)
	}
	defer tlsConn.Close()

	reader := bufio.NewReader(tlsConn)
	writer := bufio.NewWriter(tlsConn)

	var count uint64
	deadline := time.Now().Add(duration)

	for time.Now().Before(deadline) {
		// send request
		if _, err := writer.Write(rawReq); err != nil {
			break
		}
		if err := writer.Flush(); err != nil {
			break
		}

		// read response
		// we only care about consuming the headers and body (if any)
		// http.ReadResponse needs a *http.Request for context; use minimal stub
		resp, err := http.ReadResponse(reader, &http.Request{Method: "GET"})
		if err != nil {
			fmt.Println(err)
			break
		}
		// fmt.Println(resp.Status)
		// consume and discard body
		if resp.ContentLength > 0 {
			io.CopyN(io.Discard, resp.Body, resp.ContentLength)
		}
		resp.Body.Close()

		atomic.AddUint64(&count, 1)
		e.reqs++
	}

	fmt.Printf("Requests over one connection in %v: %d\n", duration, count)
}

func main() {
	fmt.Println("GALE! Hammer your server")
	// ParseArgs()

	var wg sync.WaitGroup
	var e Example

	// range is connections
	for range 100 {
		wg.Add(1)
		// go makeHTTPSRequest(&e, &wg)
		go makeHTTPRequest(&e, &wg)
	}

	wg.Wait()
	fmt.Println("Total Reqs: ", e.reqs)
}

func makeHTTPRequest(e *Example, wg *sync.WaitGroup) {
	defer wg.Done()
	const target = "localhost:8080"
	const duration = 30 * time.Second

	// Pre-build the raw HTTP/1.1 GET request bytes (keep-alive is default)
	reqLines := []string{
		"GET /api/blogs/utnmGBLv2oIOquzyXQxu HTTP/1.1",
		"Host: localhost:8080",
		"Connection: keep-alive", // explicit, though default
		"",                       // end headers
		"",
	}
	rawReq := []byte(strings.Join(reqLines, "\r\n"))

	// Open one TCP connection
	conn, err := net.Dial("tcp", target)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	var count uint64
	deadline := time.Now().Add(duration)

	for time.Now().Before(deadline) {
		// send request
		if _, err := writer.Write(rawReq); err != nil {
			break
		}
		if err := writer.Flush(); err != nil {
			break
		}

		// read response
		// we only care about consuming the headers and body (if any)
		// http.ReadResponse needs a *http.Request for context; use minimal stub
		resp, err := http.ReadResponse(reader, &http.Request{Method: "GET"})
		if err != nil {
			fmt.Println(err)
			break
		}
		// fmt.Println(resp.Status)
		// consume and discard body
		if resp.ContentLength > 0 {
			io.CopyN(io.Discard, resp.Body, resp.ContentLength)
		}
		resp.Body.Close()

		atomic.AddUint64(&count, 1)
		e.reqs++
	}

	fmt.Printf("Requests over one connection in %v: %d\n", duration, count)
}
