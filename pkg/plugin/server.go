package plugin

import (
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
)

type BaseServer struct {
	stop     chan bool
	listener net.Listener
	server   *grpc.Server
}

func NewBaseServer() (*BaseServer, error) {
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		return nil, err
	}

	return &BaseServer{
		stop:     make(chan bool),
		listener: listener,
		server:   grpc.NewServer(),
	}, nil
}

func (b *BaseServer) GetServer() *grpc.Server {
	return b.server
}

func (b *BaseServer) Run() error {
	go func() {
		if err := b.server.Serve(b.listener); err != nil {
			log.Fatalf("failed to serve: %b", err)
		}
	}()

	// set as running
	if err := b.setRunning(); err != nil {
		return err
	}

	// wait for stop
	<-b.stop
	return nil
}

func (b *BaseServer) setRunning() error {
	return os.WriteFile("plugin_ready.txt", []byte("ready"), 0644)
}

func (b *BaseServer) Stop() {
	close(b.stop)

	go func() {
		b.server.GracefulStop()
	}()
}
