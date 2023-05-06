package rtsp

import (
	"context"
	"io"
	"net"
)

// Client表示客户端
type Client struct {
	// 底层连接
	conn net.Conn
}

// Close关闭连接
func (c *Client) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// Request发起请求
func (c *Client) Request(req *Request, body io.Reader, res *Response) error {
	// header
	err := req.Encode(c.conn)
	if err != nil {
		return err
	}
	// body
	_, err = io.Copy(c.conn, body)
	if err != nil {
		return err
	}
	// response
	err = res.Decode(c.conn)
	if err != nil {
		return err
	}
	// 返回
	return nil
}

// Dial创建连接并返回新的Client
func Dial(ctx context.Context, network, address string) (*Client, error) {
	// 创建连接
	var dialer net.Dialer
	conn, err := dialer.DialContext(ctx, network, address)
	if err != nil {
		return nil, err
	}
	client := new(Client)
	client.conn = conn
	// 返回
	return client, nil
}
