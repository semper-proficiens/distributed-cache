package client

import (
	"context"
	"fmt"
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

func (cc *CacheClient) Get(ctx context.Context, key []byte) ([]byte, error) {
	cmd := &proto.CommandGet{
		Key: key,
	}

	_, err := cc.conn.Write(cmd.Bytes())
	if err != nil {
		return nil, err
	}

	resp, err := proto.ParseGetResponse(cc.conn)
	if err != nil {
		return nil, err
	}
	if resp.Status == proto.StatusNotFound {
		return nil, fmt.Errorf("key not found [%s]", key)
	}
	if resp.Status != proto.StatusOK {
		return nil, fmt.Errorf("server responded with non OK status: [%s]", resp.Status)
	}

	return resp.Value, nil
}

func (cc *CacheClient) Set(ctx context.Context, key, value []byte, ttl int) error {
	cmd := &proto.CommandSet{
		Key:   key,
		Value: value,
		TTL:   ttl,
	}

	_, err := cc.conn.Write(cmd.Bytes())
	if err != nil {
		return err
	}

	resp, err := proto.ParseSetResponse(cc.conn)
	if err != nil {
		return err
	}
	if resp.Status != proto.StatusOK {
		return fmt.Errorf("server responded with non OK status: [%s]", resp.Status)
	}

	return nil
}

func (cc *CacheClient) Close() error {
	return cc.conn.Close()
}
