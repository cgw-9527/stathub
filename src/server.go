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
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

type Statmysql struct {
	Id             string  `xorm:"not null pk comment('主节点索引') INT(11)"`
	CpuRate        float64 `xorm:"not null comment('cpu利用率') FLOAT"`
	Raw            float64 `xorm:"not null comment('内存利用率') FLOAT"`
	NetRate        string  `xorm:"not null comment('网络带宽（net i/o）') VARCHAR(255)"`
	System         string  `xorm:"not null comment('操作系统') VARCHAR(255)"`
	Load           string  `xorm:"not null comment('机器负载') VARCHAR(255)"`
	OnlineTime     string  `xorm:"not null comment('在线时长') INT(11)"`
	BlockNum       string  `xorm:"not null comment('区块高度') INT(11)"`
	NodeStatus     string  `xorm:"not null comment('主节点状态') INT(11)"`
	ExpiryProducer string  `xorm:"not null comment('主节点证书到期时间') INT(11)"`
	IsselfProblock string  `xorm:"not null comment('是否自己出块') VARCHAR(255)"`
	TrxHash        string  `xorm:"not null comment('主节点交易hash') VARCHAR(255)"`
}
type ApiResult struct {
	Status ApiStatus `json:"status"`
}
type ApiStatus struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type Updatestat struct {
	Id             string  `xorm:"not null pk comment('主节点索引') INT(11)"`
	CpuRate        float64 `xorm:"not null comment('cpu利用率') FLOAT"`
	Raw            float64 `xorm:"not null comment('内存利用率') FLOAT"`
	NetRate        string  `xorm:"not null comment('网络带宽（net i/o）') VARCHAR(255)"`
	System         string  `xorm:"not null comment('操作系统') VARCHAR(255)"`
	Load           string  `xorm:"not null comment('机器负载') VARCHAR(255)"`
	OnlineTime     string  `xorm:"not null comment('在线时长') INT(11)"`
	BlockNum       string  `xorm:"not null comment('区块高度') INT(11)"`
	NodeStatus     string  `xorm:"not null comment('主节点状态') INT(11)"`
	ExpiryProducer string  `xorm:"not null comment('主节点证书到期时间') INT(11)"`
	IsselfProblock string  `xorm:"not null comment('是否自己出块') VARCHAR(255)"`
	TrxHash        string  `xorm:"not null comment('主节点交易hash') VARCHAR(255)"`
}

// HttpService start http service
func HttpService() {
	http.HandleFunc("/sendStat", sendStat)
	http.HandleFunc("/getMasterNodeStatus", getMasterNodeStatus)
	http.HandleFunc("/getMasterNodeList", getMasterNodeList)

	SERVER_LOGGER.Info("start http service")
	err := http.ListenAndServeTLS(":15944",
		SERVER_CONFIG.BaseDir+SERVER_CONFIG.TLSCert, SERVER_CONFIG.BaseDir+SERVER_CONFIG.TLSKey, nil)
	if err != nil {
		panic(err)
	}
}

//查询最新主节点信息
func getMasterNodeList(w http.ResponseWriter, r *http.Request) {
	var updateStat Updatestat
	var updateStats []Updatestat
	engine := getEngine()
	defer engine.Close()
	datas, err := engine.QueryString("select * from updatestat;")
	if err != nil {
		log.Println(err)
		datas, _ = engine.QueryString("select * from updatestat;")
	}
	if datas == nil {
		result := `{"status": {"code": 0, "message": "数据为空"}}`
		fmt.Fprintf(w, result)
		return
	} else {
		for _, data := range datas {
			//将map转化为json
			updateStat.Id = data["id"]
			updateStat.CpuRate = stringToFloat64(data["cpu_rate"])
			updateStat.Raw = stringToFloat64(data["raw"])
			updateStat.NetRate = data["net_rate"]
			updateStat.System = data["system"]
			updateStat.Load = data["load"]
			updateStat.OnlineTime = data["online_time"]
			updateStat.BlockNum = data["block_num"]
			updateStat.NodeStatus = data["node_status"]
			updateStat.ExpiryProducer = data["expiry_producer"]
			updateStat.IsselfProblock = data["isself_problock"]
			updateStat.TrxHash = data["trx_hash"]
			updateStats = append(updateStats, updateStat)
		}

		result := `{"status": {"code": 0, "data":%v}}`
		fmt.Fprintf(w, result, updateStats)
		return
	}
}

