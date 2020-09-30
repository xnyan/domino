package node

import (
	"strings"
	"time"

	"golang.org/x/net/context"
)

type Server struct {
	ip   string
	port string
}

func NewServer(addr string) *Server {
	ipPort := strings.Split(addr, ":")
	if len(ipPort) != 2 {
		logger.Fatalf("Invalid server network address = %s", addr)
	}
	return &Server{
		ip:   ipPort[0],
		port: ipPort[1],
	}
}

func (s *Server) Probe(ctx context.Context, req *ProbeReq) (*ProbeReply, error) {
	t := time.Now().UnixNano()
	return &ProbeReply{ServerClock: t}, nil
}
