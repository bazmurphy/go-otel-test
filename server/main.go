package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/bazmurphy/go-otel-fun/proto"
	"github.com/bazmurphy/go-otel-fun/util"
)

var (
	port          = flag.Int("port", 0, "the grpc server port")
	forwardServer = flag.String("forward", "", "the address of the next grpc server")
)

func main() {
	flag.Parse()

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	myServiceServer := &MyServiceServer{
		forwardServer: *forwardServer,
	}

	pb.RegisterMyServiceServer(grpcServer, myServiceServer)

	log.Printf("🤖 MyServiceServer listening on %v", listener.Addr())

	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatalf("failed to run the grpc server: %v", err)
	}
}

type MyServiceServer struct {
	pb.UnimplementedMyServiceServer
	forwardServer string
}

func (s *MyServiceServer) MyServiceProcessData(ctx context.Context, request *pb.MyServiceRequest) (*pb.MyServiceResponse, error) {
	log.Println("🟪 MyServiceProcessData received request")
	// log.Println("DEBUG | ctx:", ctx)

	delay := time.Duration(rand.Intn(500)) * time.Millisecond
	log.Printf("⏳ artificially waiting %v...", delay)
	time.Sleep(delay)

	source := util.GetIPv4Address()
	destination := request.Source
	dataAfter := request.DataBefore + 50

	// forward request to the next server
	if s.forwardServer != "" {
		connection, err := grpc.NewClient(s.forwardServer, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("grpc client could not connect to the grpc server: %v", err)
		}
		defer connection.Close()

		client := pb.NewMyServiceClient(connection)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		request.Source = source
		request.Destination = s.forwardServer
		request.DataBefore = dataAfter
		log.Println("DEBUG | request:", request)

		log.Printf("🟨 forwarding request to: %s", *forwardServer)

		response, err := client.MyServiceProcessData(ctx, request)
		if err != nil {
			log.Printf("failed to forward request: %v", err)
			return nil, err
		}

		return response, nil
	}

	response := &pb.MyServiceResponse{
		Source:      source,
		Destination: destination,
		DataAfter:   dataAfter,
	}
	log.Println("DEBUG | response:", response)

	log.Println("🟩 MyServiceProcessData sending response")
	return response, nil
}
