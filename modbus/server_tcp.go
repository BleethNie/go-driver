package modbus

import (
	"crypto/tls"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
)

func (s *Server) accept(listen net.Listener) error {
	for {
		conn, err := listen.Accept()
		if err != nil {
			if strings.Contains(err.Error(), "use of closed network connection") {
				return nil
			}
			log.Printf("Unable to accept connections: %#v\n", err)
			return err
		}
		go func(conn net.Conn) {
			defer conn.Close()

			for {
				packet := make([]byte, 512)
				bytesRead, err := conn.Read(packet)
				if err != nil {
					if err != io.EOF {
						log.Printf("read error %v\n", err)
					}
					return
				}
				// Set the length of the packet to the number of read bytes.
				packet = packet[:bytesRead]

				frame, err := NewTCPFrame(packet)
				if err != nil {
					log.Printf("bad packet error %v\n", err)
					return
				}

				request := &Request{conn, frame}
				if len(s.IpWhiteList) > 0 {
					addr := conn.RemoteAddr().String()
					ip := strings.Split(addr, ":")[0]
					if !IsInIpList(ip, s.IpWhiteList) {
						frame := request.frame.Copy()
						frame.SetException(&ServerNetworkCheckError)
						conn.Write(frame.Bytes())
						conn.Close()
						return
					}
				}
				s.requestChan <- request
			}
		}(conn)
	}
}

// ListenTCP starts the Modbus server listening on "address:port".
func (s *Server) ListenTCP(addressPort string) (err error) {
	listen, err := net.Listen("tcp", addressPort)
	if err != nil {
		log.Printf("Failed to Listen: %v\n", err)
		return err
	}
	s.listeners = append(s.listeners, listen)
	go s.accept(listen)
	return err
}

// ListenTLS starts the Modbus server listening on "address:port".
func (s *Server) ListenTLS(addressPort string, config *tls.Config) (err error) {
	listen, err := tls.Listen("tcp", addressPort, config)
	if err != nil {
		log.Printf("Failed to Listen on TLS: %v\n", err)
		return err
	}
	s.listeners = append(s.listeners, listen)
	go s.accept(listen)
	return err
}

func IsInIpList(ip string, ipList []string) bool {
	if ip == "127.0.0.1" || ip == "localhost" {
		return true
	}
	//解析单个IP
	ipAddr := net.ParseIP(ip)
	if ipAddr == nil {
		return false // IP地址无效
	}
	// 遍历IP列表
	for _, ipRange := range ipList {
		// 分割IP和掩码长度
		parts := strings.Split(ipRange, "/")

		// 检查是否是有效的CIDR表示法
		if len(parts) != 2 {
			if ip == parts[0] {
				return true
			}
			continue
		}
		// 解析IP地址部分
		intValue, _ := strconv.ParseInt(parts[1], 0, 32)
		ipNet := net.IPNet{IP: net.ParseIP(parts[0]), Mask: net.CIDRMask(int(intValue), 32)}
		// 检查IP是否在当前的IPNet范围内
		if ipNet.Contains(ipAddr) {
			return true
		}
	}
	return false
}
