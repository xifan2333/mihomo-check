package info

import (
	"fmt"
	"net"

	"github.com/bestruirui/bestsub/config"
	"github.com/bestruirui/bestsub/utils"
)

func DeduplicateProxies(proxies []map[string]any) []map[string]any {

	seen := make(map[string]map[string]any)

	deduplicateTasks := make([]interface{}, len(proxies))
	for i, proxy := range proxies {
		deduplicateTasks[i] = proxy
	}

	concurrent := min(len(deduplicateTasks), config.GlobalConfig.Check.Concurrent)

	pool := utils.NewThreadPool(concurrent, deduplicateTask)
	pool.Start()
	pool.AddTaskArgs(deduplicateTasks)
	pool.Wait()
	results := pool.GetResults()

	for _, result := range results {
		if result.Err == nil {
			key := result.Result.(string)
			args, ok := result.Args.(map[string]any)
			if ok {
				if _, exists := seen[key]; !exists {
					seen[key] = args
				}
			} else {
				fmt.Println("Args type assertion failed")
			}
		}
	}

	result := make([]map[string]any, 0, len(seen))
	for _, proxy := range seen {
		result = append(result, proxy)
	}

	return result
}
func deduplicateTask(task interface{}) (interface{}, error) {
	proxy, ok := task.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("task is not a map[string]any")
	}
	server, serverOk := "", false
	if proxy["type"] == "vless" || proxy["type"] == "vmess" {
		if server, serverOk = proxy["servername"].(string); !serverOk {
			server, serverOk = proxy["server"].(string)
		}
	} else {
		server, serverOk = proxy["server"].(string)
	}
	port, portOk := proxy["port"].(int)

	if !serverOk || !portOk {
		return nil, fmt.Errorf("server or port is not a string or int")
	}
	serverip, err := net.LookupIP(server)
	if err != nil {
		return nil, fmt.Errorf("lookup ip failed: %v", err)
	}
	if len(serverip) == 0 {
		return nil, fmt.Errorf("no IP addresses found for server: %s", server)
	}
	key := fmt.Sprintf("%s:%v", serverip[0], port)
	return key, nil
}
