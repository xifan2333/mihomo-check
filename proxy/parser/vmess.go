package parser

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

type vmessJson struct {
	V    string      `json:"v"`
	Ps   string      `json:"ps"`
	Add  string      `json:"add"`
	Port interface{} `json:"port"`
	Id   string      `json:"id"`
	Aid  interface{} `json:"aid"`
	Scy  string      `json:"scy"`
	Net  string      `json:"net"`
	Type string      `json:"type"`
	Host string      `json:"host"`
	Path string      `json:"path"`
	Tls  string      `json:"tls"`
	Sni  string      `json:"sni"`
	Alpn string      `json:"alpn"`
	Fp   string      `json:"fp"`
}

func ParseVmess(data string) (map[string]any, error) {
	if !strings.HasPrefix(data, "vmess://") {
		return nil, fmt.Errorf("not vmess format")
	}
	data = data[8:]

	decoded, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, err
	}
	var vmessInfo vmessJson
	if err := json.Unmarshal(decoded, &vmessInfo); err != nil {
		return nil, err
	}

	var port int
	switch v := vmessInfo.Port.(type) {
	case float64:
		port = int(v)
	case string:
		var err error
		port, err = strconv.Atoi(v)
		if err != nil {
			return nil, fmt.Errorf("format error: incorrect port format")
		}
	default:
		return nil, fmt.Errorf("format error: incorrect port format")
	}

	var aid int
	switch v := vmessInfo.Aid.(type) {
	case float64:
		aid = int(v)
	case string:
		aid, err = strconv.Atoi(v)
		if err != nil {
			return nil, fmt.Errorf("format error: alterId format error")
		}
	}

	proxy := map[string]any{
		"name":       vmessInfo.Ps,
		"type":       "vmess",
		"server":     vmessInfo.Add,
		"port":       port,
		"uuid":       vmessInfo.Id,
		"alterId":    aid,
		"cipher":     "auto",
		"network":    vmessInfo.Net,
		"tls":        vmessInfo.Tls == "tls",
		"servername": vmessInfo.Sni,
	}

	switch vmessInfo.Net {
	case "ws":
		wsOpts := map[string]any{
			"path": vmessInfo.Path,
		}
		if vmessInfo.Host != "" {
			wsOpts["headers"] = map[string]any{
				"Host": vmessInfo.Host,
			}
		}
		proxy["ws-opts"] = wsOpts
	case "grpc":
		grpcOpts := map[string]any{
			"serviceName": vmessInfo.Path,
		}
		proxy["grpc-opts"] = grpcOpts
	}

	if vmessInfo.Alpn != "" {
		proxy["alpn"] = strings.Split(vmessInfo.Alpn, ",")
	}

	return proxy, nil
}
