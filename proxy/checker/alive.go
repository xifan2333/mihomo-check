package checker

import (
	"net/http"
	"time"
)

func (c *Checker) AliveTest(url string, expectedStatus int) {

	defer c.Proxy.CloseTransport()

	start := time.Now()

	req, err := http.NewRequestWithContext(c.Proxy.Ctx, http.MethodHead, url, nil)

	if err != nil {
		return
	}
	resp, err := c.Proxy.Client.Do(req)

	if err != nil {
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode == expectedStatus {
		c.Proxy.Info.Alive = true
		c.Proxy.Info.Delay = uint16(time.Since(start) / time.Millisecond)
	}
}
