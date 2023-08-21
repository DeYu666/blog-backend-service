package blog

import (
	"context"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func InjectLogger(ctx context.Context, l *zap.Logger) context.Context {

	if ginCtx, ok := ctx.(*gin.Context); ok {
		ginCtx.Set("LoggerKey", l)
		return ginCtx
	}

	ctx = context.WithValue(ctx, LoggerKey{}, l)
	return ctx
}

// Extract returns the logger stored in the context, or a no-op logger if none
// is available.
// context has two type: gin.Context and context.Context
func Extract(ctx context.Context) *zap.Logger {

	if ginCtx, ok := ctx.(*gin.Context); ok {
		return extractFromGin(ginCtx)
	}

	if l, ok := ctx.Value(LoggerKey{}).(*zap.Logger); ok {
		return l
	}

	if ginCtx, ok := ctx.Value("ginCtx").(*gin.Context); ok {
		return extractFromGin(ginCtx)
	}

	return zap.NewNop()
}

func extractFromGin(ctx *gin.Context) *zap.Logger {

	l, ok := ctx.Get("LoggerKey")
	if !ok {
		return nil
	}

	return l.(*zap.Logger)
}
