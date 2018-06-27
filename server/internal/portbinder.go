package internal

import (
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	defaultPort = 5556
)

// PortBinder binds to a tcp port
type PortBinder struct {
	Port int32
}

// Bind !
func (p *PortBinder) Bind(server *grpc.Server) {
	port := p.Port
	if port == 0 {
		port = defaultPort
	}

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", p.Port))
	log.Printf("Binding to port %d\n", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	reflection.Register(server)
	if err := server.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
