package saver

import (
	"fmt"

	"github.com/bestruirui/bestsub/config"
	"github.com/bestruirui/bestsub/proxy/info"
	"github.com/bestruirui/bestsub/utils"
	"gopkg.in/yaml.v3"
)

type ProxyCategory struct {
	Name    string
	Proxies []map[string]any
	Filter  func(result info.Proxy) bool
}

type ConfigSaver struct {
	results    []info.Proxy
	categories []ProxyCategory
	saveMethod func([]byte, string) error
}

func NewConfigSaver(results []info.Proxy) *ConfigSaver {
	return &ConfigSaver{
		results:    results,
		saveMethod: chooseSaveMethod(),
		categories: []ProxyCategory{
			{
				Name:    "all.yaml",
				Proxies: make([]map[string]any, 0),
				Filter:  func(result info.Proxy) bool { return true },
			},
			{
				Name:    "openai.yaml",
				Proxies: make([]map[string]any, 0),
				Filter:  func(result info.Proxy) bool { return result.Info.Unlock.Chatgpt },
			},
			{
				Name:    "youtube.yaml",
				Proxies: make([]map[string]any, 0),
				Filter:  func(result info.Proxy) bool { return result.Info.Unlock.Youtube },
			},
			{
				Name:    "netflix.yaml",
				Proxies: make([]map[string]any, 0),
				Filter:  func(result info.Proxy) bool { return result.Info.Unlock.Netflix },
			},
			{
				Name:    "disney.yaml",
				Proxies: make([]map[string]any, 0),
				Filter:  func(result info.Proxy) bool { return result.Info.Unlock.Disney },
			},
		},
	}
}

func SaveConfig(results []info.Proxy) {
	saver := NewConfigSaver(results)
	if err := saver.Save(); err != nil {
		utils.LogError("save config failed: %v", err)
	}
}

func (cs *ConfigSaver) Save() error {
	cs.categorizeProxies()

	for _, category := range cs.categories {
		if err := cs.saveCategory(category); err != nil {
			utils.LogError("save %s category failed: %v", category.Name, err)
			continue
		}
	}

	return nil
}

func (cs *ConfigSaver) categorizeProxies() {
	for _, result := range cs.results {
		for i := range cs.categories {
			if cs.categories[i].Filter(result) {
				cs.categories[i].Proxies = append(cs.categories[i].Proxies, result.Raw)
			}
		}
	}
}

func (cs *ConfigSaver) saveCategory(category ProxyCategory) error {
	if len(category.Proxies) == 0 {
		utils.LogWarn("%s proxies are empty, skip", category.Name)
		return nil
	}
	yamlData, err := yaml.Marshal(map[string]any{
		"proxies": category.Proxies,
	})
	if err != nil {
		return fmt.Errorf("serialize %s failed: %w", category.Name, err)
	}
	if err := cs.saveMethod(yamlData, category.Name); err != nil {
		return fmt.Errorf("save %s failed: %w", category.Name, err)
	}

	return nil
}

func chooseSaveMethod() func([]byte, string) error {
	switch config.GlobalConfig.Save.Method {
	case "r2":
		if err := ValiR2Config(); err != nil {
			utils.LogError("R2 config is incomplete: %v ,use local save", err)
			return SaveToLocal
		}
		return UploadToR2Storage
	case "gist":
		if err := ValiGistConfig(); err != nil {
			utils.LogError("Gist config is incomplete: %v ,use local save", err)
			return SaveToLocal
		}
		return UploadToGist
	case "webdav":
		if err := ValiWebDAVConfig(); err != nil {
			utils.LogError("WebDAV config is incomplete: %v ,use local save", err)
			return SaveToLocal
		}
		return UploadToWebDAV
	case "local":
		return SaveToLocal
	case "http":
		if err := ValiHTTPConfig(); err != nil {
			utils.LogError("HTTP config is incomplete: %v ,use local save", err)
			return SaveToLocal
		}
		return SaveToHTTP
	default:
		utils.LogError("unknown save method: %s, use local save", config.GlobalConfig.Save.Method)
		return SaveToLocal
	}
}
