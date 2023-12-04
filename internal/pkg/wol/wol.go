package wol

import (
	"context"
	"net"
)

type Client interface {
	SendSignal(ctx context.Context, mac string) error
}

type wakeOnLanClient struct {
	conn *net.UDPConn
}

func New(conn *net.UDPConn) Client {
	return &wakeOnLanClient{
		conn: conn,
	}
}
