package parser

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

func ParseHysteria2(data string) (map[string]any, error) {
	if !strings.HasPrefix(data, "hysteria2://") && !strings.HasPrefix(data, "hy2://") {
		return nil, fmt.Errorf("not hysteria2 format")
	}

	link, err := url.Parse(data)
	if err != nil {
		return nil, err
	}

	query := link.Query()
	server := link.Hostname()
	if server == "" {
		return nil, fmt.Errorf("hysteria2 server address error")
	}
	portStr := link.Port()
	if portStr == "" {
		return nil, fmt.Errorf("hysteria2 port error")
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, fmt.Errorf("hysteria2 port error")
	}
	_, obfs, obfsPassword, _, insecure, sni := query.Get("network"), query.Get("obfs"), query.Get("obfs-password"), query.Get("pinSHA256"), query.Get("insecure"), query.Get("sni")
	insecureBool := insecure == "1"

	return map[string]any{
		"type":             "hysteria2",
		"name":             link.Fragment,
		"server":           server,
		"port":             port,
		"ports":            query.Get("mport"),
		"password":         link.User.String(),
		"obfs":             obfs,
		"obfs-password":    obfsPassword,
		"sni":              sni,
		"skip-cert-verify": insecureBool,
		"insecure":         insecure,
		"mport":            query.Get("mport"),
	}, nil
}
