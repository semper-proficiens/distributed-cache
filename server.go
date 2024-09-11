package main

//
//import (
//	"context"
//	"fmt"
//	"log"
//	"net"
//)
//
//type ServerOpts struct {
//	ListenAddr string
//	IsLeader   bool
//	LeaderAddr string
//}
//
//type RaftServer struct {
//	ServerOpts
//	followers map[net.Conn]struct{}
//	cache     Cacher
//}
//
//func NewRaftServer(opts ServerOpts, c Cacher) *RaftServer {
//	return &RaftServer{
//		ServerOpts: opts,
//		cache:      c,
//		followers:  make(map[net.Conn]struct{}),
//	}
//}
//
//func (s *RaftServer) Start() error {
//	ln, err := net.Listen("tcp", s.ListenAddr)
//	if err != nil {
//		return fmt.Errorf("listen error: %w", err)
//	}
//	log.Printf("server starting on port [%s]\n", s.ListenAddr)
//
//	if !s.IsLeader {
//		go func() {
//			conn, err := net.Dial("tcp", s.LeaderAddr)
//			fmt.Println("connected with leader:", s.LeaderAddr)
//			if err != nil {
//				log.Fatal(err)
//			}
//			s.handleConn(conn)
//		}()
//	}
//	for {
//		// a break here means it stops accepting connections
//		conn, err := ln.Accept()
//		if err != nil {
//			log.Printf("accept error: %v\n", err)
//			continue
//		}
//		// not handling errs
//		go s.handleConn(conn)
//	}
//}
//
//func (s *RaftServer) handleConn(conn net.Conn) {
//	defer func() {
//		if err := conn.Close(); err != nil {
//			log.Printf("error closing connection: %s\n", err)
//		}
//	}()
//	buff := make([]byte, 2048)
//
//	if s.IsLeader {
//		s.followers[conn] = struct{}{}
//	}
//
//	fmt.Println("connection made:", conn.RemoteAddr())
//	for {
//		// break conn if we can't read it
//		n, err := conn.Read(buff)
//		if err != nil {
//			log.Printf("conn read error: %s\n", err)
//			break
//		}
//		msg := string(buff[:n])
//		fmt.Println(msg)
//	}
//}
//
//func (s *RaftServer) handleCommand(conn net.Conn, rawCmd []byte) {
//	msg, err := parseMessage(rawCmd)
//	if err != nil {
//		fmt.Println("failed to parse command", err)
//		return
//	}
//	switch msg.Cmd {
//	case CMDSet:
//		err = s.handleSetCmd(conn, msg)
//	case CMDGet:
//		err = s.handleGetCmd(conn, msg)
//	}
//	if err != nil {
//		conn.Write([]byte(err.Error()))
//	}
//}
//
//func (s *RaftServer) handleGetCmd(conn net.Conn, msg *Message) error {
//	val, err := s.cache.Get(msg.Key)
//	if err != nil {
//		return err
//	}
//	_, err = conn.Write(val)
//	return err
//}
//
//func (s *RaftServer) handleSetCmd(conn net.Conn, msg *Message) error {
//	if err := s.cache.Set(msg.Key, msg.Value, msg.TTL); err != nil {
//		return err
//	}
//	go s.sendToFollowers(context.TODO(), msg)
//	return nil
//}
//
//func (s *RaftServer) sendToFollowers(ctx context.Context, msg *Message) error {
//	for conn := range s.followers {
//		fmt.Println("forward to follower")
//		rawMsg := msg.ToBytes()
//		fmt.Println("forward rawMsg to follower:", rawMsg)
//		_, err := conn.Write(rawMsg)
//		if err != nil {
//			fmt.Println("write to follower error:", err)
//			continue
//		}
//	}
//	return nil
//}
