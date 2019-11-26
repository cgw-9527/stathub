package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
)

func TestHttp(t *testing.T) {
	var chainStatus ChainStatus
	url := "http://explorer.ulord.one/api/status"
	response, err := http.Get(url)
	if err != nil {
		log.Println(err)
		response, _ = http.Get(url)
	}

	defer response.Body.Close()
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
		return
	}
	err = json.Unmarshal(data, &chainStatus)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(chainStatus.Info.Blocks)

}
