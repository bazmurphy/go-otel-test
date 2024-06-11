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

	if *portFlag == 0 {
		log.Fatalf("'port' flag required")
	}

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

	log.Printf("🤖 Server | IP %s Port %v", source, *portFlag)

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
	log.Println("🟪 Server | MyServiceProcessData | received request...")

	// (!) actually the context carries the request information including source/destination
	// log.Println("DEBUG | Server | ctx:", ctx)

	valueToAdd := rand.Intn(50)

	// increment the value of data (to emulate some work)
	dataAfter := request.DataBefore + int64(valueToAdd)

	// add a delay (to emulate that work taking time)
	delay := time.Duration(rand.Intn(500)) * time.Millisecond
	log.Printf("⏳ Server | emulating work, artificially waiting %v...", delay)
	time.Sleep(delay)

	// if there is a forward server, then forward request to the next server
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
		log.Println("⬜ Server | request:", request)

		log.Printf("🟨 Server | forwarding request to: %s", *forwardServerFlag)

		response, err := client.MyServiceProcessData(ctx, request)
		if err != nil {
			log.Printf("failed to forward request: %v", err)
			return nil, err
		}

		return response, nil
	}

	// otherwise respond to the origin caller
	response := &pb.MyServiceResponse{
		Origin:      request.Origin,
		Source:      serverIP,
		Destination: request.Origin,
		DataAfter:   dataAfter,
	}
	log.Println("⬜ Server | response:", response)

	log.Println("🟦 Server | MyServiceProcessData | sending response...")
	return response, nil
}
