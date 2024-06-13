package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net"
	"strings"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.25.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/peer"

	pb "github.com/bazmurphy/go-otel-test/proto"
	"github.com/bazmurphy/go-otel-test/util"
)

var (
	portFlag          = flag.Int("port", 0, "the grpc server port")
	forwardServerFlag = flag.String("forward", "", "the address of the next grpc server")
	serverID          = flag.String("id", "", "the server's unique identifier")
	serverIP          = util.GetIPv4Address()
)

func main() {
	flag.Parse()

	if *portFlag == 0 {
		log.Fatalf("'port' flag required")
	}

	if *serverID == "" {
		log.Fatalf("'id' flag required")
	}

	// ---------- OTEL SETUP START ---------

	grpcTraceClient := otlptracegrpc.NewClient()

	ctx := context.Background()

	traceExporter, err := otlptrace.New(ctx, grpcTraceClient)
	if err != nil {
		log.Fatalf("failed to create otel trace exporter: %v", err)
	}

	defer func() {
		err := traceExporter.Shutdown(ctx)
		if err != nil {
			log.Fatalf("failed shutting down otel trace exporter: %v", err)
		}
	}()

	resource, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName("server-"+*serverID+"-service-name"),
		),
	)
	if err != nil {
		log.Fatalf("failed to create otel resource: %v", err)
	}

	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(traceExporter),
		sdktrace.WithResource(resource),
	)

	defer func() {
		err = tracerProvider.Shutdown(ctx)
		if err != nil {
			log.Fatalf("failed shutting down otel tracer provider: %v", err)
		}
	}()

	// register the global tracer provider
	otel.SetTracerProvider(tracerProvider)

	// register the W3C trace context and baggage propagators so data is propagated across services/processes
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		),
	)

	// ---------- OTEL SETUP END ---------

	// create a tracer
	tracer := otel.Tracer("server-" + *serverID + "-tracer")

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", *portFlag))
	if err != nil {
		log.Fatalf("failed to listen on tcp port %d: %v", *portFlag, err)
	}

	grpcServer := grpc.NewServer(
		grpc.StatsHandler(otelgrpc.NewServerHandler()), // (!) for otel
		// registers the OpenTelemetry gRPC server handler
		// it automatically extracts trace and span context from incoming gRPC metadata headers
		// and injects them into the server-side context
	)

	myServiceServer := &MyServiceServer{
		forwardServer: *forwardServerFlag,
		tracer:        tracer,
	}

	pb.RegisterMyServiceServer(grpcServer, myServiceServer)

	source := util.GetIPv4Address()
	log.Printf("🤖 Server%s | IP: %s Port: %v", *serverID, source, *portFlag)

	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatalf("failed to run the grpc server: %v", err)
	}
}

type MyServiceServer struct {
	pb.UnimplementedMyServiceServer
	forwardServer string
	tracer        trace.Tracer
}

