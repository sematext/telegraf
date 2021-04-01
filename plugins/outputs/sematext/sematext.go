package sematext

import (
	"fmt"
	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/plugins/outputs"
	"github.com/influxdata/telegraf/plugins/outputs/sematext/processors"
	"github.com/influxdata/telegraf/plugins/outputs/sematext/sender"
	"github.com/influxdata/telegraf/plugins/outputs/sematext/serializer"
	"net/http"
	"net/url"
)

const (
	defaultSematextMetricsReceiverURL = "https://spm-receiver.sematext.com"
)

// Sematext struct contains configuration read from Telegraf config and a few runtime objects.
// We'll use one separate instance of Telegraf for each monitored service. Therefore, token for particular service
// will be configured on Sematext output level
type Sematext struct {
	ReceiverURL string          `toml:"receiver_url"`
	Token       string          `toml:"token"`
	ProxyServer string          `toml:"proxy_server"`
	Username    string          `toml:"username"`
	Password    string          `toml:"password"`
	Log         telegraf.Logger `toml:"-"`

	metricsURL       string
	sender           *sender.Sender
	senderConfig     *sender.Config
	serializer       serializer.MetricSerializer
	metricProcessors []processors.MetricProcessor
	batchProcessors  []processors.BatchProcessor
}

// TODO add real sample config
const sampleConfig = `
  ## Sample config for Sematext output
`

// Connect is no-op for Sematext output plugin, everything was set up before in Init() method
func (s *Sematext) Connect() error {
	return nil
}

// Close Closes the Sematext output
func (s *Sematext) Close() error {
	s.sender.Close()

	for _, mp := range s.metricProcessors {
		mp.Close()
	}

	for _, bp := range s.batchProcessors {
		bp.Close()
	}

	return nil
}

// SampleConfig Returns a sample configuration for the Sematext output
func (s *Sematext) SampleConfig() string {
	return sampleConfig
}

// Description returns the description for the Sematext output
func (s *Sematext) Description() string {
	return "Use telegraf to send metrics to Sematext"
}

// Init performs full initialization of Sematext output
func (s *Sematext) Init() error {
	if len(s.Token) == 0 {
		return fmt.Errorf("'token' is a required field for Sematext output")
	}
	if len(s.ReceiverURL) == 0 {
		s.ReceiverURL = defaultSematextMetricsReceiverURL
	}

	var proxyURL *url.URL = nil

	if s.ProxyServer != "" {
		var err error
		proxyURL, err = url.Parse(s.ProxyServer)
		if err != nil {
			return fmt.Errorf("invalid url %s for the proxy server: %v", s.ProxyServer, err)
		}
	}

	s.senderConfig = &sender.Config{
		ProxyURL: proxyURL,
		Username: s.Username,
		Password: s.Password,
	}
	s.sender = sender.NewSender(s.senderConfig)
	s.metricsURL = s.ReceiverURL + "/write?db=metrics"

	s.initProcessors()

	s.serializer = serializer.NewMetricSerializer(s.Log)

	s.Log.Infof("Sematext output started with Token=%s, ReceiverUrl=%s, ProxyServer=%s", s.Token, s.ReceiverURL,
		s.ProxyServer)

	return nil
}

// initProcessors instantiates all metric processors that will be used to prepare metrics/tags for sending to Sematext
func (s *Sematext) initProcessors() {
	// add more processors as they are implemented
	s.metricProcessors = []processors.MetricProcessor{
		processors.NewToken(s.Token),
		processors.NewHost(s.Log),
	}
	s.batchProcessors = []processors.BatchProcessor{
		processors.NewHeartbeat(),
		processors.NewMetainfo(s.Log, s.Token, s.ReceiverURL, s.senderConfig),
	}
}

// Write sends metrics to Sematext backend and handles the response
func (s *Sematext) Write(metrics []telegraf.Metric) error {
	processedMetrics, err := s.processMetrics(metrics)

	if err != nil {
		// error means the whole batch should be discarded without sending it. To achieve that, we have to return
		// nil
		s.Log.Errorf("error while preparing to send metrics to Sematext, the batch will be dropped: %v", err)
		return nil
	}

	if len(processedMetrics) > 0 {
		body := s.serializer.Write(processedMetrics)

		res, err := s.sender.Request("POST", s.metricsURL, "text/plain; charset=utf-8", body)
		if err != nil {
			// TODO whether we return an error or not should depend on whether there should be a retry
			s.Log.Errorf("error while sending to %s : %s", s.ReceiverURL, err.Error())
			return err
		}
		defer res.Body.Close()

		success := res.StatusCode >= 200 && res.StatusCode < 300
		if !success {
			// TODO in the future, consider handling the retries for bad-request cases
			// badRequest := res.StatusCode >= 400 && res.StatusCode < 500
			// if !badRequest {
			return s.logAndCreateError(res)
		}
	}

	return nil
}

// processMetrics returns an error only when the whole batch of metrics should be discarded
func (s *Sematext) processMetrics(metrics []telegraf.Metric) ([]telegraf.Metric, error) {
	for _, p := range s.batchProcessors {
		var err error
		metrics, err = p.Process(metrics)

		if err != nil {
			s.Log.Errorf("error while running batch processors in Sematext output: %v", err)
			return metrics, err
		}
	}

	processedMetrics := make([]telegraf.Metric, 0, len(metrics))

	for _, metric := range metrics {
		metricOk := true
		for _, p := range s.metricProcessors {
			err := p.Process(metric)

			if err != nil {
				// log the message, mark the metric to be skipped, skip other processors
				s.Log.Warnf("can't process metric: %s in Sematext output, error : %s", metric, err.Error())
				metricOk = false
				break
			}
		}
		if metricOk {
			processedMetrics = append(processedMetrics, metric)
		}
	}
	return processedMetrics, nil
}

// TODO may not be needed as we have to rework how retry logic works depending on response status codes; sometimes
// we'll log the message, sometimes return an error, possibly never have to do both
// logAndCreateError logs the error message and forms an error object
func (s *Sematext) logAndCreateError(res *http.Response) error {
	errorMsg := fmt.Sprintf("received %d status code, message = '%s' while sending to %s", res.StatusCode,
		res.Status, s.ReceiverURL)
	s.Log.Error(errorMsg)
	return fmt.Errorf(errorMsg)
}

func init() {
	outputs.Add("sematext", func() telegraf.Output {
		return &Sematext{}
	})
}
