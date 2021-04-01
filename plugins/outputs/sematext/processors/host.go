package processors

import (
	"io/ioutil"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/influxdata/telegraf"
)

const (
	sematextHostTag  = "os.host"
	telegrafHostTag  = "host"
	hostnameFileName = ".resolved-hostname"
)

type Host struct {
	hostname     string
	lock         sync.RWMutex
	reloadTicker *time.Ticker
	stopReload   chan bool
	Log          telegraf.Logger
}

// NewHost creates and initializes an instance of Host processor. It also starts periodic host reload goroutine.
func NewHost(log telegraf.Logger) MetricProcessor {
	// do the initial load before spawning a goroutine which will periodically reload the hostname
	hostnameFileName := getHostnameFileName()

	var host string
	var err error
	if hostnameFileName != "" {
		host, err = loadHostname(hostnameFileName)
		if err != nil {
			log.Warnf("can't load the hostname from the file %s, error: %v", hostnameFileName, err)
		}
	}

	h := &Host{
		hostname: host,
		Log:      log,
	}

	// if the Sematext dir (which might hold the hostname file) doesn't exist, no point in starting the ticker
	if hostnameFileName != "" {
		h.reloadTicker = time.NewTicker(5 * time.Minute)
		h.stopReload = make(chan bool, 1)

		go func() {
			for {
				select {
				case <-h.stopReload:
					return
				case <-h.reloadTicker.C:
					host, err = loadHostname(hostnameFileName)
					if err != nil {
						log.Warnf("can't load the hostname from the file %s, error: %v", hostnameFileName, err)
					}

					h.lock.Lock()
					h.hostname = host
					h.lock.Unlock()
				}
			}
		}()
	}

	return h
}

// Process adjusts the host tag to be compliant with Sematext backend
func (h *Host) Process(metric telegraf.Metric) error {
	// locking because of h.hostname which might be written to by a separate goroutine
	h.lock.RLock()
	defer h.lock.RUnlock()

	adjustHostname(metric, h.hostname)

	return nil
}

// Close clears the resources processor used
func (h *Host) Close() {
	if h.stopReload != nil {
		h.stopReload <- true
	}
}

func adjustHostname(metric telegraf.Metric, loadedHostname string) {
	if loadedHostname != "" {
		metric.RemoveTag(telegrafHostTag)
		metric.AddTag(sematextHostTag, loadedHostname)
	} else {
		h, set := metric.GetTag(telegrafHostTag)
		if set {
			metric.RemoveTag(telegrafHostTag)
			metric.AddTag(sematextHostTag, h)
		}
	}
}

// getHostnameFileName returns the full path of the hostname
func getHostnameFileName() string {
	if root := GetRootDir(); root != "" {
		return path.Join(root, hostnameFileName)
	}
	return ""
}

func loadHostname(hostnameFileName string) (string, error) {
	data, err := ioutil.ReadFile(hostnameFileName)
	if err != nil {
		return "", err
	}

	fullStr := string(data)
	return strings.Split(fullStr, "\n")[0], nil
}
