// Package mbserver implments a Modbus server (slave).
package modbus

import (
	"encoding/hex"
	"fmt"
	"github.com/goburrow/serial"
	"io"
	"net"
	"sync"
)

// Server is a Modbus slave with allocated memory for discrete inputs, coils, etc.
type Server struct {
	// Debug enables more verbose messaging.
	Debug            bool
	listeners        []net.Listener
	ports            []serial.Port
	portsWG          sync.WaitGroup
	portsCloseChan   chan struct{}
	requestChan      chan *Request
	function         [256](func(*Server, Framer) ([]byte, *Exception))
	DiscreteInputs   []byte
	Coils            []byte
	HoldingRegisters []uint16
	InputRegisters   []uint16
	IpWhiteList      []string
}

// Request contains the connection and Modbus frame.
type Request struct {
	conn  io.ReadWriteCloser
	frame Framer
}

// NewServer creates a new Modbus server (slave).
func NewServer() *Server {
	s := &Server{}

	// Allocate Modbus memory maps.
	s.DiscreteInputs = make([]byte, 65536)
	s.Coils = make([]byte, 65536)
	s.HoldingRegisters = make([]uint16, 65536)
	s.InputRegisters = make([]uint16, 65536)
	s.IpWhiteList = make([]string, 0)

	// Add default functions.
	s.function[1] = ReadCoils
	s.function[2] = ReadDiscreteInputs
	s.function[3] = ReadHoldingRegisters
	s.function[4] = ReadInputRegisters
	s.function[5] = WriteSingleCoil
	s.function[6] = WriteHoldingRegister
	s.function[15] = WriteMultipleCoils
	s.function[16] = WriteHoldingRegisters

	s.requestChan = make(chan *Request)
	s.portsCloseChan = make(chan struct{})

	go s.handler()

	return s
}

// 默认支持localhost和127.0.0.1
func (s *Server) SetIpWhiteList(ips []string) {
	s.IpWhiteList = ips
}

// RegisterFunctionHandler override the default behavior for a given Modbus function.
func (s *Server) RegisterFunctionHandler(funcCode uint8, function func(*Server, Framer) ([]byte, *Exception)) {
	s.function[funcCode] = function
}

func (s *Server) handle(request *Request) Framer {
	var exception *Exception
	var data []byte
	response := request.frame.Copy()

	function := request.frame.GetFunction()
	if s.function[function] != nil {
		data, exception = s.function[function](s, request.frame)
		response.SetData(data)
		hexStr := hex.EncodeToString(response.Bytes())
		fmt.Println("获取到的数据为：" + hexStr)
	} else {
		exception = &IllegalFunction
	}

	if exception != nil && exception != &Success {
		response.SetException(exception)
	}
	return response
}

// All requests are handled synchronously to prevent modbus memory corruption.
func (s *Server) handler() {
	for {
		request := <-s.requestChan
		response := s.handle(request)
		request.conn.Write(response.Bytes())
	}
}

// Close stops listening to TCP/IP ports and closes serial ports.
func (s *Server) Close() {
	for _, listen := range s.listeners {
		listen.Close()
	}
	//发送关闭信号
	s.portsCloseChan <- struct{}{}
	s.portsWG.Wait()

	for _, port := range s.ports {
		port.Close()
	}
}
