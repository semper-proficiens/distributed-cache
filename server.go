package main

import (
	"encoding/binary"
	"fmt"
	"github.com/semper-proficiens/distributed-cache/proto"
	"io"
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
	members map[net.Conn]struct{}
	cache   Cacher
}

func NewServer(opts ServerOpts, c Cacher) *Server {
	return &Server{
		ServerOpts: opts,
		members:    make(map[net.Conn]struct{}),
		cache:      c,
	}
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.ListenAddr)
	if err != nil {
		return fmt.Errorf("listen error: %w", err)
	}
	log.Printf("server starting on port [%s]\n", s.ListenAddr)

	// dial leader if not leader
	if !s.IsLeader && len(s.LeaderAddr) != 0 {
		log.Printf("Attempting to connect to leader at: %s", s.LeaderAddr)
		go func() {
			if err = s.dialLeader(); err != nil {
				log.Println(err)
			}
		}()
	}

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

func (s *Server) dialLeader() error {
	conn, err := net.Dial("tcp", s.LeaderAddr)
	if err != nil {
		return fmt.Errorf("failed to connect to leader at [%s]: %v", s.LeaderAddr, err)
	}
	log.Println("connected to leader:", s.LeaderAddr)

	binary.Write(conn, binary.LittleEndian, proto.CmdJoin)
	s.handleConn(conn)
	return nil
}

func (s *Server) handleConn(conn net.Conn) {
	defer func() {
		if err := conn.Close(); err != nil {
			log.Printf("error closing connection: %s\n", err)
		}
	}()

	for {
		cmd, err := proto.ParseCommand(conn)
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Println("parse command error:", err)
			break
		}

		go s.handleCommand(conn, cmd)
	}
}

func (s *Server) handleCommand(conn net.Conn, cmd any) {
	switch v := cmd.(type) {
	case *proto.CommandSet:
		s.handleSetCommand(conn, v)
	case *proto.CommandGet:
		s.handleGetCommand(conn, v)
	case *proto.CommandJoin:
		s.handleJoinCommand(conn, v)
	}
}

func (s *Server) handleJoinCommand(conn net.Conn, cmd *proto.CommandJoin) error {
	fmt.Println("member just joined the cluster:", conn.RemoteAddr())
	s.members[conn] = struct{}{}
	return nil
}

func (s *Server) handleGetCommand(conn net.Conn, cmd *proto.CommandGet) error {
	//log.Printf("GET %s", cmd.Key)
	resp := proto.ResponseGet{}
	value, err := s.cache.Get(cmd.Key)
	if err != nil {
		resp.Status = proto.StatusError
		_, err = conn.Write(resp.Bytes())
		return err
	}
	resp.Status = proto.StatusOK
	resp.Value = value
	_, err = conn.Write(resp.Bytes())

	return err
}

func (s *Server) handleSetCommand(conn net.Conn, cmd *proto.CommandSet) error {
	//log.Printf("SET %s to %s", cmd.Key, cmd.Value)
	resp := proto.ResponseSet{}
	if err := s.cache.Set(cmd.Key, cmd.Value, time.Duration(cmd.TTL)); err != nil {
		resp.Status = proto.StatusError
		conn.Write(resp.Bytes())
		return err
	}
	resp.Status = proto.StatusOK
	_, err := conn.Write(resp.Bytes())

	return err
}
