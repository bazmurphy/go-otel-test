package main

import (
	"context"
	"flag"
	"log"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
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
	destination = flag.String("destination", "", "the address of the grpc server to connect to")
	clientID    = flag.String("id", "", "the client's unique identifier")
	clientIP    = util.GetIPv4Address()
)

func main() {
	flag.Parse()

	if *destination == "" {
		log.Fatalf("'destination' flag required")
	}

	if *clientID == "" {
		log.Fatalf("'id' flag required")
	}

	log.Printf("⬜ Client%s | IP: %s", *clientID, clientIP)

	// ---------- OTEL START ---------

	ctx := context.Background()

	// create an otel trace exporter
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
			semconv.ServiceName("client-"+*clientID+"-service-name"),
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
	defer func() {
		_ = tracerProvider.Shutdown(ctx)
	}()

	// register the tracer provider as the global tracer provider
	otel.SetTracerProvider(tracerProvider)

	// (!!!) this was super important... but I don't understand why the default propagator doesn't automatically work?
	// otel.SetTextMapPropagator(
	// 	propagation.NewCompositeTextMapPropagator(
	// 		propagation.TraceContext{},
	// 		propagation.Baggage{},
	// 	),
	// )

	// sets the TraceContext propagator as the global propagator for the OpenTelemetry SDK
	otel.SetTextMapPropagator(propagation.TraceContext{})

	// create a tracer
	tracer := tracerProvider.Tracer("client-" + *clientID + "-tracer")

	// ---------- OTEL END ---------

	connection, err := grpc.NewClient(
		*destination,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()), // (!) for otel,
		// grpc.WithStatsHandler is a gRPC client option that allows you to specify a custom stats handler for the client. The stats handler is responsible for collecting and processing various statistics and metrics related to the gRPC client's operations.
		// otelgrpc.NewClientHandler() is a function provided by the OpenTelemetry gRPC instrumentation library (otelgrpc). It creates a new OpenTelemetry client stats handler.
	)
	if err != nil {
		log.Fatalf("grpc client could not connect to the grpc server: %v", err)
	}
	defer connection.Close()

	client := pb.NewMyServiceClient(connection)

	// create a new span
	ctx, span := tracer.Start(context.Background(), "client-"+*clientID+"-span-test")
	defer span.End()
	// log.Printf("🔍 Client | span : %v", span)

	spanContext := trace.SpanContextFromContext(ctx)
	// log.Printf("🔍 Client | spanContext : %v", spanContext)
	traceID := spanContext.TraceID().String()
	spanID := spanContext.SpanID().String()
	log.Printf("🔍 Client%s | Trace ID: %s Span ID: %s", *clientID, traceID, spanID)

	request := &pb.MyServiceRequest{
		Origin:      clientIP,
		Source:      clientIP,
		Destination: *destination,
		Data:        0,
	}
	log.Printf("⬜ Client%s | Created Request: %v", *clientID, request)

	// grpcMetadata, _ := metadata.FromOutgoingContext(ctx)
	// log.Printf("🔍 Client | outgoing gRPC metadata: %v", grpcMetadata)

	log.Printf("🟦 Client%s | Sending Request to: %s", *clientID, *destination)

	start := time.Now()

	response, err := client.MyServiceProcessData(ctx, request)
	if err != nil {
		log.Printf("failed to send request: %v", err)
	}

	end := time.Now()
	duration := end.Sub(start)

	log.Printf("🟩 Client%s | Received Response: %v", *clientID, response)
	log.Printf("🟩 Client%s | Total Duration: %v", *clientID, duration)
}
