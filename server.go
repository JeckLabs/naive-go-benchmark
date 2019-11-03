package main

import (
	"log"
	"net"
    "encoding/gob"
    "runtime"
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


func handler(conn net.Conn) {
    defer conn.Close()

    enc := gob.NewEncoder(conn)
    dec := gob.NewDecoder(conn)

    var req Request

    for {
        err := dec.Decode(&req)
        if err != nil {
            return
        }
        enc.Encode(Response{MessageType: MessagePong, Payload: []byte("pong")})
        if err != nil {
            return
        }
    }
}

func main() {
    runtime.GOMAXPROCS(1)

    addr := flag.String("addr", "0.0.0.0:6666", "host and port to listen on")
    flag.Parse()

    l, err := net.Listen("tcp4", *addr)
    if err != nil {
        log.Fatal("Error listening:", err.Error())
    }
    defer l.Close()

    for {
        conn, err := l.Accept()
        if err != nil {
            log.Fatal("Error accepting: ", err.Error())
        }
        go handler(conn)
    }
}
