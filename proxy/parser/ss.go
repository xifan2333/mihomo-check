package parser

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

func ParseShadowsocks(data string) (map[string]any, error) {
	if !strings.HasPrefix(data, "ss://") {
		return nil, fmt.Errorf("not ss format")
	}
	data = data[5:]

	if !strings.Contains(data, "@") {
		if strings.Contains(data, "#") {
			temp := strings.SplitN(data, "#", 2)
			data = DecodeBase64(temp[0]) + "#" + temp[1]
		} else {
			data = DecodeBase64(data)
		}
	}
	if !strings.Contains(data, "@") && !strings.Contains(data, "#") {
		return nil, fmt.Errorf("format error: missing @ or # separator")
	}

	name := ""
	if idx := strings.LastIndex(data, "#"); idx != -1 {
		name = data[idx+1:]
		name, _ = url.QueryUnescape(name)
		data = data[:idx]
	}

	parts := strings.SplitN(data, "@", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("format error: missing @ separator")
	}

	parts[0] = DecodeBase64(parts[0])

	methodAndPassword := strings.SplitN(parts[0], ":", 2)
	if len(methodAndPassword) != 2 {
		return nil, fmt.Errorf("format error: incorrect encryption method and password format")
	}

	method := methodAndPassword[0]

	password := DecodeBase64(methodAndPassword[1])

	hostPort := strings.Split(parts[1], ":")
	if len(hostPort) != 2 {
		return nil, fmt.Errorf("format error: incorrect server address format")
	}
	port, err := strconv.Atoi(hostPort[1])
	if err != nil {
		return nil, fmt.Errorf("format error: incorrect port format")
	}

	proxy := map[string]any{
		"name":     name,
		"type":     "ss",
		"server":   hostPort[0],
		"port":     port,
		"cipher":   method,
		"password": password,
	}

	return proxy, nil
}
