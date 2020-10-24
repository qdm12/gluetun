package httpproxy

import (
	"context"
	"net"
)

type ContextConn interface {
	Read(b []byte) (n int, err error)
	Write(b []byte) (n int, err error)
	Close() error
}

func newContextConn(ctx context.Context, conn net.Conn) ContextConn {
	return &contextConn{
		ctx:  ctx,
		conn: conn,
	}
}

type contextConn struct {
	ctx  context.Context
	conn net.Conn
}

func (c *contextConn) Read(b []byte) (n int, err error) {
	if err := c.ctx.Err(); err != nil {
		return 0, err
	}
	return c.conn.Read(b)
}

func (c *contextConn) Write(b []byte) (n int, err error) {
	if err := c.ctx.Err(); err != nil {
		return 0, err
	}
	return c.conn.Write(b)
}

func (c *contextConn) Close() error {
	return c.conn.Close()
}
