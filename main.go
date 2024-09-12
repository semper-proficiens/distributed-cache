package main

import (
	"context"
	"flag"
	"github.com/semper-proficiens/distributed-cache/client"
	"log"
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
		cc, err := client.NewCacheClient(":3000", client.Options{})
		if err != nil {
			log.Fatal(err)
		}
		for i := 0; i < 10; i++ {
			SendCommand(cc)
		}
		// test that we can close safely with no error. TODO remove this
		cc.Close()
		time.Sleep(time.Second * 1)
	}()

	server := NewServer(opts, NewCache())
	server.Start()
}

func SendCommand(cc *client.CacheClient) {
	_, err := cc.Set(context.Background(), []byte("Ernie"), []byte("GO"), 0)
	if err != nil {
		log.Fatal(err)
	}
}
