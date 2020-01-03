package main

import (
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestHttp(t *testing.T) {
	url := "http://175.6.144.117:9879"
	jsonStr := `{"method":"masternode","params":["current"]}`
	req, err := http.NewRequest("POST", url, strings.NewReader(jsonStr))
	if err != nil {
		Nlog("get master node height Post:", err)
	}
	req.Close = true
	req.Header.Set("Content-ype", "text/plain;")
	req.Header.Add("Authorization", "Basic  VWxvcmQwMzpVbG9yZDAz")

	client := &http.Client{
		Timeout: time.Duration(5 * time.Second),
	}
	resp, err := client.Do(req)
	resp = nil
	if resp == nil {
		return
	}
	defer resp.Body.Close()
}
