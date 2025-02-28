package proxies

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/bestruirui/mihomo-check/config"
	"github.com/bestruirui/mihomo-check/proxy/parser"
	"github.com/bestruirui/mihomo-check/utils"
	"github.com/metacubex/mihomo/log"
	"gopkg.in/yaml.v3"
)

func GetProxies() ([]map[string]any, error) {
	log.Infoln("当前共设置了%d个订阅链接", len(config.GlobalConfig.SubUrls))

	subUrls := make([]interface{}, len(config.GlobalConfig.SubUrls))
	for i, url := range config.GlobalConfig.SubUrls {
		subUrls[i] = url
	}
	// 根据 len(subUrls) 和 config.GlobalConfig.PrintProgress 计算出需要多少个线程
	numWorkers := min(len(subUrls), config.GlobalConfig.Concurrent)

	pool := utils.NewThreadPool(numWorkers, taskGetProxies)
	pool.Start()
	pool.AddTaskArgs(subUrls)
	pool.Wait()
	results := pool.GetResults()
	var mihomoProxies []map[string]any

	for _, result := range results {
		if result.Result != nil {
			mihomoProxies = append(mihomoProxies, result.Result.([]map[string]any)...)
		}
	}
	return mihomoProxies, nil
}

func taskGetProxies(args interface{}) (interface{}, error) {
	subUrl := args.(string)
	var mihomoProxies []map[string]any
	data, err := GetDateFromSubs(subUrl)
	if err != nil {
		return nil, err
	}
	var config map[string]any
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		reg, _ := regexp.Compile("(ssr|ss|vmess|trojan|vless|hysteria|hy2|hysteria2)://")
		// 如果不匹配则base64解码
		if !reg.Match(data) {
			data = []byte(parser.DecodeBase64(string(data)))
		}
		if reg.Match(data) {
			proxies := strings.Split(string(data), "\n")

			for _, proxy := range proxies {
				parseProxy, err := ParseProxy(proxy)
				if err != nil {
					continue
				}
				// 如果proxy为空，则跳过
				if parseProxy == nil {
					continue
				}
				mihomoProxies = append(mihomoProxies, parseProxy)
			}
		}
	}
	proxyInterface, ok := config["proxies"]
	if !ok || proxyInterface == nil {
		log.Errorln("订阅链接: %s 没有proxies", subUrl)
		return nil, fmt.Errorf("订阅链接: %s 没有proxies", subUrl)
	}

	proxyList, ok := proxyInterface.([]any)
	if !ok {
		return nil, fmt.Errorf("订阅链接: %s 没有proxies", subUrl)
	}

	for _, proxy := range proxyList {
		proxyMap, ok := proxy.(map[string]any)
		if !ok {
			continue
		}
		mihomoProxies = append(mihomoProxies, proxyMap)
	}
	return mihomoProxies, nil
}

// 订阅链接中获取数据
func GetDateFromSubs(subUrl string) ([]byte, error) {
	maxRetries := config.GlobalConfig.SubUrlsReTry
	var lastErr error

	client := utils.NewHTTPClient()

	for i := 0; i < maxRetries; i++ {
		if i > 0 {
			time.Sleep(time.Second)
		}

		req, err := http.NewRequest("GET", subUrl, nil)
		if err != nil {
			lastErr = err
			continue
		}

		req.Header.Set("User-Agent", "clash.meta")

		resp, err := client.Do(req)
		if err != nil {
			lastErr = err
			continue
		}
		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			lastErr = fmt.Errorf("订阅链接: %s 返回状态码: %d", subUrl, resp.StatusCode)
			continue
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			lastErr = err
			continue
		}
		return body, nil
	}

	return nil, fmt.Errorf("重试%d次后失败: %v", maxRetries, lastErr)
}
