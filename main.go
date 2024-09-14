package main

import (
	"context"
	"crypto/rand"
	"flag"
	"fmt"
	"github.com/semper-proficiens/distributed-cache/client"
	"io"
	"log"
)

func main() {
	var (
		listenAddr = flag.String("listenaddr", ":3000", "listening address of the server")
		leaderAddr = flag.String("leaderaddr", "", "listening address of the leader")
	)
	flag.Parse()
	opts := ServerOpts{
		ListenAddr: *listenAddr,
		IsLeader:   len(*leaderAddr) == 0,
		LeaderAddr: *leaderAddr,
	}

	//go func() {
	//	// time to boot server
	//	time.Sleep(time.Second * 2)
	//	sendStuff()
	//}()

	server := NewServer(opts, NewCache())
	server.Start()
}

func randomBytes(n int) []byte {
	b := make([]byte, n)
	io.ReadFull(rand.Reader, b)
	return b
}

func sendStuff() {
	for i := 0; i < 100; i++ {
		go func() {
			cc, err := client.NewCacheClient(":3000", client.Options{})
			if err != nil {
				log.Fatal(err)
			}

			var (
				key   = []byte(fmt.Sprintf("key-%d", i))
				value = []byte(fmt.Sprintf("value-%d", i))
			)

			err = cc.Set(context.Background(), key, value, 0)
			if err != nil {
				log.Fatal(err)
			}

			resp, err := cc.Get(context.Background(), key)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(string(resp))

			cc.Close()
		}()
	}
}
