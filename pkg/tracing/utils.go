package tracing

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/metadata"
)

func StartGrpcServerTracerSpan(ctx context.Context, operationName string) (context.Context, trace.Span) {
	ctx = InjectTextMapCarrierToGrpcMetaData(ctx)
	tracer := otel.GetTracerProvider().Tracer("account_service_tracer")
	ctx, span := tracer.Start(ctx, operationName)
	ctx = InjectTextMapCarrierToGrpcMetaData(ctx)

	return ctx, span
}

func InjectTextMapCarrierToGrpcMetaData(ctx context.Context) context.Context {
	textMapCarrier := make(propagation.MapCarrier)
	otel.GetTextMapPropagator().Inject(ctx, textMapCarrier)
	md := metadata.New(textMapCarrier)
	md2, _ := metadata.FromIncomingContext(ctx)
	for key, values := range md2 {
		md[key] = append(md[key], values...)
	}
	ctx = metadata.NewOutgoingContext(ctx, md)

	return ctx
}
