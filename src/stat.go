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
	"fmt"
	"github.com/likexian/host-stat-go"
	"os/exec"
	"strings"
	"strconv"
	"encoding/json"
)

// Stat storing stat data
type Stat struct {
	Id        string  `json:"id"`//
	OSRelease string  `json:"os_release"`//
	Uptime    uint64  `json:"uptime"`//
	Load      string  `json:"load"`//
	CPURate   float64 `json:"cpu_rate"`//
	MemRate   float64 `json:"mem_rate"`//
	NetRead   uint64  `json:"net_read"`//
	NetWrite  uint64  `json:"net_write"`//
	TrxHash         string `json:"trx_hash"`
	NodeStatus     string `json:"node_status"`
	BlockNum       string `json:"block_num"`
	ExpiryProducer string `json:"expiry_producer"`
	IsselfProblock string `json:"isself_problock"`
}

type MasterNode struct {
	TrxHash         string
	NodeStatus     string
	Id             string
	BlockNum       string
	ExpiryProducer string
	IsselfProblock string
}
type Produce struct {
	Height    int `json:"height:"`
	Produceno int `json:"produceno:"`
}
func getMasterNodeList()[]MasterNode {
	var masterNode MasterNode
	var produce Produce
	var masterNodeList []MasterNode
	cmd := exec.Command("ulord-cli", "masternodelist", "full")
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(err)
	}
	cmd = exec.Command("ulord-cli", "masternode", "current")
	out1, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(err)
	}
	err = json.Unmarshal(out1, &produce)
	if err != nil {
		fmt.Println(err)
	}
	str := strings.TrimSuffix(strings.TrimPrefix(strings.TrimSpace(string(out)), "{"), "}")
	linesData := strings.Split(str, ",")

	for _, lineData := range linesData {
		s := strings.SplitN(lineData, ":", 2)
		txHash := strings.Split(strings.TrimSuffix(strings.TrimPrefix(strings.TrimSpace(s[0]), `"`), `"`), "-")
		masterNode.TrxHash = txHash[0]

		data := strings.Split(strings.TrimSuffix(strings.TrimSpace(strings.TrimPrefix(s[1], " "+`"`)), `"`), " ")
		masterNode.Id = data[7]
		masterNode.BlockNum = data[6]
		masterNode.ExpiryProducer = data[4]
		masterNode.NodeStatus = data[8]
		if strconv.Itoa(produce.Produceno) == data[7] {
			masterNode.IsselfProblock = "true"
		} else {
			masterNode.IsselfProblock = "false"
		}
		masterNodeList = append(masterNodeList, masterNode)
	}
	return masterNodeList
}
// GetStat return stat data
func GetStat(id string, name string) Stat {

	stat := Stat{}

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
