package main

import (
	"flag"
	"log"
	"net"
	"time"
)

func main() {
	var (
		listenAddr = flag.String("listenaddr", ":3000", "listening address of the server")
		leaderAddr = flag.String("leaderaddr", ":", "listening address of the leader")
	)
	flag.Parse()
	opts := ServerOpts{
		ListenAddr: *listenAddr,
		IsLeader:   len(*leaderAddr) == 0,
		LeaderAddr: *leaderAddr,
	}

	go func() {
		time.Sleep(time.Second * 2)
		for i := 0; i < 10; i++ {
			SendCommand()
			time.Sleep(time.Millisecond * 200)
		}
	}()

	server := NewServer(opts, NewCache())
	server.Start()
}

func SendCommand() {
	cmd := &CommandSet{
		Key:   []byte("Foo"),
		Value: []byte("Bar"),
		TTL:   2,
	}

	conn, err := net.Dial("tcp", ":3000")
	if err != nil {
		log.Fatal(err)
	}

	conn.Write(cmd.Bytes())
}
