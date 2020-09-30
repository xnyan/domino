package node

import (
	"net"
	"time"

	"github.com/op/go-logging"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var logger = logging.MustGetLogger("Node")

// Blocking
func StartServer(s *Server) {
	// Creates the gRPC server instance
	grpcServer := grpc.NewServer()
	RegisterLatencyServer(grpcServer, s)
	reflection.Register(grpcServer)

	// Starts RPC service
	rpcListener, err := net.Listen("tcp", ":"+s.port)
	if err != nil {
		logger.Errorf("Fails to listen on port %s \nError: %v", s.port, err)
	}

	logger.Infof("Starting gRPC servers")

	err = grpcServer.Serve(rpcListener)
	if err != nil {
		logger.Errorf("Cannot start gRPC servers. \nError: %v", err)
	}
}

// Blocking
func StartClient(c *Client, inv, duration time.Duration) {
	c.StartProbing(inv, duration)
}
