package parser

import (
	"strings"
)

func ParseProxy(proxy string) (map[string]any, error) {
	if strings.HasPrefix(proxy, "ss://") {
		return ParseShadowsocks(proxy)
	}
	if strings.HasPrefix(proxy, "trojan://") {
		return ParseTrojan(proxy)
	}
	if strings.HasPrefix(proxy, "vmess://") {
		return ParseVmess(proxy)
	}
	if strings.HasPrefix(proxy, "vless://") {
		return ParseVless(proxy)
	}
	if strings.HasPrefix(proxy, "hysteria2://") {
		return ParseHysteria2(proxy)
	}
	if strings.HasPrefix(proxy, "hy2://") {
		return ParseHysteria2(proxy)
	}
	if strings.HasPrefix(proxy, "ssr://") {
		return ParseSsr(proxy)
	}
	return nil, nil
}
