package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
	"time"
)

func TestHttp(t *testing.T) {
	url := "http://175.6.81.115:15944/getVersion"
	type Version struct {
		Code    int    `json:"code"`
		Version string `json:"version"`
	}
	var version Version
	res, err := http.Get(url)
	if err != nil {
		res, _ = http.Get(url)
		log.Println(err)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		time.Sleep(1 * time.Minute)

	}
	err = json.Unmarshal(body, &version)
	log.Println(version.Version)
}
