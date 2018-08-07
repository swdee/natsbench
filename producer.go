package main

import (
	"flag"
	"net/url"
	"runtime"
	"./pkg/producer"
)

func main() {

	msgSize := flag.Int("s", 32, "Size of message to send in bytes")
	msgRate := flag.Int("r", 1, "Number of messages per second to send")
	addr := flag.String("a", ":8081", "Address to listen on for /metrics endpoint")
	dest := flag.String("d", "", "Enter target endpoint to send messages to.  For NATS use nats://<ip>:<port> or direct connection to Worker worker://<ip>:<port>")
	mode := flag.String("m","shard", "Subject publishing mode on NATS.  shard or single")

	flag.Parse()

	// start http server with prometheus /metrics endpoint
	h := producer.NewHttpServer(*addr)
	go h.Start()

	// create a new producer that generates random strings for messages to send
	// according to the given msg/sec rate
	p := producer.NewProducer(*msgSize, *msgRate, *mode)

	// determine if we are sending to NATS or Worker and set the Producer to
	// connect to that receiver type
	u, err := url.Parse(*dest)

	if err != nil {
		panic("Error parsing Destation URL "+ err.Error() )
	}

	switch u.Scheme {
	case "nats":
		err = p.ConnectNats(*dest)
		if err != nil {
			panic("Error connecting to NATS "+ err.Error() )
		}

	case "worker":
		err = p.ConnectWorker(u.Host)
		if err != nil {
			panic("Error connecting to Worker "+ err.Error() )
		}

	default:
		panic("Unknown target destination")
	}


	runtime.Goexit()
}
