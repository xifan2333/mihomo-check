package checker

func (c *Checker) GoogleTest() {
	resp, err := c.Proxy.Client.Get("http://www.google.com/generate_204")
	if err != nil {
	}
	defer resp.Body.Close()
	if resp.StatusCode == 204 {
		c.Proxy.Info.Unlock.Google = true
	}
}
