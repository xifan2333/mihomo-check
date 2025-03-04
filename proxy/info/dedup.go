package info

import (
	"fmt"
	"net"
	"sync"

	"github.com/bestruirui/bestsub/config"
	"github.com/panjf2000/ants/v2"
)

var (
	dedupProxies = make(map[string]map[string]any)
	dedupMutex   sync.Mutex
)

func addDedupProxy(key string, arg map[string]any) {
	dedupMutex.Lock()
	defer dedupMutex.Unlock()
	if _, exists := dedupProxies[key]; !exists {
		dedupProxies[key] = arg
	}
}

func DeduplicateProxies(proxies []map[string]any) []map[string]any {
	var wg sync.WaitGroup

	pool, _ := ants.NewPool(config.GlobalConfig.Check.Concurrent)
	defer pool.Release()

	for _, proxy := range proxies {
		wg.Add(1)
		pool.Submit(func() {
			defer wg.Done()
			deduplicateTask(proxy)
		})
	}
	wg.Wait()

	result := make([]map[string]any, 0, len(dedupProxies))
	for _, proxy := range dedupProxies {
		result = append(result, proxy)
	}

	return result
}

func deduplicateTask(arg map[string]any) {

	server, serverOk := "", false
	if arg["type"] == "vless" || arg["type"] == "vmess" {
		server, serverOk = arg["servername"].(string)
		if !serverOk || server == "" {
			server, serverOk = arg["server"].(string)
		}
	} else {
		server, serverOk = arg["server"].(string)
	}
	port, portOk := arg["port"].(int)

	if !serverOk || !portOk {
		return
	}
	serverip, err := net.LookupIP(server)
	if err != nil {
		return
	}
	if len(serverip) == 0 {
		return
	}
	key := fmt.Sprintf("%s:%v", serverip[0], port)

	addDedupProxy(key, arg)
}
