package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/bestruirui/bestsub/config"
)

type httpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type versionResponse struct {
	Version string `json:"version"`
}

type providersResponse struct {
	Providers map[string]struct {
		VehicleType string `json:"vehicleType"`
	} `json:"providers"`
}

func makeRequest(client httpClient, method, url string) ([]byte, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request failed: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", config.GlobalConfig.MihomoApiSecret))

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("execute request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNoContent {
			return nil, nil
		}
		return nil, fmt.Errorf("API returned non-200 status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body failed: %w", err)
	}

	return body, nil
}

func UpdateSubs() {
	if config.GlobalConfig.MihomoApiUrl == "" {
		LogWarn("MihomoApiUrl not configured, skipping update")
		return
	}

	version, err := getVersion(http.DefaultClient)
	if err != nil {
		LogError("get version failed: %v", err)
		return
	}

	LogInfo("Mihomo version: %s", version)

	names, err := getNeedUpdateNames(http.DefaultClient)
	if err != nil {
		LogError("get need update subs failed: %v", err)
		return
	}

	if err := updateSubs(http.DefaultClient, names); err != nil {
		LogError("update subs failed: %v", err)
		return
	}
	LogInfo("subs updated")
}

func getVersion(client httpClient) (string, error) {
	url := fmt.Sprintf("%s/version", config.GlobalConfig.MihomoApiUrl)
	body, err := makeRequest(client, http.MethodGet, url)
	if err != nil {
		return "", err
	}

	var version versionResponse
	if err := json.Unmarshal(body, &version); err != nil {
		return "", fmt.Errorf("parse version info failed: %w", err)
	}
	return version.Version, nil
}

func getNeedUpdateNames(client httpClient) ([]string, error) {
	url := fmt.Sprintf("%s/providers/proxies", config.GlobalConfig.MihomoApiUrl)
	body, err := makeRequest(client, http.MethodGet, url)
	if err != nil {
		return nil, err
	}

	var response providersResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("parse provider info failed: %w", err)
	}

	var names []string
	for name, provider := range response.Providers {
		if provider.VehicleType == "HTTP" {
			names = append(names, name)
		}
	}
	return names, nil
}

func updateSubs(client httpClient, names []string) error {
	for _, name := range names {
		url := fmt.Sprintf("%s/providers/proxies/%s", config.GlobalConfig.MihomoApiUrl, name)
		if _, err := makeRequest(client, http.MethodPut, url); err != nil {
			LogError("update sub %s failed: %v", name, err)
		}
		LogInfo("update sub %s success", name)
	}
	return nil
}
