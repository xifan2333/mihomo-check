package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/bestruirui/bestsub/config"
	"github.com/bestruirui/bestsub/proxy"
	"github.com/bestruirui/bestsub/proxy/checker"
	"github.com/bestruirui/bestsub/proxy/info"
	"github.com/bestruirui/bestsub/proxy/saver"
	"github.com/bestruirui/bestsub/utils"
	"github.com/fsnotify/fsnotify"
	"github.com/metacubex/mihomo/log"
	"gopkg.in/yaml.v3"
)

type App struct {
	configPath  string
	interval    int
	watcher     *fsnotify.Watcher
	reloadTimer *time.Timer
}

func NewApp() *App {
	configPath := flag.String("f", "", "config file path")
	flag.Parse()

	return &App{
		configPath: *configPath,
	}
}

func (app *App) Initialize() error {
	if err := app.initConfigPath(); err != nil {
		return fmt.Errorf("init config path failed: %w", err)
	}

	if err := app.loadConfig(); err != nil {
		return fmt.Errorf("load config failed: %w", err)
	}

	if err := app.initConfigWatcher(); err != nil {
		return fmt.Errorf("init config watcher failed: %w", err)
	}
	if config.GlobalConfig.Proxy.Type == "http" {
		utils.LogInfo("use http proxy: %s", config.GlobalConfig.Proxy.Address)
	} else if config.GlobalConfig.Proxy.Type == "socks" {
		utils.LogInfo("use socks proxy: %s", config.GlobalConfig.Proxy.Address)
	} else {
		utils.LogInfo("not use proxy")
	}
	app.interval = config.GlobalConfig.CheckInterval
	log.SetLevel(log.ERROR)
	if config.GlobalConfig.SaveMethod == "http" {
		saver.StartHTTPServer()
	}
	return nil
}

func (app *App) initConfigPath() error {
	if app.configPath == "" {
		execPath := utils.GetExecutablePath()
		configDir := filepath.Join(execPath, "config")

		if err := os.MkdirAll(configDir, 0755); err != nil {
			return fmt.Errorf("create config dir failed: %w", err)
		}

		app.configPath = filepath.Join(configDir, "config.yaml")
	}
	return nil
}

func (app *App) loadConfig() error {
	yamlFile, err := os.ReadFile(app.configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return app.createDefaultConfig()
		}
		return fmt.Errorf("read config file failed: %w", err)
	}

	if err := yaml.Unmarshal(yamlFile, &config.GlobalConfig); err != nil {
		return fmt.Errorf("parse config file failed: %w", err)
	}

	utils.LogInfo("read config file success")
	return nil
}

func (app *App) createDefaultConfig() error {
	utils.LogInfo("config file not found, create default config file")

	if err := os.WriteFile(app.configPath, []byte(config.DefaultConfigTemplate), 0644); err != nil {
		return fmt.Errorf("write default config file failed: %w", err)
	}

	utils.LogInfo("default config file created")
	utils.LogInfo("please edit config file: %v", app.configPath)
	os.Exit(0)
	return nil
}

func (app *App) initConfigWatcher() error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("create file watcher failed: %w", err)
	}

	app.watcher = watcher
	app.reloadTimer = time.NewTimer(0)
	<-app.reloadTimer.C

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					if app.reloadTimer != nil {
						app.reloadTimer.Stop()
					}
					app.reloadTimer.Reset(100 * time.Millisecond)

					go func() {
						<-app.reloadTimer.C
						utils.LogInfo("config file changed, reloading")
						if err := app.loadConfig(); err != nil {
							utils.LogError("reload config file failed: %v", err)
							return
						}
						app.interval = config.GlobalConfig.CheckInterval
					}()
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				utils.LogError("config file watcher error: %v", err)
			}
		}
	}()

	if err := watcher.Add(app.configPath); err != nil {
		return fmt.Errorf("add config file watcher failed: %w", err)
	}

	utils.LogInfo("config file watcher started")
	return nil
}

func (app *App) Run() {
	defer func() {
		app.watcher.Close()
		if app.reloadTimer != nil {
			app.reloadTimer.Stop()
		}
	}()

	utils.LogInfo("progress display: %v", config.GlobalConfig.PrintProgress)

	for {
		maintask()
		nextCheck := time.Now().Add(time.Duration(app.interval) * time.Minute)
		utils.LogInfo("next check time: %v", nextCheck.Format("2006-01-02 15:04:05"))
		time.Sleep(time.Duration(app.interval) * time.Minute)
	}
}

func main() {

	app := NewApp()

	if err := app.Initialize(); err != nil {
		utils.LogError("initialize failed: %v", err)
		os.Exit(1)
	}

	app.Run()
}
func maintask() {

	proxies, err := proxy.GetProxies()
	if err != nil {
	}

	utils.LogInfo("get proxies success: %v", len(proxies))

	proxies = info.DeduplicateProxies(proxies)

	utils.LogInfo("deduplicate proxies: %v", len(proxies))

	proxyTasks := make([]interface{}, len(proxies))
	for i, proxy := range proxies {
		proxyTasks[i] = proxy
	}

	pool := utils.NewThreadPool(config.GlobalConfig.Concurrent, proxyAliveTask)
	pool.Start()
	pool.AddTaskArgs(proxyTasks)
	pool.Wait()
	results := pool.GetResults()
	var success int
	var successProxies []info.Proxy
	for _, result := range results {
		if result.Err != nil {
			continue
		}
		proxy := result.Result.(*info.Proxy)
		if proxy.Info.Alive {
			success++
			proxy.Id = success
			successProxies = append(successProxies, *proxy)
		}
	}
	utils.LogInfo("success proxies: %v", success)

	proxyTasks = make([]interface{}, len(successProxies))
	for i, proxy := range successProxies {
		proxyTasks[i] = proxy
	}
	pool = utils.NewThreadPool(config.GlobalConfig.Concurrent, proxyRenameTask)
	pool.Start()
	pool.AddTaskArgs(proxyTasks)
	pool.Wait()
	results = pool.GetResults()
	var resultProxies []map[string]any
	for _, result := range results {
		if result.Err != nil {
			continue
		}
		proxy := result.Result.(info.Proxy)
		resultProxies = append(resultProxies, proxy.Raw)
	}
	utils.LogInfo("rename end")
	saver.SaveConfig(successProxies)
}
func proxyAliveTask(task interface{}) (interface{}, error) {
	proxy := proxy.NewProxy(task.(map[string]any))
	checker := checker.NewChecker(proxy)
	checker.AliveTest("https://gstatic.com/generate_204", 204)
	for _, item := range config.GlobalConfig.CheckItems {
		switch item {
		case "openai":
			checker.OpenaiTest()
		case "youtube":
			checker.YoutubeTest()
		case "netflix":
			checker.NetflixTest()
		case "disney":
			checker.DisneyTest()
		}
	}
	return proxy, nil
}
func proxyRenameTask(task interface{}) (interface{}, error) {
	proxy := task.(info.Proxy)
	if config.GlobalConfig.RenameMethod == "api" {
		proxy.CountryCodeFromApi()
		proxy.CountryFlag()
		name := fmt.Sprintf("%v %v %03d", proxy.Info.Flag, proxy.Info.Country, proxy.Id)
		proxy.Raw["name"] = name
	} else if config.GlobalConfig.RenameMethod == "match" {
	}
	return proxy, nil
}
