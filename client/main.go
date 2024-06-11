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
	clientIP    = util.GetIPv4Address()
)

func main() {
	flag.Parse()

	if *destination == "" {
		log.Fatalf("'destination' flag required")
	}

	log.Println("⬜ Client | clientIP:", clientIP)

	connection, err := grpc.NewClient(*destination, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("grpc client could not connect to the grpc server: %v", err)
	}
	defer connection.Close()

	client := pb.NewMyServiceClient(connection)

	// continuously make a request every 3 seconds
	for {
		log.Println("--------------------------------------------------------------")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

		request := &pb.MyServiceRequest{
			Origin:      clientIP,
			Source:      clientIP,
			Destination: *destination,
			DataBefore:  100,
		}
		log.Println("⬜ Client | request:", request)

		log.Printf("🟦 Client | making request to: %s", *destination)

		start := time.Now()

		response, err := client.MyServiceProcessData(ctx, request)
		if err != nil {
			log.Printf("failed to send request: %v", err)
		}

		end := time.Now()
		duration := end.Sub(start)

		log.Println("🟩 Client | received response:", response)
		log.Printf("🟩 Client | total duration: %v", duration)

		// cancel the context (without defer above)
		cancel()

		time.Sleep(3 * time.Second)
	}
}
