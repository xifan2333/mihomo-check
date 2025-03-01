package info

import (
	"net/http"
)

type Unlock struct {
	Google     bool
	Chatgpt    bool
	Netflix    bool
	Disney     bool
	Youtube    bool
	Cloudflare bool
}

type ProxyInfo struct {
	Unlock  Unlock
	Speed   int
	Delay   uint16
	Alive   bool
	Country string
	Flag    string
}

type Proxy struct {
	Raw    map[string]any
	Id     int
	Client *http.Client
	Info   ProxyInfo
}
