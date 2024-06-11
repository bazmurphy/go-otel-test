package main

import (
	"context"
	"flag"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/bazmurphy/go-otel-fun/proto"
	"github.com/bazmurphy/go-otel-fun/util"
)

var (
	destination = flag.String("destination", "", "the address of the grpc server to connect to")
)

func main() {
	flag.Parse()

	connection, err := grpc.NewClient(*destination, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("grpc client could not connect to the grpc server: %v", err)
	}
	defer connection.Close()

	client := pb.NewMyServiceClient(connection)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	source := util.GetIPv4Address()

	request := &pb.MyServiceRequest{
		Source:      source,
		Destination: *destination,
		DataBefore:  100,
	}
	log.Println("DEBUG | request:", request)

	log.Printf("🟧 client making request to: %s", *destination)

	start := time.Now()

	response, err := client.MyServiceProcessData(ctx, request)
	if err != nil {
		log.Printf("failed to send request: %v", err)
	}

	end := time.Now()
	duration := end.Sub(start)

	log.Println("🟩 client received response:", response)
	log.Printf("duration: %v", duration)
}
