package main

import "fmt"
import "log"
import "net"
import "os"
import "rand"
import "time"

// Delay the response by only allowing for 1 byte every [min,max)
// nanoseconds.
func TrickleResponse(host, url string, min, max uint, counter *Counter) {
	conn, err := net.Dial("tcp", "", host)
	if err != nil {
		log.Printf("Failed to connect to web server: %s", err)
		return
	}

	counter.Increment()
	fmt.Fprintf(conn, "GET /%s HTTP/1.1\n\n", url)

	// Trickle receive the response one byte at a time
	var buf []byte = make([]byte, 1, 1)

	for {
		n, err := conn.Read(buf)
		if n == 0 && err == os.EOF {
			break
		}
		if err != nil {
			log.Printf("Got an error when reading: %s", err)
			break
		}

		delay := int(min) - rand.Intn(int(max - min))
		time.Sleep(int64(delay))
	}

	counter.Decrement()
	if conn != nil {
		conn.Close()
	}
}

// Delay the response by only allowing for 1 byte every [min,max)
// nanoseconds.
func TrickleRequest(host, url string, payload, min, max uint, counter *Counter) {
	conn, err := net.Dial("tcp", "", host)
	if err != nil {
		log.Printf("Failed to connect to web server: %s", err)
		return
	}

	counter.Increment()

	stupid := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	request := fmt.Sprintf("GET /%s HTTP/1.1\n", url)
	for len(request) < int(payload) {
		request = fmt.Sprintf("%sFakeData%d: \"%s\"\n", request, len(request), stupid)
	}

	bytes := []byte(request)
	for idx, _ := range bytes {
		_, err := conn.Write(bytes[idx:idx+1])
		if err != nil {
			break
		}
		delay := int(min) - rand.Intn(int(max - min))
		time.Sleep(int64(delay))
	}

	// Trickle receive the response one byte at a time
	size_32k := 32 * 1024
	var buf []byte = make([]byte, 0, size_32k)

	for {
		n, err := conn.Read(buf)
		if n == 0 && err == os.EOF {
			break
		}
		if err != nil {
			log.Printf("Got an error when reading response: %s", err)
			break
		}

	}

	counter.Decrement()
	if conn != nil {
		conn.Close()
	}
}

