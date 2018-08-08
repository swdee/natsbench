package producer

import (
	"bufio"
	"math/rand"
	"net"
	"time"
	"github.com/nats-io/go-nats"
	"crypto/sha256"
	"encoding/hex"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

var err error

type Producer struct {
	size     int
	rate     float64
	mode     string
	ticker   *time.Ticker
	worker   net.Conn
	writer   *bufio.Writer
	conn     *nats.Conn
	isWorker bool
}

func NewProducer(s int, r int, m string) *Producer {
	return &Producer{
		size: s,
		rate: float64(r),
		mode: m,
	}
}

func (p *Producer) ConnectWorker(url string) error {

	p.worker, err = net.Dial("tcp", url)

	if err != nil {
		return err
	}

	p.isWorker = true
	p.writer = bufio.NewWriter(p.worker)

	go p.Listen()

	return nil
}

func (p *Producer) ConnectNats(url string) error {

	p.conn, err = nats.Connect(url)

	if err != nil {
		return err
	}

	p.isWorker = false

	go p.Listen()

	return nil
}

func (p *Producer) Listen() {

	initRand := rand.New(rand.NewSource(time.Now().UnixNano()))

	p.setMsgRate(p.rate)

	for {
		select {
		case <-p.ticker.C:
			if p.isWorker == true {
				// send to Worker
				_, err := p.writer.WriteString(p.randString(initRand, p.size) + "\n")

				if err != nil {
					panic("Error writing to Worker: "+ err.Error())
				}

				err = p.writer.Flush()

				if err != nil {
					panic("Error flushing write to Worker: "+ err.Error())
				}
			} else {
				// send to NATS
				data := []byte(p.randString(initRand, p.size))

				subj := "single"

				if p.mode == "shard" {
					// use the hash of the message for subject publishing
					hasher := sha256.New()
					hasher.Write(data)
					hash := hex.EncodeToString(hasher.Sum(nil))

					subj = "txn." + hash[0:1] + "." + hash[1:2] + "." + hash[2:3]
				}

				p.conn.Publish(subj, data)
				//p.conn.Flush()

				if err := p.conn.LastError(); err != nil {
					panic("Error publising to NATS, error: "+ err.Error())
				}
			}

			// increment prometheus counter for total msgs sent
			msgsSent.Inc()
		}
	}
}

func (p *Producer) setMsgRate(r float64) {

	// stop current ticker
	if p.ticker != nil {
		p.ticker.Stop()
	}

	p.rate = r

	// init ticker, a nanosecond = 1 billion, create ticker for number of msgs/sec
	if p.rate > 0 {
		p.ticker = time.NewTicker(time.Duration(int(1000000000/r)) * time.Nanosecond)
	}
}

func (p *Producer) randString(r *rand.Rand, str_length int) string {
	b := make([]byte, str_length)
	for i := range b {
		b[i] = letterBytes[r.Intn(len(letterBytes))]
	}
	return string(b)
}
