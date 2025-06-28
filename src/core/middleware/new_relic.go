package middleware

import (
	"github.com/KhaiHust/email-notification-service/core/properties"
	"github.com/golibs-starter/golib"
	"github.com/newrelic/go-agent/v3/newrelic"
	"go.uber.org/fx"
)

func NewRelicOpt() fx.Option {
	return fx.Options(
		golib.ProvideProps(properties.NewNewrelicProperties),
		fx.Provide(NewRelicApp),
	)
}

func NewRelicApp(app *golib.App, props *properties.NewrelicProperties) *newrelic.Application {
	nrApp, err := newrelic.NewApplication(
		newrelic.ConfigAppName(props.AppName),
		newrelic.ConfigLicense(props.LicenseKey),
		newrelic.ConfigAppLogForwardingEnabled(true),
		newrelic.ConfigEnabled(props.Enabled),
		//newrelic.ConfigDebugLogger(os.Stdout),
		func(cfg *newrelic.Config) {
			cfg.ErrorCollector.RecordPanics = true
		},
	)
	if err != nil {
		panic("Could connect connect to newrelic: " + err.Error())
	}
	return nrApp
}

func NewSegment(name string, transaction *newrelic.Transaction) *newrelic.Segment {
	if transaction == nil {
		return nil
	}
	segment := newrelic.Segment{}
	segment.Name = name
	segment.StartTime = transaction.StartSegmentNow()
	return &segment
}

func NewDatastoreSegment(name string,
	product newrelic.DatastoreProduct,
	host,
	port,
	database,
	operation string,
	transaction *newrelic.Transaction) *newrelic.DatastoreSegment {
	if transaction == nil {
		return nil
	}
	segment := &newrelic.DatastoreSegment{
		Product:      product,
		Operation:    operation,
		Host:         host,
		PortPathOrID: port,
		DatabaseName: database,
	}
	segment.StartTime = transaction.StartSegmentNow()
	return segment
}

func EndSegment(segment *newrelic.Segment) {
	if segment == nil {
		return
	}
	segment.End()
}
