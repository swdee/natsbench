package main

import (
	"./pkg/worker"
	"flag"
	"fmt"
	"runtime"
)

func main() {

	addr := flag.String("a", ":8082", "Address to listen on for /metrics endpoint")
	serv := flag.String("s", "", "Set Worker to recieve messages directly from Producer.  Enter the Address the Worker should listen on")
	nats := flag.String("n", "", "Set Worker to recieve messages from NATS. Enter address of NATS server in format nats://<ip>:<port>")
	mode := flag.String("m","shard", "Subject publishing mode on NATS.  shard or single")

	flag.Parse()

	if *serv == "" && *nats == "" {
		panic("Worker needs to be started to receive msgs directly or from NATS")
	}

	// start http server with prometheus /metrics endpoint
	h := worker.NewHttpServer(*addr)
	go h.Start()

	if *nats != "" {
		runNats(*nats, *mode)
	} else if *serv != "" {
		runServer(*serv)
	}

	runtime.Goexit()
}

func runNats(url string, mode string) {

	// create new worker of given type to receive messages sent by producer
	w := worker.NewWorker()

	w.ConnectNats(url)

	if mode == "shard" {
		subjs := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "a", "b", "c", "d", "e", "f"}

		for _, sub := range subjs {
			err := w.Subscribe(fmt.Sprintf("txn.%s.>", sub))

			if err != nil {
				panic("Error subscribing to NATS " + err.Error())
			}
		}
	} else {
		// single mode
		err := w.Subscribe("single")

		if err != nil {
			panic("Error subscribing to NATS " + err.Error())
		}
	}

	fmt.Println("Awaiting messages from NATS")
}

func runServer(addr string) {

	w := worker.NewWorker()
	err := w.StartServer(addr)

	if err != nil {
		panic("Unable to bind Worker Server to address "+err.Error() )
	}

	fmt.Printf("Worker is listening on %s, now awaiting messages from Producer\n", addr)
}
