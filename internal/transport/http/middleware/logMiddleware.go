package middleware

import (
	"context"
	"ozon-tesk-task/pkg/logger"

	"github.com/99designs/gqlgen/graphql"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

const requestIDHeader = "X-Request-ID"

func LogMiddleware(log logger.Logger) graphql.OperationMiddleware {
	return func(ctx context.Context, next graphql.OperationHandler) graphql.ResponseHandler {
		opCtx := graphql.GetOperationContext(ctx)

		req := opCtx.Headers.Get(requestIDHeader)
		if req == "" {
			newUUID, err := uuid.NewUUID()
			if err == nil {
				req = newUUID.String()
			}
		}

		ctx = context.WithValue(ctx, logger.RequestID, req)

		log.Debug(ctx, "request", zap.String("operation name", opCtx.OperationName))

		return next(ctx)
	}
}
