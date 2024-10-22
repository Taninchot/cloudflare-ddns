package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type DNSRecord struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	Name    string `json:"name"`
	Content string `json:"content"`
}

type CloudflareAPIResponse struct {
	Success bool          `json:"success"`
	Errors  []interface{} `json:"errors"`
	Result  DNSRecord     `json:"result"`
}

func main() {
	//print(getPublicIP())
	apiToken := getEnv("CLOUDFLARE_API_TOKEN")
	zoneId := getEnv("CLOUDFLARE_ZONE_ID")
	recordName := getEnv("CLOUDFLARE_RECORD_NAME")
	checkPublicIpInterval := getEnv("CHECK_PUBLIC_IP_INTERVAL")

	checkInterval, err := time.ParseDuration(checkPublicIpInterval + "ms")
	if err != nil || checkInterval <= 0 {
		fmt.Println("Invalid check interval: " + checkPublicIpInterval)
		os.Exit(1)
	}

	for {
		publicIP := getPublicIP()
		dnsRecord := getCurrentDNSRecord(zoneId, recordName, apiToken)
		if dnsRecord.Content != publicIP {
			fmt.Println("Public IP has changed from " + dnsRecord.Content + " to " + publicIP)
			fmt.Println("Updating DNS record")
			updateDNSRecord(apiToken, zoneId, dnsRecord.ID, recordName, publicIP)
		}
		time.Sleep(checkInterval)
	}

}

func getEnv(key string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		fmt.Println("Environment variable not found: " + key)
		os.Exit(1)
	}
	return value
}

func getPublicIP() string {
	response, err := http.Get("https://ipv4.icanhazip.com")
	if err != nil {
		fmt.Println("Failed to get public IP: " + err.Error())
		os.Exit(1)
	}
	defer closeResponseBody(response.Body)
	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Failed to read public IP response: " + err.Error())
		os.Exit(1)
	}
	return string(body)
}

func getCurrentDNSRecord(zoneId string, recordName string, apiToken string) DNSRecord {
	client := &http.Client{}
	url := fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/dns_records?name=%s", zoneId, recordName)
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		fmt.Println("Failed to create request: " + err.Error())
		os.Exit(1)
	}

	req.Header.Add("Authorization", "Bearer "+apiToken)
	req.Header.Add("Content-Type", "application/json")

	response, err := client.Do(req)
	if err != nil {
		fmt.Println("Failed to get current DNS record: " + err.Error())
		os.Exit(1)
	}
	defer closeResponseBody(response.Body)

	var result struct {
		Success bool          `json:"success"`
		Errors  []interface{} `json:"errors"`
		Result  []DNSRecord   `json:"result"`
	}

	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		fmt.Println("Failed to decode response body: " + err.Error())
		os.Exit(1)
	}
	if !result.Success || len(result.Result) == 0 {
		fmt.Println("Failed to get current DNS record: " + response.Status)
		os.Exit(1)
	}
	return result.Result[0]
}

func updateDNSRecord(apiToken, zoneId, recordId, recordName, newIP string) {
	client := &http.Client{}
	url := fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/dns_records/%s", zoneId, recordId)

	data := map[string]interface{}{
		"type":    "A",
		"name":    recordName,
		"content": newIP,
		"ttl":     0,
		"proxied": false,
	}
	payload, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Failed to marshal data: " + err.Error())
		os.Exit(1)
	}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(payload))
	if err != nil {
		fmt.Println("Failed to create request: " + err.Error())
		os.Exit(1)
	}

	req.Header.Add("Authorization", "Bearer "+apiToken)
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Failed to update DNS record: " + err.Error())
		os.Exit(1)
	}
	defer closeResponseBody(resp.Body)

	var result CloudflareAPIResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		fmt.Println("Failed to decode response body: " + err.Error())
		os.Exit(1)
	}

	if !result.Success {
		fmt.Println("Failed to update DNS record: " + result.Errors[0].(string))
		os.Exit(1)
	}

	fmt.Println("DNS record updated successfully")
}

func closeResponseBody(Body io.ReadCloser) {
	err := Body.Close()
	if err != nil {
		fmt.Println("Failed to close response body: " + err.Error())
		os.Exit(1)
	}
}
