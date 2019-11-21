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
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
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
	http.HandleFunc("/api/stat", apiStatHandler)
	http.HandleFunc("/getMasterNodeByIndex", getMasterNodeByIndex)
	http.HandleFunc("/getMasterNode", getMasterNode)

	SERVER_LOGGER.Info("start http service")
	err := http.ListenAndServeTLS(":15944",
		SERVER_CONFIG.BaseDir+SERVER_CONFIG.TLSCert, SERVER_CONFIG.BaseDir+SERVER_CONFIG.TLSKey, nil)
	if err != nil {
		panic(err)
	}
}

//查询最新主节点信息
func getMasterNode(w http.ResponseWriter, r *http.Request) {
	engine := getEngine()
	defer engine.Close()
	data, err := engine.QueryString("select * from updatestat;")
	if err != nil {
		log.Println(err)
		data, _ = engine.QueryString("select * from updatestat;")
	}
	if data == nil {
		result := `{"status": {"code": 0, "message": "数据为空"}}`
		fmt.Fprintf(w, result)
		return
	} else {
		result := `{"status": {"code": 0, "data":%v}}`
		fmt.Fprintf(w, result, data)
		return
	}
}

//根据索引查主节点信息
func getMasterNodeByIndex(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	index := r.Form["index"][0]

	engine := getEngine()
	defer engine.Close()
	data, err := engine.QueryString("select * from statmysql where id=?;", index)
	if err != nil {
		log.Println(err)
		data, _ = engine.QueryString("select * from statmysql where id=?;", index)
	}
	if data == nil {
		result := `{"status": {"code": 0, "message": "数据为空"}}`
		fmt.Fprintf(w, result)
		return
	} else {
		result := `{"status": {"code": 0, "data":%v}}`
		fmt.Fprintf(w, result, data)
		return
	}
}

func apiStatHandler(w http.ResponseWriter, r *http.Request) {
	var stat Stat
	var statMysql Statmysql
	var updateStat Updatestat
	ip := getHTTPHeader(r, "X-Real-Ip")
	if ip == "" {
		ip = strings.Split(r.RemoteAddr, ":")[0]
	}

	clientKey := getHTTPHeader(r, "X-Client-Key")
	if clientKey == "" {
		result := `{"status": {"code": 0, "message": "key invalid"}}`
		fmt.Fprintf(w, result)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		result := `{"status": {"code": 0, "message": "body invalid"}}`
		fmt.Fprintf(w, result)
		return
	}

	// text := string(body)
	json.Unmarshal(body, &stat)
	if stat.Id != "" {
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
		_, err := engine.Insert(statMysql)
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
	}
	if body == nil {
		result := `{"status": {"code": 0, "message": "数据为空"}}`
		fmt.Fprintf(w, result)
		return
	}
	result := `{"status": {"code": 1, "message": "ok"}}`
	fmt.Fprintf(w, result)
}

func getHTTPHeader(r *http.Request, name string) string {
	if line, ok := r.Header[name]; ok {
		return line[0]
	}

	return ""
}

//获取engine对象
func getEngine() *xorm.Engine {
	engine, err := xorm.NewEngine("mysql", "root:123456@tcp(127.0.0.1:3306)/data?parseTime=true")
	if err != nil {
		log.Println("生成engine对象失败", err)
		engine, _ = xorm.NewEngine("mysql", "root:123456@tcp(127.0.0.1:3306)/data?parseTime=true")
	}
	return engine
}
