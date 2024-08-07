package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/contrib/processors/baggage/baggagetrace"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
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
	"github.com/google/uuid"
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

	log.Printf("🧑 Client%s | IP: %s", *clientID, clientIP)

	// ---------- OTEL SETUP START ---------

	ctx := context.Background()

	grpcTraceClient := otlptracegrpc.NewClient()

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
			semconv.ServiceName("client-"+*clientID+"-service-name"),
		),
	)
	if err != nil {
		log.Fatalf("failed to create otel resource: %v", err)
	}

	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(traceExporter),
		sdktrace.WithResource(resource),
		sdktrace.WithSpanProcessor(baggagetrace.New()), // (!) for passing baggage down
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
	tracer := tracerProvider.Tracer("client-" + *clientID + "-tracer")

	ctx = context.Background()

	// create a new root span
	ctx, rootSpan := tracer.Start(ctx, "client-"+*clientID+"-root-span")
	defer rootSpan.End()
	// log.Printf("🔍 Client%s | rootSpan : %v", *clientID, rootSpan)

	// generate a uuid for the request id to add as baggage
	requestID := uuid.NewString()

	// create a new baggage member with the request_id
	requestIDMember, err := baggage.NewMember("request_id", requestID)
	if err != nil {
		log.Fatalf("failed to create baggage member: %v", err)
	}

	// create a new baggage
	requestIDBaggage, err := baggage.New(requestIDMember)
	if err != nil {
		log.Fatalf("failed to create baggage: %v", err)
	}

	// add the baggage to the context
	ctx = baggage.ContextWithBaggage(ctx, requestIDBaggage)

	// add the baggage as an attribute on the root span
	// (!!!) we have to do this because the Span Processor above doesn't seem to do this on the root span??
	// is this an order of operations issue?
	rootSpan.SetAttributes(attribute.String("request_id", requestID))

	spanContext := trace.SpanContextFromContext(ctx)
	// log.Printf("🔍 Client | spanContext : %v", spanContext)
	traceID := spanContext.TraceID().String()
	spanID := spanContext.SpanID().String()
	log.Printf("🔍 Client%s | Trace ID: %s Span ID: %s", *clientID, traceID, spanID)

	baggageCheck := baggage.FromContext(ctx)
	// log.Printf("🧳 Client%s | Baggage: %v", *clientID, baggageCheck)

	requestIDFromBaggage := baggageCheck.Member("request_id")
	log.Printf("🧳 Client%s | Request ID: %v", *clientID, requestIDFromBaggage.Value())

	connection, err := grpc.NewClient(
		*destination,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()), // (!) for otel,
		// grpc.WithStatsHandler is a gRPC client option that allows you to specify a custom stats handler for the client.
		// the stats handler is responsible for collecting and processing various statistics and metrics related to the gRPC client's operations.
		// otelgrpc.NewClientHandler() is a function provided by the OpenTelemetry gRPC instrumentation library (otelgrpc). It creates a new OpenTelemetry client stats handler.
	)
	if err != nil {
		log.Fatalf("grpc client could not connect to the grpc server: %v", err)
	}
	defer connection.Close()

	client := pb.NewMyServiceClient(connection)

	request := &pb.ProcessDataRequest{
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

	rootSpan.AddEvent(fmt.Sprintf("Client%s Sent Request", *clientID))

	response, err := client.ProcessData(ctx, request)
	if err != nil {
		log.Fatalf("failed to send request and receive response: %v", err)
	}

	rootSpan.AddEvent(fmt.Sprintf("Client%s Received Response", *clientID))

	end := time.Now()
	duration := end.Sub(start)

	log.Printf("🟩 Client%s | Received Response: %v", *clientID, response)
	log.Printf("🟩 Client%s | Total Duration: %v", *clientID, duration)
}
