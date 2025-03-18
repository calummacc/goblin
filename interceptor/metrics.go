package interceptor

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// MetricsInterceptor implements request metrics collection
type MetricsInterceptor struct {
	BaseInterceptor
	mu      sync.RWMutex
	metrics map[string]*RouteMetrics
}

// RouteMetrics holds metrics for a specific route
type RouteMetrics struct {
	TotalRequests    int64
	TotalErrors      int64
	TotalTimeNs      int64
	AverageTimeNs    int64
	LastRequestTime  time.Time
	LastResponseTime time.Time
}

// NewMetricsInterceptor creates a new metrics interceptor
func NewMetricsInterceptor() *MetricsInterceptor {
	return &MetricsInterceptor{
		metrics: make(map[string]*RouteMetrics),
	}
}

// Before implements Interceptor.Before
func (i *MetricsInterceptor) Before(ctx *ExecutionContext) error {
	// Store start time in context
	ctx.Data["metricsStartTime"] = time.Now()

	i.mu.Lock()
	defer i.mu.Unlock()

	routeKey := fmt.Sprintf("%s %s", ctx.Method, ctx.Path)
	metrics, ok := i.metrics[routeKey]
	if !ok {
		metrics = &RouteMetrics{}
		i.metrics[routeKey] = metrics
	}

	metrics.TotalRequests++
	metrics.LastRequestTime = time.Now()

	return nil
}

// After implements Interceptor.After
func (i *MetricsInterceptor) After(ctx *ExecutionContext) error {
	startTime, ok := ctx.Data["metricsStartTime"].(time.Time)
	if !ok {
		startTime = time.Now()
	}

	duration := time.Since(startTime)

	i.mu.Lock()
	defer i.mu.Unlock()

	routeKey := fmt.Sprintf("%s %s", ctx.Method, ctx.Path)
	metrics := i.metrics[routeKey]

	metrics.TotalTimeNs += duration.Nanoseconds()
	metrics.AverageTimeNs = metrics.TotalTimeNs / metrics.TotalRequests
	metrics.LastResponseTime = time.Now()

	if ctx.GinContext.Writer.Status() >= 400 {
		metrics.TotalErrors++
	}

	return nil
}

// OnRegister implements LifecycleInterceptor.OnRegister
func (i *MetricsInterceptor) OnRegister(ctx context.Context) error {
	fmt.Println("Metrics interceptor registered")
	return nil
}

// OnShutdown implements LifecycleInterceptor.OnShutdown
func (i *MetricsInterceptor) OnShutdown(ctx context.Context) error {
	i.mu.RLock()
	defer i.mu.RUnlock()

	fmt.Println("\nFinal Metrics:")
	for route, metrics := range i.metrics {
		fmt.Printf("\n%s:\n", route)
		fmt.Printf("  Total Requests: %d\n", metrics.TotalRequests)
		fmt.Printf("  Total Errors: %d\n", metrics.TotalErrors)
		fmt.Printf("  Average Response Time: %v\n", time.Duration(metrics.AverageTimeNs))
		fmt.Printf("  Last Request: %v\n", metrics.LastRequestTime)
		fmt.Printf("  Last Response: %v\n", metrics.LastResponseTime)
	}

	return nil
}

// GetMetrics returns a copy of the current metrics
func (i *MetricsInterceptor) GetMetrics() map[string]RouteMetrics {
	i.mu.RLock()
	defer i.mu.RUnlock()

	result := make(map[string]RouteMetrics)
	for k, v := range i.metrics {
		result[k] = *v
	}

	return result
}
