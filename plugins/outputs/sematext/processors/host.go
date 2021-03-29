package processors

import (
	"bufio"
	"github.com/influxdata/telegraf"
	"os"
	"path"
	"sync"
	"time"
)

const (
	sematextHostTag  = "os.host"
	telegrafHostTag  = "host"
	hostnameFileName = ".resolved-hostname"
)

type Host struct {
	hostname string
	lock     sync.RWMutex
}

// NewHost creates and initializes an instance of Host processor. It also starts periodic host reload goroutine.
func NewHost() *Host {
	// do the initial load before spawning a goroutine which will periodically reload the hostname
	hostnameFileName := getHostnameFileName()
	h := &Host{
		hostname: loadHostname(hostnameFileName),
	}

	// if the Sematext dir (which might hold the hostname file) doesn't exist, no point in starting the goroutine
	if hostnameFileName != "" {
		go func() {
			for {
				time.Sleep(2 * time.Minute)

				h.lock.Lock()
				loadHostname(hostnameFileName)
				h.lock.Unlock()
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

func getHostnameFileName() string {
	if root := GetRootDir(); root != "" {
		return path.Join(root, hostnameFileName)
	}
	return ""
}

func loadHostname(hostnameFile string) string {
	file, err := os.Open(hostnameFile)
	if err != nil {
		return ""
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	hostname := scanner.Text()

	return hostname
}
