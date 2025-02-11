package middleware

import (
	"context"
	"crypto/sha256"
	"encoding/binary"
	"ozon-tesk-task/pkg/logger"

	"github.com/99designs/gqlgen/graphql"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

const (
	requestIDHeader = "X-Request-ID"
	userAgentHeader = "User-Agent"
)

func LogMiddleware(log logger.Logger) graphql.OperationMiddleware {
	return func(ctx context.Context, next graphql.OperationHandler) graphql.ResponseHandler {
		opCtx := graphql.GetOperationContext(ctx)

		var (
			req    string
			user   string
			userId int32
		)

		if opCtx != nil {
			req = opCtx.Headers.Get(requestIDHeader)
			if req == "" {
				newUUID, err := uuid.NewUUID()
				if err == nil {
					req = newUUID.String()
				}
			}

			user = opCtx.Headers.Get(userAgentHeader)
			if user == "" {
				user = "unknown"
			}
		}

		ctx = context.WithValue(ctx, logger.RequestID, req)

		if user != "" {
			hash := sha256.Sum256([]byte(user))
			userId = int32(binary.BigEndian.Uint32(hash[:8]))
		}
		ctx = context.WithValue(ctx, "user_id", userId)

		log.Debug(ctx, "request", zap.String("operation name", opCtx.OperationName), zap.Int32("user", userId))

		return next(ctx)
	}
}
