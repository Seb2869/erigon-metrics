package metrics

import (
	"fmt"
	"io"
	"sync/atomic"
)

// NewCounter registers and returns new counter with the given name.
//
// name must be valid Prometheus-compatible metric with possible labels.
// For instance,
//
//   - foo
//   - foo{bar="baz"}
//   - foo{bar="baz",aaa="b"}
//
// The returned counter is safe to use from concurrent goroutines.
func NewCounter(name string, isGauge ...bool) *Counter {
	c := defaultSet.NewCounter(name)
	if len(isGauge) > 0 {
		c.isGauge.Store(isGauge[0])
	}

	return c
}

// Counter is a counter.
//
// It may be used as a gauge if Dec and Set are called.
type Counter struct {
	n       atomic.Uint64
	isGauge atomic.Bool
}

func (c *Counter) IsGauge() bool {
	return c.isGauge.Load()
}

// Inc increments c.
func (c *Counter) Inc() {
	c.n.Add(1)
}

// Dec decrements c.
func (c *Counter) Dec() {
	c.n.Add(^uint64(0))
	c.isGauge.Store(true)
}

// Add adds n to c.
func (c *Counter) Add(n int) {
	if n < 0 {
		c.isGauge.Store(true)
	}
	c.n.Add(uint64(n))
}

// Get returns the current value for c.
func (c *Counter) Get() uint64 {
	return c.n.Load()
}

// Set sets c value to n.
func (c *Counter) Set(n uint64) {
	if n < c.n.Load() {
		c.isGauge.Store(true)
	}

	c.n.Store(n)
}

// marshalTo marshals c with the given prefix to w.
func (c *Counter) marshalTo(prefix string, w io.Writer) {
	v := c.Get()
	fmt.Fprintf(w, "%s %d\n", prefix, v)
}

// GetOrCreateCounter returns registered counter with the given name
// or creates new counter if the registry doesn't contain counter with
// the given name.
//
// name must be valid Prometheus-compatible metric with possible labels.
// For instance,
//
//   - foo
//   - foo{bar="baz"}
//   - foo{bar="baz",aaa="b"}
//
// The returned counter is safe to use from concurrent goroutines.
//
// Performance tip: prefer NewCounter instead of GetOrCreateCounter.
func GetOrCreateCounter(name string, isGauge ...bool) *Counter {
	c := defaultSet.GetOrCreateCounter(name)
	if len(isGauge) > 0 {
		c.isGauge.Store(isGauge[0])
	}
	return c
}
