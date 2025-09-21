package transport

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"sync"
)

type TcpServer struct {
	listener net.Listener
	port     string
	handlers map[string]TcpHandler
	mu       sync.RWMutex
	ctx      context.Context
	cancel   context.CancelFunc
}

type TcpHandler func([]byte) ([]byte, error)

type TcpMessage struct {
	Action string          `json:"action"`
	Data   json.RawMessage `json:"data"`
}

type TcpResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func NewTcpServer(port string) *TcpServer {
	ctx, cancel := context.WithCancel(context.Background())
	return &TcpServer{
		port:     port,
		handlers: make(map[string]TcpHandler),
		ctx:      ctx,
		cancel:   cancel,
	}
}

func (s *TcpServer) RegisterHandler(action string, handler TcpHandler) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.handlers[action] = handler
}

func (s *TcpServer) Start() error {
	listener, err := net.Listen("tcp", ":"+s.port)
	if err != nil {
		return err
	}
	s.listener = listener

	fmt.Printf("TCP Server listening on port %s\n", s.port)

	for {
		select {
		case <-s.ctx.Done():
			return nil
		default:
			conn, err := listener.Accept()
			if err != nil {
				continue
			}
			go s.handleConnection(conn)
		}
	}
}

func (s *TcpServer) handleConnection(conn net.Conn) {
	defer conn.Close()

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		data := scanner.Bytes()
		response := s.processMessage(data)

		responseData, err := json.Marshal(response)
		if err != nil {
			continue
		}

		conn.Write(responseData)
		conn.Write([]byte("\n"))
	}
}

func (s *TcpServer) processMessage(data []byte) TcpResponse {
	var msg TcpMessage
	if err := json.Unmarshal(data, &msg); err != nil {
		return TcpResponse{
			Success: false,
			Error:   "Invalid message format",
		}
	}

	s.mu.RLock()
	handler, exists := s.handlers[msg.Action]
	s.mu.RUnlock()

	if !exists {
		return TcpResponse{
			Success: false,
			Error:   fmt.Sprintf("Unknown action: %s", msg.Action),
		}
	}

	result, err := handler(msg.Data)
	if err != nil {
		return TcpResponse{
			Success: false,
			Error:   err.Error(),
		}
	}

	var resultData interface{}
	json.Unmarshal(result, &resultData)

	return TcpResponse{
		Success: true,
		Data:    resultData,
	}
}

func (s *TcpServer) Stop() error {
	s.cancel()
	if s.listener != nil {
		return s.listener.Close()
	}
	return nil
}

type TcpClient struct {
	conn net.Conn
	host string
	port string
}

func NewTcpClient(host, port string) *TcpClient {
	return &TcpClient{host: host, port: port}
}

func (c *TcpClient) Connect() error {
	conn, err := net.Dial("tcp", c.host+":"+c.port)
	if err != nil {
		return err
	}
	c.conn = conn
	return nil
}

func (c *TcpClient) SendMessage(action string, data interface{}) (*TcpResponse, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	msg := TcpMessage{
		Action: action,
		Data:   jsonData,
	}

	msgData, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}

	_, err = c.conn.Write(msgData)
	if err != nil {
		return nil, err
	}

	_, err = c.conn.Write([]byte("\n"))
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(c.conn)
	if scanner.Scan() {
		var response TcpResponse
		err := json.Unmarshal(scanner.Bytes(), &response)
		if err != nil {
			return nil, err
		}
		return &response, nil
	}

	return nil, fmt.Errorf("no response received")
}

func (c *TcpClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}