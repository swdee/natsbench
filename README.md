# NATS.io Message Throughput Tests

Code consists of two programs;
 * Producer - used for creating random messages to send at a specified size and rate
 * Worker - receives messages from Producer or NATS depending on the runtime mode

These program allow us to compare message through put of NATS versus a simple TCP server.

## Building

1. Checkout GIT repository 

2. Build executable
```
go build producer.go
go build worker.go
```


## Usage

### Worker

The Worker can be run in two modes, the first connects to NATS and subscribes to the Subject names, the 
second mode acts as a stand alone TCP server that receives messages directly from the Producer.

To run via NATS
```
./worker2 -n nats://<ip>:<port>
```

To run as a standalone TCP server on port 8095
```
./worker2 -s :8095
```

### Producer

The Producer is run to send messages to NATS or to the Worker directly for comparing message through put of simple TCP server.

To run via NATS at message rate 10,000 per second.  Default message size 32 bytes.
```
./producer2 -d nats://<ip>:4222 -r 10000
```

To send messages direct to Worker TCP server
```
./producer2 -d worker://<ip>:8095 -r 10000
```


