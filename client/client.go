package client

import (
	"context"
	"github.com/semper-proficiens/distributed-cache/proto"
	"net"
)

type CacheClient struct {
	conn net.Conn
}

type Options struct {
}

func NewCacheClient(endpoint string, opts Options) (*CacheClient, error) {
	conn, err := net.Dial("tcp", endpoint)
	if err != nil {
		return nil, err
	}
	return &CacheClient{
		conn: conn,
	}, nil
}

func (cc *CacheClient) Set(ctx context.Context, key, value []byte, ttl int) (any, error) {
	cmd := &proto.CommandSet{
		Key:   key,
		Value: value,
		TTL:   ttl,
	}

	_, err := cc.conn.Write(cmd.Bytes())
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (cc *CacheClient) Close() error {
	return cc.conn.Close()
}
