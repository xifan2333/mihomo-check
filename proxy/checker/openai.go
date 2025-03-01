package checker

import (
	"io"
	"strings"
)

func (c *Checker) OpenaiTest() {
	resp, err := c.Proxy.Client.Get("https://android.chat.openai.com")
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
