package main

import (
	"fmt"
	"log"
	"net"
	"time"
)

type ServerOpts struct {
	ListenAddr string
	IsLeader   bool
	LeaderAddr string
}

type Server struct {
	ServerOpts
	cache Cacher
}

func NewServer(opts ServerOpts, c Cacher) *Server {
	return &Server{
		ServerOpts: opts,
		cache:      c,
	}
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.ListenAddr)
	if err != nil {
		return fmt.Errorf("listen error: %w", err)
	}
	log.Printf("server starting on port [%s]\n", s.ListenAddr)

	for {
		// a break here means it stops accepting connections
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("accept error: %v\n", err)
			continue
		}
		// not handling errs
		go s.handleConn(conn)
	}
}

func (s *Server) handleConn(conn net.Conn) {
	defer func() {
		if err := conn.Close(); err != nil {
			log.Printf("error closing connection: %s\n", err)
		}
	}()

	for {
		cmd, err := ParseCommand(conn)
		if err != nil {
			log.Println("parse command error:", err)
			break
		}
		fmt.Println(cmd)

		go s.handleCommand(conn, cmd)
	}
}

func (s *Server) handleCommand(conn net.Conn, cmd any) {
	switch v := cmd.(type) {
	case *CommandSet:
		s.handleSetCommand(conn, v)
	case *CommandGet:
	}
}

func (s *Server) handleSetCommand(conn net.Conn, cmd *CommandSet) error {
	log.Printf("SET %s to %s", cmd.Key, cmd.Value)
	return s.cache.Set(cmd.Key, cmd.Value, time.Duration(cmd.TTL))
}
