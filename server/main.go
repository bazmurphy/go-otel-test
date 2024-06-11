package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.25.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/bazmurphy/go-otel-test/proto"
	"github.com/bazmurphy/go-otel-test/util"
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

	// ---------- OTEL START ---------

	ctx := context.Background()

	// create an otel exporter
	// traceExporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	traceExporter, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithEndpoint("jaeger:4317"),
		otlptracegrpc.WithInsecure(),
	)
	if err != nil {
		log.Fatalf("failed to create otel exporter: %v", err)
	}

	// create an otel resource
	resource, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName("server"),
		),
	)
	if err != nil {
		log.Fatalf("failed to create otel resource: %v", err)
	}

	// create an otel tracer provider
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(traceExporter),
		sdktrace.WithResource(resource),
	)

	defer func() { _ = tracerProvider.Shutdown(context.Background()) }()

	// register the tracer provider as the global tracer provider
	otel.SetTracerProvider(tracerProvider)

	// otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	// ---------- OTEL END ---------

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", *portFlag))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer(
		grpc.StatsHandler(otelgrpc.NewServerHandler()), // (!) for otel
		// registers the OpenTelemetry gRPC server handler
		// it automatically extracts trace and span context from incoming gRPC metadata headers
		// and injects them into the server-side context
	)

	myServiceServer := &MyServiceServer{
		forwardServer: *forwardServerFlag,
	}

	pb.RegisterMyServiceServer(grpcServer, myServiceServer)

	source := util.GetIPv4Address()
	log.Printf("🤖 Server | IP: %s Port: %v", source, *portFlag)

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
	log.Println("🟪 Server | received request...")

	// (!) actually the context carries the request information including source/destination
	// log.Println("DEBUG | Server | ctx:", ctx)
	// DEBUG | Server | ctx: context.Background.WithValue(type transport.connectionKey, val <not Stringer>).WithValue(type peer.peerKey, val Peer{Addr: '192.168.32.7:54688', LocalAddr: '192.168.32.2:8081', AuthInfo: <nil>}).WithCancel.WithValue(type metadata.mdIncomingKey, val MD{:authority=[server1:8081], content-type=[application/grpc], user-agent=[grpc-go/1.64.0], grpc-accept-encoding=[gzip]}).WithValue(type grpc.serverKey, val <not Stringer>).WithValue(type trace.traceContextKeyType, val <not Stringer>).WithValue(type trace.traceContextKeyType, val <not Stringer>).WithValue(type otelgrpc.gRPCContextKey, val <not Stringer>).WithValue(type grpc.streamKey, val <not Stringer>)

	// span := trace.SpanFromContext(ctx)
	// log.Printf("🔍 Server | span : %v", span)
	spanContext := trace.SpanContextFromContext(ctx)
	// log.Printf("🔍 Server | spanContext : %v", spanContext)
	traceID := spanContext.TraceID().String()
	spanID := spanContext.SpanID().String()
	log.Printf("🔍 Server | Trace ID: %s Span ID: %s", traceID, spanID)

	// (!) HELP....
	// how to make a child span?
	// do we even need to? why can't it auto follow through the request path...

	tracer := otel.Tracer("server")

	// (!) if the context.Context provided in `ctx` contains a Span then the newly-created Span will be a child of that span
	ctx, childSpan := tracer.Start(ctx, "server-span-test")
	defer childSpan.End()

	valueToAdd := rand.Intn(50)
	// increment the value of data (to emulate some work)
	dataAfter := request.Data + int64(valueToAdd)

	// add a delay (to emulate that work taking time)
	delay := time.Duration(rand.Intn(500)) * time.Millisecond
	log.Printf("⏳ Server | emulating work, adding %d to data, then waiting %v...", valueToAdd, delay)
	time.Sleep(delay)

	// if there is a forward server, then forward request to it
	if s.forwardServer != "" {
		connection, err := grpc.NewClient(
			s.forwardServer,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithStatsHandler(otelgrpc.NewClientHandler()), // for otel
		)
		if err != nil {
			log.Fatalf("grpc client could not connect to the grpc server: %v", err)
		}
		defer connection.Close()

		client := pb.NewMyServiceClient(connection)

		// (!) HELP (!)
		// i need to propagate the incoming request's context (and tracing info) through to the next request
		// but isn't this done automatically by the otelgrpc???

		request.Source = serverIP
		request.Destination = s.forwardServer
		request.Data = dataAfter
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
		Data:        dataAfter,
	}
	log.Println("⬜ Server | response:", response)

	log.Println("🟦 Server | sending response...")
	return response, nil
}
