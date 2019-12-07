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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"strings"
	"time"

	hoststat "github.com/likexian/host-stat-go"
)

//Check master status
func checkStatus() {
	var produce Produce
	//If the master node is stopped, start it
	str := "ps aux|grep ulordd|grep -v grep"
	cmd := exec.Command("sh", "-c", str)
	out1, err := cmd.CombinedOutput()
	if err != nil {
		log.Println(err)
	}
	if string(out1) == "" {
		cmd = exec.Command("sh", "-c", "ulordd")
		go func() { cmd.Run() }()
		time.Sleep(30 * time.Minute)
	}
	//restart master node
	for {
		// log.Println("current--------------------")
		cmd := exec.Command("ulord-cli", "masternode", "current")
		out1, err := cmd.CombinedOutput()
		if err != nil {
			log.Println(err)
		}
		err = json.Unmarshal(out1, &produce)
		if err != nil {
			log.Println(err)
		}
		// log.Println("getChainHeight():", getChainHeight(), "produce.Height:", produce.Height)
		//If the current machine falls behind 6 blocks on the chain, restart the machine
		if getChainHeight()-produce.Height > 25 {
			cmd := exec.Command("ulord-cli", "stop")
			stop, _ := cmd.CombinedOutput()
			log.Println("stop:", string(stop))
			time.Sleep(60 * time.Second)

			str := "ps aux|grep ulordd|grep -v grep"
			cmd = exec.Command("sh", "-c", str)
			out1, err := cmd.CombinedOutput()
			if err != nil {
				log.Println(err)
			}
			if string(out1) != "" {
				time.Sleep(60 * time.Second)
				continue
			}
			s := "ulordd"
			cmd = exec.Command("sh", "-c", s)
			go func() {
				cmd.Run()
				select {}
			}()
			time.Sleep(30 * time.Minute)
		}
		produce.Height = 0
		time.Sleep(5 * time.Second)
	}

}
func checkVersion() {
	var version Version
	url := "http://175.6.81.115:15944/getVersion"
	for {
		cmd := exec.Command("ulord-cli", "-version")
		out, err := cmd.CombinedOutput()
		if err != nil || string(out) == "" {
			log.Println(err)
			time.Sleep(1 * time.Minute)
			continue
		}

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
			continue
		}
		err = json.Unmarshal(body, &version)
		if err != nil {
			log.Println(err)
			time.Sleep(1 * time.Minute)
			continue
		}
		if version.Version != string(out) {
			//下载ulord包
			cmd := exec.Command("bash", "update.sh")
			_, err := cmd.CombinedOutput()
			if err != nil {
				log.Println(err)
				time.Sleep(1 * time.Minute)
				continue
			}
			//获取文件大小 s[0]
			cmd = exec.Command("du", "-sh", "../ulord/ulord_1_1_86.tgz")
			size, err := cmd.CombinedOutput()
			if err != nil {
				log.Println(err)
				time.Sleep(1 * time.Minute)
				continue
			}
			s := strings.Split(string(size), "M")
			log.Println("s[0]:", s[0])
			//小于36M就说明没下完
			if s[0] < "36" {
				time.Sleep(1 * time.Minute)
				continue
			}

			cmd = exec.Command("ulord-cli", "stop")
			stop, _ := cmd.CombinedOutput()
			log.Println("stop:", string(stop))
			time.Sleep(60 * time.Second)

			str := "ps aux|grep ulordd|grep -v grep"
			cmd = exec.Command("sh", "-c", str)
			out1, err := cmd.CombinedOutput()
			if err != nil {
				log.Println(err)
			}
			if string(out1) != "" {
				time.Sleep(60 * time.Second)
				continue
			}
			cmd = exec.Command("sh", "-c", "ulordd")
			go func() {
				cmd.Run()
				select {}
			}()
			time.Sleep(30 * time.Minute)
		}
		time.Sleep(6 * time.Hour)
	}
}

// GetStat return stat data
func GetStat(id string, name string) Stat {
	var statusInfo StatusInfo
	stat := Stat{}
	//发送版本信息
	cmd := exec.Command("ulord-cli", "-version")
	out, err := cmd.CombinedOutput()
	if err != nil || string(out) == "" {
		log.Println(err)
	}
	stat.Version = string(out)
	//发送IP id信息
	cmd = exec.Command("ulord-cli", "masternode", "status")
	status, err := cmd.CombinedOutput()
	if err != nil || string(out) == "" {
		log.Println(err)
	}
	err = json.Unmarshal(status, &statusInfo)
	if err != nil {
		json.Unmarshal(status, &statusInfo)
		log.Println(err)
	}
	stat.Ip = statusInfo.Service
	stat.Id = statusInfo.Masternodeindex
	//发送机器信息
	hostInfo, err := hoststat.GetHostInfo()
	if err != nil {
		SERVER_LOGGER.ErrorOnce("get host info failed: %s", err.Error())
	}
	stat.OSRelease = hostInfo.Release + " " + hostInfo.OSBit
	if err != nil {
		SERVER_LOGGER.ErrorOnce("get cpu info failed: %s", err.Error())
	}

	cpuStat, err := hoststat.GetCPUStat()
	if err != nil {
		SERVER_LOGGER.ErrorOnce("get cpu stat failed: %s", err.Error())
	}
	stat.CPURate = Round(100-cpuStat.IdleRate, 2)

	memStat, err := hoststat.GetMemStat()
	if err != nil {
		SERVER_LOGGER.ErrorOnce("get mem stat failed: %s", err.Error())
	}
	stat.MemRate = memStat.MemRate

	netStat, err := hoststat.GetNetStat()
	if err != nil {
		SERVER_LOGGER.ErrorOnce("get net stat failed: %s", err.Error())
	}
	netWrite := uint64(0)
	netRead := uint64(0)
	for _, v := range netStat {
		if v.Device != "lo" {
			netWrite += v.TXBytes
			netRead += v.RXBytes
		}
	}
	stat.NetWrite = netWrite
	stat.NetRead = netRead

	uptimeStat, err := hoststat.GetUptimeStat()
	if err != nil {
		SERVER_LOGGER.ErrorOnce("get uptime stat failed: %s", err.Error())
	}
	stat.Uptime = uint64(uptimeStat.Uptime)

	loadStat, err := hoststat.GetLoadStat()
	if err != nil {
		SERVER_LOGGER.ErrorOnce("get load stat failed: %s", err.Error())
	}
	stat.Load = fmt.Sprintf("%.2f %.2f %.2f", loadStat.LoadNow, loadStat.LoadPre, loadStat.LoadFar)
	return stat
}

//Get the height of the master node on the chain
func getChainHeight() int {
	var masterNodeHeight MasterNodeHeight
	url := "http://175.6.144.117:9879"
	jsonStr := `{"method":"masternode","params":["current"],"id":"curltest"}`
	req, err := http.NewRequest("POST", url, strings.NewReader(jsonStr))
	if err != nil {
		log.Println("get master node height Post:", err)
	}
	req.Header.Set("Content-ype", "text/plain;")
	req.Header.Add("Authorization", "Basic  VWxvcmQwMzpVbG9yZDAz")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("get master node height resp:", err)
	}
	defer resp.Body.Close()
	data, _ := ioutil.ReadAll(resp.Body)

	err = json.Unmarshal(data, &masterNodeHeight)

	return masterNodeHeight.Result.Height
}
