package checker

import (
	"io"
	"net/http"
	"time"

	"github.com/bestruirui/bestsub/config"
)

func (c *Checker) CheckSpeed() {
	speedClient := &http.Client{
		Timeout:   time.Duration(config.GlobalConfig.Check.DownloadTimeout) * time.Second,
		Transport: c.Proxy.Client.Transport,
	}

	resp, err := speedClient.Get(config.GlobalConfig.Check.SpeedTestUrl)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	buffer := make([]byte, 32*1024)
	totalBytes := 0
	var startTime time.Time
	firstRead := true

	for {
		n, err := resp.Body.Read(buffer)
		if firstRead && n > 0 {
			startTime = time.Now()
			firstRead = false
		}
		totalBytes += n

		if err != nil {
			if err == io.EOF {
				break
			}
			if totalBytes > 0 {
				break
			}
			return
		}
	}

	if firstRead {
		return
	}

	duration := time.Since(startTime).Milliseconds()
	if duration == 0 {
		duration = 1
	}

	c.Proxy.Info.Speed = int(float64(totalBytes) / 1024 * 1000 / float64(duration))

}
