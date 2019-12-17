/*
 * Copyright 2015-2019 Li Kexian
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * A smart Hub for holding server stat
 * https://www.likexian.com/
 */

package main

import (
	"log"
	"time"

	"github.com/likexian/simplejson-go"
)

// StatService statSend loop
func StatService() {
	log.Println("start stat service")
	for {
		statSend()
		time.Sleep(300 * time.Second)
	}
}

// statSend get host stat and send to server
func statSend() {
	stat := GetStat(SERVER_CONFIG.Id, SERVER_CONFIG.Name)
	//将结构体数据转为string
	data := simplejson.New(stat)
	result, err := data.Dumps()
	if err != nil {
		Nlog("statsend struct=>string: %v", err.Error())
		return
	}

	log.Println("get stat data: %v", result)
	for i := 0; i < 3; i++ {
		err := httpSend(SERVER_CONFIG.ServerUrl, result)
		if err != nil {
			Nlog("send stat failed, %v", err.Error())
			time.Sleep(3 * time.Second)
		} else {
			log.Println("send stat to server successful")
			break
		}
	}

}
