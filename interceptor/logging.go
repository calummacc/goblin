package interceptor

import (
	"fmt"
	"time"
)

// LoggingInterceptor implements request/response logging
type LoggingInterceptor struct {
	BaseInterceptor
}

// NewLoggingInterceptor creates a new logging interceptor
func NewLoggingInterceptor() *LoggingInterceptor {
	return &LoggingInterceptor{}
}

// Before implements Interceptor.Before
func (i *LoggingInterceptor) Before(ctx *ExecutionContext) error {
	// Store start time in context
	ctx.Data["startTime"] = time.Now()

	fmt.Printf("[%s] -> %s %s\n",
		time.Now().Format(time.RFC3339),
		ctx.Method,
		ctx.Path,
	)

	return nil
}

// After implements Interceptor.After
func (i *LoggingInterceptor) After(ctx *ExecutionContext) error {
	startTime, ok := ctx.Data["startTime"].(time.Time)
	if !ok {
		startTime = time.Now()
	}

	duration := time.Since(startTime)

	fmt.Printf("[%s] <- %s %s (took %v)\n",
		time.Now().Format(time.RFC3339),
		ctx.Method,
		ctx.Path,
		duration,
	)

	return nil
}
