package main

import (
	"log"
	"net"
    "encoding/gob"
    "time"
    "sync/atomic"
    "flag"
)

type MessageType int
const (
    MessagePing = iota
    MessagePong
)

type Request struct {
    MessageType MessageType
    Payload []byte
}

type Response struct {
    MessageType MessageType
    Payload []byte
}

var requests uint64

func doRequests(conn net.Conn, c chan<- bool) {
    defer conn.Close()
    defer func () {c <- true}()

    enc := gob.NewEncoder(conn)
    dec := gob.NewDecoder(conn)
    var res Response

    go func () {
        for {
            enc.Encode(Request{MessageType: MessagePing, Payload: []byte("random payload...")})
        }
    }()

    for {
        err := dec.Decode(&res)
        if err != nil {
            log.Println("Error reading:", err.Error())
            return
        }
        atomic.AddUint64(&requests, 1)
    }
}

func printStat() {
    for {
        last := requests
        time.Sleep(time.Second)
        log.Printf("Requests: %d rps: %d", requests, requests - last)
    }
}

func main() {
    addr := flag.String("addr", "127.0.0.1:6666", "host and port to connect to")
    connections := flag.Int("c", 4, "number of concurrent connections")
    flag.Parse()

    c := make(chan bool)

    socks := make([]net.Conn, *connections)

    for i := 0; i < *connections; i++ {
        conn, err := net.Dial("tcp4", *addr)
        if err != nil {
            log.Fatal("Error listening:", err.Error())
        }
        socks[i] = conn
    }

    for i := 0; i < *connections; i++ {
        go doRequests(socks[i], c)
    }

    go printStat()

    for i := 0; i < *connections; i++ {
        <-c
    }

}
