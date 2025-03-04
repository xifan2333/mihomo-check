package checker

import (
	"io"
	"net/http"
	"strings"
)

func (c *Checker) OpenaiTest() {
	req, err := http.NewRequestWithContext(c.Proxy.Ctx, "GET", "https://android.chat.openai.com", nil)
	if err != nil {
		return
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")

	resp, err := c.Proxy.Client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == 403 {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return
		}
		if strings.Contains(string(body), "Request is not allowed. Please try again later.") {
			c.Proxy.Info.Unlock.Chatgpt = true
		}
	}
}
