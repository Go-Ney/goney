package transport

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
)

type NatsClient struct {
	conn *nats.Conn
	url  string
}

type NatsHandler func([]byte) ([]byte, error)

type NatsSubscription struct {
	subscription *nats.Subscription
}

func NewNatsClient(url string) *NatsClient {
	return &NatsClient{url: url}
}

func (c *NatsClient) Connect() error {
	conn, err := nats.Connect(c.url,
		nats.ReconnectWait(time.Second*2),
		nats.MaxReconnects(5),
		nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
			fmt.Printf("NATS disconnected: %v\n", err)
		}),
		nats.ReconnectHandler(func(nc *nats.Conn) {
			fmt.Printf("NATS reconnected to %v\n", nc.ConnectedUrl())
		}),
	)
	if err != nil {
		return err
	}
	c.conn = conn
	return nil
}

func (c *NatsClient) Publish(subject string, data interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return c.conn.Publish(subject, jsonData)
}

func (c *NatsClient) Request(subject string, data interface{}, timeout time.Duration) (*nats.Msg, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return c.conn.Request(subject, jsonData, timeout)
}

func (c *NatsClient) Subscribe(subject string, handler NatsHandler) (*NatsSubscription, error) {
	sub, err := c.conn.Subscribe(subject, func(msg *nats.Msg) {
		response, err := handler(msg.Data)
		if err != nil {
			fmt.Printf("Error handling message: %v\n", err)
			return
		}
		if msg.Reply != "" {
			c.conn.Publish(msg.Reply, response)
		}
	})
	if err != nil {
		return nil, err
	}
	return &NatsSubscription{subscription: sub}, nil
}

func (c *NatsClient) QueueSubscribe(subject, queue string, handler NatsHandler) (*NatsSubscription, error) {
	sub, err := c.conn.QueueSubscribe(subject, queue, func(msg *nats.Msg) {
		response, err := handler(msg.Data)
		if err != nil {
			fmt.Printf("Error handling message: %v\n", err)
			return
		}
		if msg.Reply != "" {
			c.conn.Publish(msg.Reply, response)
		}
	})
	if err != nil {
		return nil, err
	}
	return &NatsSubscription{subscription: sub}, nil
}

func (c *NatsClient) Close() {
	if c.conn != nil {
		c.conn.Close()
	}
}

func (s *NatsSubscription) Unsubscribe() error {
	return s.subscription.Unsubscribe()
}