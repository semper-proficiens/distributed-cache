package main

//
//import (
//	"flag"
//	"fmt"
//	"log"
//	"net"
//	"time"
//)
//
//func main() {
//
//	//conn, err := net.Dial("tcp", ":3000")
//	//if err != nil {
//	//	log.Fatal(err)
//	//}
//	//
//	//_, err = conn.Write([]byte("SET Foo Bar 40000"))
//	//if err != nil {
//	//	log.Fatal(err)
//	//}
//	//
//	//select {}
//
//	var (
//		listenAddr = flag.String("listenaddr", ":3000", "listening address of the server")
//		leaderAddr = flag.String("leaderaddr", ":", "listening address of the leader")
//	)
//	flag.Parse()
//	opts := ServerOpts{
//		ListenAddr: *listenAddr,
//		IsLeader:   len(*leaderAddr) == 0,
//		LeaderAddr: *leaderAddr,
//	}
//	go func() {
//		time.Sleep(time.Second * 2)
//		conn, err := net.Dial("tcp", opts.ListenAddr)
//		if err != nil {
//			log.Fatalf("Failed to connect to server: %v", err)
//		}
//		conn.Write([]byte("SET Foo Bar 250000"))
//
//		time.Sleep(time.Second * 2)
//		conn.Write([]byte("GET Foo Bar 250000"))
//
//		buf := make([]byte, 1024)
//		n, _ := conn.Read(buf)
//		fmt.Println(string(buf[:n]))
//	}()
//
//	server := NewRaftServer(opts, NewCache())
//	server.Start()
//}
