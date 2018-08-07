package worker

import (
	"bufio"
	"github.com/nats-io/go-nats"
	"net"
)

type Worker struct {
	ln       net.Listener
	conn     *nats.Conn
	isWorker bool
}

var err error

func NewWorker() *Worker {
	return &Worker{}
}

func (w *Worker) ConnectNats(url string) error {

	w.conn, err = nats.Connect(url)

	if err != nil {
		return err
	}

	w.isWorker = false

	return nil
}

func (w *Worker) Subscribe(subj string) error {
	w.conn.Subscribe(subj, w.processMsg)
	w.conn.Flush()

	if err := w.conn.LastError(); err != nil {
		return err
	}

	return nil
}

func (w *Worker) processMsg(m *nats.Msg) {
	// increment prometheus counter for total msgs recieved
	msgsRecv.Inc()
}

func (w *Worker) StartServer(addr string) error {

	w.ln, err = net.Listen("tcp", addr)

	if err != nil {
		return err
	}

	w.isWorker = true

	go w.Listen()

	return nil
}

func (w *Worker) Listen() {
	for {
		conn, err := w.ln.Accept()

		if err != nil {
			// error accepting incoming connection
			continue
		}

		go w.handleConnection(conn)
	}
}

func (w *Worker) handleConnection(conn net.Conn) {

	reader := bufio.NewReader(conn)

	for {
		_, err := reader.ReadBytes('\n')

		if err != nil {
			// connection closed or read error
			break
		}

		// increment prometheus counter for total msgs received
		msgsRecv.Inc()
	}
}
