package middleware

import (
	"context"
	"ozon-tesk-task/pkg/logger"

	"github.com/99designs/gqlgen/graphql"
	"go.uber.org/zap"
)

func LogMiddleware(logger logger.Logger) graphql.OperationMiddleware {
	return func(ctx context.Context, next graphql.OperationHandler) graphql.ResponseHandler {
		opCtx := graphql.GetOperationContext(ctx)

		logger.Debug(ctx, "request", zap.String("operation name", opCtx.OperationName))

		return next(ctx)
	}
}
