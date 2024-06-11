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
	portFlag          = flag.Int("port", 0, "the grpc server port")
	forwardServerFlag = flag.String("forward", "", "the address of the next grpc server")
	serverIP          = util.GetIPv4Address()
)

func main() {
	flag.Parse()

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", *portFlag))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	myServiceServer := &MyServiceServer{
		forwardServer: *forwardServerFlag,
	}

	pb.RegisterMyServiceServer(grpcServer, myServiceServer)

	source := util.GetIPv4Address()

	log.Printf("🤖 MyServiceServer %s %v", source, listener.Addr())

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
	log.Println("🟪 Server | MyServiceProcessData received request...")
	// log.Println("DEBUG | ctx:", ctx)

	delay := time.Duration(rand.Intn(500)) * time.Millisecond
	log.Printf("⏳ Server | artificially waiting %v...", delay)
	time.Sleep(delay)

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

		request.Source = serverIP
		request.Destination = s.forwardServer
		request.DataBefore = dataAfter
		log.Println("ℹ️ Server | request:", request)

		log.Printf("🟨 Server | forwarding request to: %s", *forwardServerFlag)

		response, err := client.MyServiceProcessData(ctx, request)
		if err != nil {
			log.Printf("failed to forward request: %v", err)
			return nil, err
		}

		return response, nil
	}

	response := &pb.MyServiceResponse{
		Origin:      request.Origin,
		Source:      serverIP,
		Destination: destination,
		DataAfter:   dataAfter,
	}
	log.Println("ℹ️ Server | response:", response)

	log.Println("🟦 Server | MyServiceProcessData sending response...")
	return response, nil
}
