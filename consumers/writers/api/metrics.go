// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package api

import (
	"time"

	"git.willowglen.ca/sq/third-party/mainflux/consumers"
	"github.com/go-kit/kit/metrics"
)

var _ consumers.Consumer = (*metricsMiddleware)(nil)

type metricsMiddleware struct {
	counter  metrics.Counter
	latency  metrics.Histogram
	consumer consumers.Consumer
}

// MetricsMiddleware returns new message repository
// with Save method wrapped to expose metrics.
func MetricsMiddleware(consumer consumers.Consumer, counter metrics.Counter, latency metrics.Histogram) consumers.Consumer {
	return &metricsMiddleware{
		counter:  counter,
		latency:  latency,
		consumer: consumer,
	}
}

func (mm *metricsMiddleware) Consume(msgs interface{}) error {
	defer func(begin time.Time) {
		mm.counter.With("method", "consume").Add(1)
		mm.latency.With("method", "consume").Observe(time.Since(begin).Seconds())
	}(time.Now())
	return mm.consumer.Consume(msgs)
}
