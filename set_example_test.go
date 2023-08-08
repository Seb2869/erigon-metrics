package metrics_test

import (
	"bytes"
	"fmt"
	"github.com/VictoriaMetrics/metrics"
)

func ExampleSet() {
	// Create a set with a counter
	s := metrics.NewSet()
	sc := s.NewCounter("set_counter")
	sc.Inc()
	s.NewGauge(`set_gauge{foo="bar"}`, func() float64 { return 42 })

	// Dump metrics from s.
	var bb bytes.Buffer
	s.WritePrometheus(&bb)
	fmt.Printf("set metrics:\n%s\n", bb.String())

	// Output:
	// set metrics:
	// # HELP set_counter
	// # TYPE set_counter counter
	// set_counter 1
	// # HELP set_gauge
	// # TYPE set_gauge gauge
	// set_gauge{foo="bar"} 42
}
