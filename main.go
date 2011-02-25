package main

import "flag"
import "log"
import "rand"
import "time"

// Generate a stream of 'num' bad clients against a given host and url. The
// ratio of different types of clients is random, but you can specify the
// options of the different types of clients on the commandline.

var host *string = flag.String("host", "localhost:12345", "The host and port to contact")
var url *string = flag.String("url", "/", "The URL to request")
var num *int = flag.Int("num", 500, "The number of bad clients to sustain")
var reqmin *float64 = flag.Float64("reqmin", 1.0, "Minimum delay time for trickle request (in seconds)")
var reqmax *float64 = flag.Float64("reqmax", 5.0, "Maximum delay time for trickle request (in seconds)")
var reqlen *uint = flag.Uint("reqlen", 1024, "Payload length for request data")
var respmin *float64 = flag.Float64("respmin", 1.0, "Minimum delay time for trickle response (in seconds)")
var respmax *float64 = flag.Float64("respmax", 5.0, "Maximum delay time for trickle response (in seconds)")
var delay *float64 = flag.Float64("delay", 5.0, "Delay time for status report")

func main() {
	flag.Parse()

	counter := NewCounter(0)

	for {
		numClients := counter.Get()
		log.Printf("%d clients are currently active", numClients)

		needed := *num - numClients

		if needed > 0 {
			tresp := 0
			treq := 0
			for i := needed; i > 0; i-- {
				ctype := rand.Intn(2)
				switch ctype {
				case 0: // trickle response
					min := uint(*respmin * 1e9)
					max := uint(*respmax * 1e9)
					go TrickleResponse(*host, *url, min, max, counter)
					tresp++

				case 1: // trickle request
					min := uint(*reqmin * 1e9)
					max := uint(*reqmax * 1e9)
					go TrickleRequest(*host, *url, *reqlen, min, max, counter)
					treq++

				default:
					panic("this should not happen")
				}
			}

			log.Printf("Opened %d new clients (%d trickle_req, %d trickle_resp)", needed, treq, tresp)
		}

		// Sleep for a bit, then do everything again!
		time.Sleep(int64(*delay * 1e9))
	}
}
