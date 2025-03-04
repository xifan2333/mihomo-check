package parser

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

func ParseTrojan(data string) (map[string]any, error) {
	if !strings.HasPrefix(data, "trojan://") {
		return nil, fmt.Errorf("not trojan format")
	}

	u, err := url.Parse(data)
	if err != nil {
		return nil, err
	}

	password := u.User.String()
	hostPort := strings.Split(u.Host, ":")
	if len(hostPort) != 2 {
		return nil, nil
	}

	name := ""
	if fragment := u.Fragment; fragment != "" {
		name = fragment
	}

	params := u.Query()
	port, err := strconv.Atoi(hostPort[1])
	if err != nil {
		return nil, fmt.Errorf("format error: incorrect port format")
	}

	proxy := map[string]any{
		"name":     name,
		"type":     "trojan",
		"server":   hostPort[0],
		"port":     port,
		"password": password,
		"network": func() string {
			if t := params.Get("type"); t != "" {
				return t
			} else {
				return "original"
			}
		}(),
		"skip-cert-verify": params.Get("allowInsecure") == "1",
		"allowInsecure":    params.Get("allowInsecure"),
	}

	if params.Get("security") == "tls" {
		proxy["tls"] = true
		if sni := params.Get("sni"); sni != "" {
			proxy["sni"] = sni
		}
	}

	switch params.Get("type") {
	case "ws":
		wsOpts := map[string]any{
			"path": params.Get("path"),
		}
		if host := params.Get("host"); host != "" {
			wsOpts["headers"] = map[string]any{
				"Host": host,
			}
		}
		proxy["ws-opts"] = wsOpts
	case "grpc":
		if serviceName := params.Get("serviceName"); serviceName != "" {
			proxy["grpc-opts"] = map[string]any{
				"serviceName": serviceName,
			}
		}
	}

	return proxy, nil
}
