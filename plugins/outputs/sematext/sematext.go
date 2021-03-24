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
	defaultSematextMetricsReceiverUrl = "https://spm-receiver.sematext.com"
)

// Sematext struct contains configuration read from Telegraf config and a few runtime objects.
// We'll use one separate instance of Telegraf for each monitored service. Therefore, token for particular service
// will be configured on Sematext output level
type Sematext struct {
	ReceiverUrl string          `toml:"receiver_url"`
	Token       string          `toml:"token"`
	ProxyServer string          `toml:"proxy_server"`
	Username    string          `toml:"username"`
	Password    string          `toml:"password"`
	Log         telegraf.Logger `toml:"-"`

	metricsUrl   string
	sender       *sender.Sender
	senderConfig *sender.Config
	serializer   serializer.MetricSerializer
	processors   []processors.Processor
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
	if len(s.ReceiverUrl) == 0 {
		s.ReceiverUrl = defaultSematextMetricsReceiverUrl
	}

	proxyURL, err := url.Parse(s.ProxyServer)
	if err != nil {
		return fmt.Errorf("invalid url %s for the proxy server: %v", s.ProxyServer, err)
	}
	s.senderConfig = &sender.Config{
		ProxyURL: proxyURL,
		Username: s.Username,
		Password: s.Password,
	}
	s.sender = sender.NewSender(s.senderConfig)
	s.metricsUrl = s.ReceiverUrl + "/write?db=metrics"

	s.initProcessors()

	s.serializer = serializer.NewMetricSerializer()

	s.Log.Infof("Sematext output started with Token=%s, ReceiverUrl=%s, ProxyServer=%s", s.Token, s.ReceiverUrl,
		s.ProxyServer)

	return nil
}

// initProcessors instantiates all metric processors that will be used to prepare metrics/tags for sending to Sematext
func (s *Sematext) initProcessors() {
	// add more processors as they are implemented
	s.processors = []processors.Processor{
		&processors.Token{
			Token: s.Token,
		},
	}
}

// Write sends metrics to Sematext backend and handles the response
func (s *Sematext) Write(metrics []telegraf.Metric) error {
	for _, p := range s.processors {
		err := p.Process(metrics)

		if err != nil {
			s.Log.Errorf("can't process metrics in Sematext output, error : %s", err.Error())
			return err
		}
	}

	body := s.serializer.Write(metrics)

	res, err := s.sender.Request("POST", s.metricsUrl, "text/plain; charset=utf-8", body)
	if res != nil && res.Body != nil {
		defer res.Body.Close()
	}
	if err != nil {
		s.Log.Errorf("error while sending to %s : %s", s.ReceiverUrl, err.Error())
		return err
	}

	success := res.StatusCode >= 200 && res.StatusCode < 300
	if !success {
		// TODO in the future, consider handling the retries for bad-request cases
		// badRequest := res.StatusCode >= 400 && res.StatusCode < 500
		// if !badRequest {
		return s.logAndCreateError(res)
	}

	return nil
}

// logAndCreateError logs the error message and forms an error object
func (s *Sematext) logAndCreateError(res *http.Response) error {
	errorMsg := fmt.Sprintf("received %d status code, message = '%s' while sending to %s", res.StatusCode,
		res.Status, s.ReceiverUrl)
	s.Log.Error(errorMsg)
	return fmt.Errorf(errorMsg)
}

func init() {
	outputs.Add("sematext", func() telegraf.Output {
		return &Sematext{}
	})
}
