package datadog_test

import (
	"context"
	"time"

	"github.com/DataDog/sketches-go/ddsketch"
	"go.opentelemetry.io/contrib/exporters/metric/datadog"
	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/api/metric"
	"go.opentelemetry.io/otel/sdk/metric/controller/push"
	"go.opentelemetry.io/otel/sdk/metric/selector/simple"
)

func ExampleExporter() {
	selector := simple.NewWithSketchDistribution(ddsketch.NewDefaultConfig())
	exp, err := datadog.NewExporter(datadog.Options{
		Tags: []string{"env:dev"},
	})
	if err != nil {
		panic(err)
	}
	defer exp.Close()

	pusher := push.New(selector, exp, push.WithPeriod(time.Second*10))
	defer pusher.Stop()
	pusher.Start()
	global.SetMeterProvider(pusher.Provider())
	meter := global.Meter("marwandist")
	m := metric.Must(meter).NewInt64Counter("mycounter")
	meter.RecordBatch(context.Background(), nil, m.Measurement(19))
}