func (s *MyServiceServer) MyServiceProcessData(ctx context.Context, request *pb.MyServiceRequest) (*pb.MyServiceResponse, error) {
	log.Printf("🟩 Server%s | Received Request...", *serverID)

	// (!) actually the context carries the request information including source/destination
	// log.Printf("🟫 Server%s | ctx: %v", *serverID, ctx)

	// BEFORE SetTextMapPropagator:
	// DEBUG | Server | ctx: context.Background.WithValue(type transport.connectionKey, val <not Stringer>).WithValue(type peer.peerKey, val Peer{Addr: '192.168.32.7:54688', LocalAddr: '192.168.32.2:8081', AuthInfo: <nil>}).WithCancel.WithValue(type metadata.mdIncomingKey, val MD{:authority=[server1:8081], content-type=[application/grpc], user-agent=[grpc-go/1.64.0], grpc-accept-encoding=[gzip]}).WithValue(type grpc.serverKey, val <not Stringer>).WithValue(type trace.traceContextKeyType, val <not Stringer>).WithValue(type trace.traceContextKeyType, val <not Stringer>).WithValue(type otelgrpc.gRPCContextKey, val <not Stringer>).WithValue(type grpc.streamKey, val <not Stringer>)
	// AFTER SetTextMapPropagator:
	// DEBUG | Server | ctx: context.Background.WithValue(type transport.connectionKey, val <not Stringer>).WithValue(type peer.peerKey, val Peer{Addr: '192.168.176.7:35092', LocalAddr: '192.168.176.3:8081', AuthInfo: <nil>}).WithCancel.WithValue(type metadata.mdIncomingKey, val MD{user-agent=[grpc-go/1.64.0], :authority=[server1:8081], grpc-accept-encoding=[gzip], traceparent=[00-95952fc50af45f098f646169bd9ee53c-c655a120963fee03-01], content-type=[application/grpc]}).WithValue(type grpc.serverKey, val <not Stringer>).WithValue(type trace.traceContextKeyType, val <not Stringer>).WithValue(type trace.traceContextKeyType, val <not Stringer>).WithValue(type trace.traceContextKeyType, val <not Stringer>).WithValue(type otelgrpc.gRPCContextKey, val <not Stringer>).WithValue(type grpc.streamKey, val <not Stringer>)

	// get the current span from context
	span := trace.SpanFromContext(ctx)
	// log.Printf("🔍 Server%s | span : %v", *serverID, span)

	spanContext := trace.SpanContextFromContext(ctx)
	// log.Printf("🔍 Server%s | spanContext : %v", *serverID, spanContext)

	traceID := spanContext.TraceID().String()
	spanID := spanContext.SpanID().String()
	log.Printf("🔍 Server%s | Trace ID: %s Span ID: %s", *serverID, traceID, spanID)

	// (!) if the ctx already contains a span the newly created span will be a child of that span
	// (in this case the span contained within the ctx is the automatically created span by otelgrpc)
	// ctx, childSpan := s.tracer.Start(ctx, "server-"+*serverID+"-child-span")
	// defer childSpan.End()
	// log.Printf("🔍 Server%s | childSpan : %v", *serverID, childSpan)

	span.AddEvent(fmt.Sprintf("Server%s Received Request", *serverID))

	p, _ := peer.FromContext(ctx)
	requestFrom := p.Addr.String()
	requestTo := p.LocalAddr.String()
	log.Printf("🟫 Server%s | From: %s To: %s", *serverID, requestFrom, requestTo)

	// grpcMetadata, ok := metadata.FromIncomingContext(ctx)
	// if ok {
	// 	log.Printf("🔍 Server%s | Incoming gRPC Metadata: %v", *serverID, grpcMetadata)
	// }

	span.SetAttributes(
		attribute.String("baz.source", requestFrom),
		attribute.String("baz.destination", requestTo),
	)

	span.AddEvent(fmt.Sprintf("Server%s Started Work On Data", *serverID))
	span.SetAttributes(
		attribute.Int64("baz.data.before", request.Data),
	)

	valueToAdd := rand.Intn(50)
	// adjust the value of data (to emulate some work happening)
	dataAfter := request.Data + int64(valueToAdd)

	// add a delay (to emulate that work taking time)
	delay := time.Duration(rand.Intn(500)) * time.Millisecond
	log.Printf("⏳ Server%s | Emulating work, Adding %d to data, then waiting %v...", *serverID, valueToAdd, delay)
	time.Sleep(delay)

	span.AddEvent(fmt.Sprintf("Server%s Finished Work On Data", *serverID))
	span.SetAttributes(
		attribute.Int64("baz.data.after", dataAfter),
	)

	if s.forwardServer != "" {
		// if there is a forward server, then forward the request to it
		connection, err := grpc.NewClient(
			s.forwardServer,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithStatsHandler(otelgrpc.NewClientHandler()), // (!) for otel
		)
		if err != nil {
			log.Fatalf("grpc client could not connect to the grpc server: %v", err)
		}
		defer connection.Close()

		client := pb.NewMyServiceClient(connection)

		// update the request
		request.Source = serverIP
		request.Destination = s.forwardServer
		request.Data = dataAfter

		log.Printf("⬜ Server%s | Created Request: %v", *serverID, request)

		log.Printf("🟨 Server%s | Forwarding Request to: %s", *serverID, *forwardServerFlag)

		span.AddEvent(fmt.Sprintf("Server%s Forwarding Request", *serverID), trace.WithAttributes(
			attribute.String("baz.forward.server", s.forwardServer),
		))

		response, err := client.MyServiceProcessData(ctx, request)
		if err != nil {
			log.Printf("failed to forward request: %v", err)
			return nil, err
		}

		log.Printf("🟩 Server%s | Received Response: %v", *serverID, response)

		span.AddEvent(fmt.Sprintf("Server%s Received Response", *serverID))

		// span.SetStatus(codes.Ok, "Request processed successfully")

		// update the response
		response.Source = serverIP
		response.Destination = strings.Split(requestFrom, ":")[0]

		log.Printf("⬜ Server%s | Updated Response: %v", *serverID, response)

		log.Printf("🟦 Server%s | Sending Response...", *serverID)

		span.AddEvent(fmt.Sprintf("Server%s Returning Response", *serverID))

		return response, nil
	} else {
		// otherwise respond to the origin caller
		response := &pb.MyServiceResponse{
			Origin:      request.Origin,
			Source:      serverIP,
			Destination: request.Source,
			Data:        dataAfter,
		}
		log.Printf("⬜ Server%s | Created Response: %v", *serverID, response)

		log.Printf("🟦 Server%s | Sending Response...", *serverID)

		span.AddEvent(fmt.Sprintf("Server%s Returning Response", *serverID), trace.WithAttributes(
			attribute.String("baz.return.server", request.Origin),
		))

		return response, nil
	}
}
