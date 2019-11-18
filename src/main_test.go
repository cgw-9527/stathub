package main

import (
	"bytes"
	"crypto/md5"
	"crypto/tls"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/likexian/simplejson-go"
)

func TestHttp(t *testing.T) {
	for i := 0; i < 3; i++ {
		err := HttpSend("https://114.67.37.245:15944", "bc97a48eeb86aece4ab7685bab3971e1", "")
		if err != nil {
			log.Println("send stat failed ", err.Error())
			time.Sleep(3 * time.Second)
		} else {
			log.Println("send stat to server successful")
			break
		}
	}

}
func Md4(str, key string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(str+key)))
}

// httpSend send data to stat api
func HttpSend(server, key, stat string) (err error) {
	surl := server + "/api/stat"
	skey := Md4(key, stat)

	request, err := http.NewRequest("POST", surl, bytes.NewBuffer([]byte(stat)))
	if err != nil {
		return
	}

	request.Header.Set("X-Client-Key", skey)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("User-Agent", "Stat Hub API Client/"+"0.102.4"+" (i@likexian.com)")

	tr := &http.Transport{
		// If not self-signed certificate please disabled this.
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{
		Timeout:   time.Duration(30 * time.Second),
		Transport: tr,
	}

	response, err := client.Do(request)
	if err != nil {
		return
	}

	defer response.Body.Close()
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}

	jsonData, err := simplejson.Loads(string(data))
	if err != nil {
		return
	}

	status := jsonData.Get("status.code").MustInt(0)
	if status != 1 {
		message := jsonData.Get("status.message").MustString("unknown error")
		return errors.New("server return: " + message)
	}

	return
}
