// Package metrics constructs the metrics that the application will track.
package metrics

import (
	"context"
	"expvar"
)

// This holds the single instance of the metrics value needed for collecting metrics. The expvar
// package is already based on a singleton for different metrics that are registered with the
// package so there isn't much choice here.
var m *metrics

// Metrics represents the set of metrics we gather. These fields are safe to be accessed
// concurrently thanks to expvar. No extra abstraction is required.
type metrics struct {
	goroutines *expvar.Int
	requests   *expvar.Int
	errors     *expvar.Int
	panics     *expvar.Int
}

// init constructs the metrics value that will be used to capture metrics. The metrics value is
// stored inside a package level variable since everything inside expvar is registered as a
// singleton. The use of once will make sure this initialization happens only once.
func init() {
	m = &metrics{
		goroutines: expvar.NewInt("goroutines"),
		requests:   expvar.NewInt("requests"),
		errors:     expvar.NewInt("errors"),
		panics:     expvar.NewInt("panics"),
	}
}

// Metrics gathering will be supported through the context.
// ctxKeyMetrics represents the type of value for the context key.
type ctxKeyMetrics int

// metricsKey is how metric values are stored and retrieved.
const metricsKey ctxKeyMetrics = 1

// Set sets the metrics data in to the context.
func Set(ctx context.Context) context.Context {
	return context.WithValue(ctx, metricsKey, m)
}

// AddGoroutine increments the goroutines metric by 1.
func AddGoroutine(ctx context.Context) {
	if v, ok := ctx.Value(metricsKey).(*metrics); ok {
		if v.requests.Value()%100 == 0 {
			v.goroutines.Add(1)
		}
	}
}

// AddRequest increments the requests metric by 1.
func AddRequest(ctx context.Context) {
	if v, ok := ctx.Value(metricsKey).(*metrics); ok {
		v.requests.Add(1)
	}
}

// AddError increments the errors metric by 1.
func AddError(ctx context.Context) {
	if v, ok := ctx.Value(metricsKey).(*metrics); ok {
		v.errors.Add(1)
	}
}

// AddPanic increments the panics metric by 1.
func AddPanic(ctx context.Context) {
	if v, ok := ctx.Value(metricsKey).(*metrics); ok {
		v.panics.Add(1)
	}
}
