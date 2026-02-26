package clash

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
)

const (
	defaultClashURL = "http://127.0.0.1:9090"
	proxiesPath     = "/proxies"
)

var (
	mockMode  = os.Getenv("MOCK_CLASH") == "1"
	apiSecret = os.Getenv("MIHOMO_SECRET")
)

type Client struct {
	baseURL       string
	secret        string
	httpClient    *http.Client
	mockProxies   map[string]Proxy
	mockProxiesMu sync.RWMutex
}

type Proxy struct {
	Name    string                 `json:"name"`
	Type    string                 `json:"type"`
	Now     string                 `json:"now"`
	All     []string               `json:"all"`
	History []ProxyHistory         `json:"history"`
	Uptime  string                 `json:"uptime"`
	Extra   map[string]interface{} `json:"extra"`
}

type ProxyHistory struct {
	Time  string `json:"time"`
	Delay int    `json:"delay"`
}

type ProxiesResponse struct {
	Proxies map[string]Proxy `json:"proxies"`
}

func NewClient(baseURL, secret string) *Client {
	if baseURL == "" {
		baseURL = defaultClashURL
	}
	if secret == "" {
		secret = apiSecret
	}
	return &Client{
		baseURL:    baseURL,
		secret:     secret,
		httpClient: &http.Client{},
	}
}

func (c *Client) addAuthHeader(req *http.Request) {
	if c.secret != "" {
		req.Header.Set("Authorization", "Bearer "+c.secret)
	}
}

func (c *Client) GetProxies() (*ProxiesResponse, error) {
	if mockMode {
		c.mockProxiesMu.RLock()
		if c.mockProxies != nil {
			defer c.mockProxiesMu.RUnlock()
			return &ProxiesResponse{Proxies: c.mockProxies}, nil
		}
		c.mockProxiesMu.RUnlock()

		// Initialize mock data for testing
		c.mockProxiesMu.Lock()
		defer c.mockProxiesMu.Unlock()
		c.mockProxies = make(map[string]Proxy)
		c.mockProxies["Proxy Group A"] = Proxy{
			Name: "Proxy Group A",
			Type: "Selector",
			Now:  "Proxy-1",
			All:  []string{"Proxy-1", "Proxy-2", "Proxy-3", "Proxy-4", "Proxy-5", "Proxy-6", "Proxy-7"},
		}
		c.mockProxies["Proxy Group B"] = Proxy{
			Name: "Proxy Group B",
			Type: "URLTest",
			Now:  "Auto-2",
			All:  []string{"Auto-1", "Auto-2", "Auto-3", "Auto-4", "Auto-5", "Auto-6"},
		}
		c.mockProxies["Proxy Group C"] = Proxy{
			Name: "Proxy Group C",
			Type: "Selector",
			Now:  "Direct-1",
			All:  []string{"Direct-1", "Direct-2", "Direct-3", "Direct-4", "Direct-5", "Direct-6", "Direct-7", "Direct-8"},
		}
		return &ProxiesResponse{Proxies: c.mockProxies}, nil
	}

	url := c.baseURL + proxiesPath
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	c.addAuthHeader(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get proxies: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, string(body))
	}

	var result ProxiesResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

func (c *Client) SelectProxy(groupName, proxyName string) error {
	if mockMode {
		c.mockProxiesMu.Lock()
		defer c.mockProxiesMu.Unlock()

		if proxy, ok := c.mockProxies[groupName]; ok {
			for _, p := range proxy.All {
				if p == proxyName {
					proxy.Now = proxyName
					c.mockProxies[groupName] = proxy
					return nil
				}
			}
			return fmt.Errorf("proxy %s not found in group %s", proxyName, groupName)
		}
		return fmt.Errorf("group %s not found", groupName)
	}

	url := c.baseURL + proxiesPath + "/" + groupName

	payload := map[string]string{
		"name": proxyName,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	c.addAuthHeader(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to select proxy: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

func (c *Client) TestDelay(groupName, proxyName string, testURL string) (int, error) {
	url := c.baseURL + proxiesPath + "/" + groupName + "/delay"

	if testURL == "" {
		testURL = "http://www.gstatic.com/generate_204"
	}

	payload := map[string]string{
		"url":     testURL,
		"timeout": "5000",
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return 0, fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequest("GET", url, bytes.NewBuffer(body))
	if err != nil {
		return 0, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	c.addAuthHeader(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("failed to test delay: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return 0, fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, string(body))
	}

	var result map[string]int
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, fmt.Errorf("failed to decode response: %w", err)
	}

	return result["delay"], nil
}
