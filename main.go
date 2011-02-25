package main

import "flag"
import "fmt"
import "log"
import "net"
import "os"
import "time"

type Counter struct {
	value int
	delta chan int
	result chan int
	status chan int
	report func()
}

func NewCounter(tick int64) (*Counter, chan int) {
	status := make(chan int)
	counter := &Counter{
		value: 0,
		delta: make(chan int),
		result: make(chan int),
		status: status,
	}
	counter.report = func() {
		counter.delta <- 0
		counter.status <- <-counter.result
		time.AfterFunc(tick, counter.report)
	}

	time.AfterFunc(tick, counter.report)
	go counter.serve()
	return counter, counter.status
}

func (c *Counter) Stop() {
	close(c.delta)
	close(c.result)
	close(c.status)
}

func (c *Counter) serve() {
	for delta := range c.delta {
		c.value += delta
		c.result <- c.value
	}
}

func (c *Counter) Get() int {
	c.delta <- 0
	return <-c.result
}

func (c *Counter) Increment() int {
	c.delta <- 1
	return <-c.result
}

func (c *Counter) Decrement() int {
	c.delta <- -1
	return <-c.result
}

func TrickleClient(host, url string, delay int64, counter *Counter) {
	conn, err := net.Dial("tcp", "", host)
	if err != nil {
		log.Printf("Failed to connect to web server: %s", err)
		time.Sleep(delay)
	} else {
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
				time.Sleep(delay)
			}

		counter.Decrement()
	}

	if conn != nil {
		conn.Close()
	}

	TrickleClient(host, url, delay, counter)
}

func main() {
	var host *string = flag.String("host", "localhost:12345", "The host to open a connection with")
	var url *string = flag.String("url", "/", "The URL to request")
	var num *int = flag.Int("num", 5, "The number of slow clients to open")
	var delay *float64 = flag.Float64("delay", 1.0, "The delay in receiving each byte (in seconds)")

	flag.Parse()

	counter, status := NewCounter(1e9 * 2)
	delaynano := int64(*delay * 1e9)

	for i := 0; i < *num; i++ {
		go TrickleClient(*host, *url, delaynano, counter)
	}

	for val := range status {
		log.Printf("There are %d clients open", val)
	}
}
