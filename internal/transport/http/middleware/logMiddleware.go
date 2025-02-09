package middleware

import (
	"context"
	"ozon-tesk-task/pkg/logger"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"go.uber.org/zap"
)

func LogMiddleware(logger logger.Logger) graphql.OperationMiddleware {
	return func(ctx context.Context, next graphql.OperationHandler) graphql.ResponseHandler {
		opCtx := graphql.GetOperationContext(ctx)

		start := time.Now()

		logger.Debug(ctx, "request started", zap.String("operation name", opCtx.OperationName))

		responseHandeler := next(ctx)

		dur := time.Since(start)

		logger.Debug(ctx, "request finished", zap.String("operation name", opCtx.OperationName), zap.Any("duration", dur))

		return responseHandeler
	}
}