//根据索引查主节点信息
func getMasterNodeStatus(w http.ResponseWriter, r *http.Request) {
	var statMysql Statmysql
	var statMysqls []Statmysql
	r.ParseForm()
	index := r.Form["index"][0]
	pageNum, _ := strconv.Atoi(r.Form["pageNum"][0])
	pageSize, _ := strconv.Atoi(r.Form["pageSize"][0])

	start := (pageNum - 1) * pageSize
	offset := pageSize

	if pageNum == 0 || pageSize == 0 {
		result := `{"status": {"code": 0,"message":"pageNum或者pageSize不能为空"}}`
		fmt.Fprintf(w, result)
		return
	}
	engine := getEngine()
	defer engine.Close()
	datas, err := engine.QueryString("select * from statmysql where id=? limit ?,?;", index, start, offset)
	if err != nil {
		log.Println(err)
		datas, _ = engine.QueryString("select * from statmysql where id=? limit ?,?;", index, start, offset)
	}
	if datas == nil {
		result := `{"status": {"code": 0, "message": "数据为空"}}`
		fmt.Fprintf(w, result)
		return
	} else {
		for _, data := range datas {
			statMysql.Id = data["id"]
			statMysql.CpuRate = stringToFloat64(data["cpu_rate"])
			statMysql.Raw = stringToFloat64(data["raw"])
			statMysql.NetRate = data["net_rate"]
			statMysql.System = data["system"]
			statMysql.Load = data["load"]
			statMysql.OnlineTime = data["online_time"]
			statMysql.BlockNum = data["block_num"]
			statMysql.NodeStatus = data["node_status"]
			statMysql.ExpiryProducer = data["expiry_producer"]
			statMysql.IsselfProblock = data["isself_problock"]
			statMysql.TrxHash = data["trx_hash"]
			statMysqls = append(statMysqls, statMysql)
		}
		result := `{"status": {"code": 0, "data":%v}}`
		fmt.Fprintf(w, result, statMysqls)
		return
	}
}

func sendStat(w http.ResponseWriter, r *http.Request) {
	var stat Stat
	var statMysql Statmysql
	var updateStat Updatestat

	apiResult := ApiResult{
		Status: ApiStatus{},
	}

	clientKey := getHTTPHeader(r, "X-Client-Key")
	if clientKey == "" {
		result := `{"status": {"code": 0, "message": "key X-Client-Key invalid"}}`
		fmt.Fprintf(w, result)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		result := `{"status": {"code": 0, "message": "body invalid"}}`
		fmt.Fprintf(w, result)
		return
	}
	serverKey := Md5(SERVER_CONFIG.ServerKey, string(body))
	if serverKey != clientKey {
		result := `{"status": {"code": 0, "message": "key invalid"}}`
		fmt.Fprintf(w, result)
		return
	}
	// text := string(body)
	json.Unmarshal(body, &stat)
	//保存数据
	statMysql.Id = stat.Id
	statMysql.CpuRate = stat.CPURate
	statMysql.Raw = stat.MemRate
	statMysql.NetRate = fmt.Sprintf("%v / %v", stat.NetRead, stat.NetWrite)
	statMysql.System = stat.OSRelease
	statMysql.Load = stat.Load
	statMysql.OnlineTime = fmt.Sprintf("%v", stat.Uptime)
	statMysql.BlockNum = stat.BlockNum
	statMysql.NodeStatus = stat.NodeStatus
	statMysql.ExpiryProducer = stat.ExpiryProducer
	statMysql.IsselfProblock = stat.IsselfProblock
	statMysql.TrxHash = stat.TrxHash

	engine := getEngine()
	defer engine.Close()

	_, err = engine.Insert(statMysql)
	if err != nil {
		log.Println(err)
		engine.Insert(statMysql)
	}
	//更新数据
	updateStat.Id = stat.Id
	updateStat.CpuRate = stat.CPURate
	updateStat.Raw = stat.MemRate
	updateStat.NetRate = fmt.Sprintf("%v / %v", stat.NetRead, stat.NetWrite)
	updateStat.System = stat.OSRelease
	updateStat.Load = stat.Load
	updateStat.OnlineTime = fmt.Sprintf("%v", stat.Uptime)
	updateStat.BlockNum = stat.BlockNum
	updateStat.NodeStatus = stat.NodeStatus
	updateStat.ExpiryProducer = stat.ExpiryProducer
	updateStat.IsselfProblock = stat.IsselfProblock
	updateStat.TrxHash = stat.TrxHash
	result, err := engine.QueryString("select * from updatestat where id=?;", stat.Id)
	if err != nil {
		log.Println(err)
		result, _ = engine.QueryString("select * from updatestat where id=?;", stat.Id)
	}
	if result != nil {
		_, err := engine.Id(stat.Id).Update(updateStat)
		if err != nil {
			log.Println(err)
			engine.Id(stat.Id).Update(updateStat)
		}
	} else {
		_, err := engine.Insert(updateStat)
		if err != nil {
			log.Println(err)
			engine.Insert(updateStat)
		}

	}
	if body == nil {
		apiResult.Status.Code = 0
		apiResult.Status.Message = "数据为空"
		return
	}
	result1 := `{"status": {"code": 1, "message": "ok"}}`
	fmt.Fprintf(w, result1)
	return
}

func getHTTPHeader(r *http.Request, name string) string {
	if line, ok := r.Header[name]; ok {
		return line[0]
	}

	return ""
}
