package worker

import (
	"fmt"
	"net/http"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var msgsRecv = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "natsbench_worker_msgs_sent_total",
		Help: "Total number of messages received by Worker",
	},
)

type HttpServer struct {
	addr string
	http *http.ServeMux
}

func init() {
	prometheus.MustRegister(msgsRecv)
}

func NewHttpServer(a string) *HttpServer {
	return &HttpServer{
		addr: a,
	}
}

func (s *HttpServer) Start() {

	s.http = http.NewServeMux()

	s.http.Handle("/metrics", promhttp.Handler())

	fmt.Println("HTTP server listening on", s.addr)

	err := http.ListenAndServe(s.addr, s.http)

	if err != nil {
		panic("Error starting web server: " + err.Error())
	}
}
