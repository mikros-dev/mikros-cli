package plugin

import (
	"log"
	"net"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
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
	println("Listening on " + listener.Addr().String())

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
	// Register gRPC health server
	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(b.server, healthServer)

	println("Starting server")
	go func() {
		if err := b.server.Serve(b.listener); err != nil {
			log.Fatalf("failed to serve: %b", err)
		}
	}()

	healthServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_NOT_SERVING)
	// Wait for the server to become available
	for {
		conn, err := net.Dial("tcp", "localhost:50051")
		if err == nil {
			_ = conn.Close()
			break
		}

		time.Sleep(100 * time.Millisecond)
	}

	// set as running
	time.Sleep(5 * time.Second)
	healthServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)
	//if err := b.setRunning(); err != nil {
	//	return err
	//}

	// wait for stop
	<-b.stop
	return nil
}

func (b *BaseServer) setRunning() error {

	for {
		conn, err := net.Dial("tcp", "localhost:50051")
		if err == nil {
			_ = conn.Close()
			time.Sleep(100 * time.Millisecond)
			break
		}

		time.Sleep(100 * time.Millisecond)
	}

	return os.WriteFile("plugin_ready.txt", []byte("ready"), 0644)
}

func (b *BaseServer) Stop() {
	close(b.stop)

	go func() {
		b.server.GracefulStop()
	}()
}
