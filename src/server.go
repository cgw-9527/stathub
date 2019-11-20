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
	"os"
	"strings"
	"text/template"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/likexian/simplejson-go"
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
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/logout", logoutHandler)
	http.HandleFunc("/passwd", passwdHandler)
	http.HandleFunc("/help", helpHandler)
	http.HandleFunc("/node", nodeHandler)
	http.HandleFunc("/pkgs/", pkgsHandler)
	http.HandleFunc("/static/", staticHandler)
	http.HandleFunc("/robots.txt", robotsTXTHandler)
	http.HandleFunc("/api/stat", apiStatHandler)
	http.HandleFunc("/api/node", apiNodeHandler)

	SERVER_LOGGER.Info("start http service")
	err := http.ListenAndServeTLS(":15944",
		SERVER_CONFIG.BaseDir+SERVER_CONFIG.TLSCert, SERVER_CONFIG.BaseDir+SERVER_CONFIG.TLSKey, nil)
	if err != nil {
		panic(err)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if isRobots(w, r) {
		httpError(w, r, http.StatusForbidden)
		return
	}

	if !isLogin(w, r) {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	if r.URL.Path != "/" {
		httpError(w, r, http.StatusNotFound)
		return
	}

	tpl, err := template.New("index").Parse(TPL_TEMPLATE["layout.html"])
	if err != nil {
		httpError(w, r, http.StatusInternalServerError)
		return
	}

	tpl, err = tpl.Parse(TPL_TEMPLATE["index.html"])
	if err != nil {
		httpError(w, r, http.StatusInternalServerError)
		return
	}

	if DEBUG {
		tpl, err = template.ParseFiles("template/layout.html", "template/index.html")
		if err != nil {
			httpError(w, r, http.StatusInternalServerError)
			return
		}
	}

	status := ReadStatus(SERVER_CONFIG.DataDir)
	data := map[string]interface{}{
		"data":    status,
		"version": Version(),
	}
	tpl.Execute(w, data)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if isRobots(w, r) {
		httpError(w, r, http.StatusForbidden)
		return
	}

	if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			httpError(w, r, http.StatusInternalServerError)
			return
		}

		password := r.PostForm.Get("password")
		if Md5(SERVER_CONFIG.ServerKey, password) != SERVER_CONFIG.PassWord {
			http.Redirect(w, r, "/login", http.StatusFound)
		} else {
			value := Md5(SERVER_CONFIG.ServerKey, SERVER_CONFIG.PassWord)
			cookie := http.Cookie{Name: "id", Value: value, HttpOnly: true}
			http.SetCookie(w, &cookie)
			http.Redirect(w, r, "/", http.StatusFound)
		}
	} else {
		tpl, err := template.New("login").Parse(TPL_TEMPLATE["layout.html"])
		if err != nil {
			httpError(w, r, http.StatusInternalServerError)
			return
		}

		tpl, err = tpl.Parse(TPL_TEMPLATE["login.html"])
		if err != nil {
			httpError(w, r, http.StatusInternalServerError)
			return
		}

		if DEBUG {
			tpl, err = template.ParseFiles("template/layout.html", "template/login.html")
			if err != nil {
				httpError(w, r, http.StatusInternalServerError)
				return
			}
		}

		tpl.Execute(w, map[string]string{"action": "login"})
	}
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	expires := time.Now()
	expires = expires.AddDate(0, 0, -1)
	cookie := http.Cookie{Name: "id", Value: "", Expires: expires, HttpOnly: true}
	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/login", http.StatusFound)
	return
}

func passwdHandler(w http.ResponseWriter, r *http.Request) {
	if !isLogin(w, r) {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			httpError(w, r, http.StatusInternalServerError)
			return
		}

		password := r.PostForm.Get("password")
		if password == "" {
			http.Redirect(w, r, "/passwd", http.StatusFound)
		} else {
			SERVER_CONFIG.PassWord = Md5(SERVER_CONFIG.ServerKey, password)
			err := SaveConfig(SERVER_CONFIG)
			if err != nil {
				httpError(w, r, http.StatusInternalServerError)
				return
			}
			http.Redirect(w, r, "/", http.StatusFound)
		}
	} else {
		tpl, err := template.New("passwd").Parse(TPL_TEMPLATE["layout.html"])
		if err != nil {
			httpError(w, r, http.StatusInternalServerError)
			return
		}

		tpl, err = tpl.Parse(TPL_TEMPLATE["login.html"])
		if err != nil {
			httpError(w, r, http.StatusInternalServerError)
			return
		}

		if DEBUG {
			tpl, err = template.ParseFiles("template/layout.html", "template/login.html")
			if err != nil {
				httpError(w, r, http.StatusInternalServerError)
				return
			}
		}

		tpl.Execute(w, map[string]string{"action": "passwd"})
	}
}

func helpHandler(w http.ResponseWriter, r *http.Request) {
	if !isLogin(w, r) {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	tpl, err := template.New("help").Parse(TPL_TEMPLATE["layout.html"])
	if err != nil {
		httpError(w, r, http.StatusInternalServerError)
		return
	}

	tpl, err = tpl.Parse(TPL_TEMPLATE["help.html"])
	if err != nil {
		httpError(w, r, http.StatusInternalServerError)
		return
	}

	if DEBUG {
		tpl, err = template.ParseFiles("template/layout.html", "template/help.html")
		if err != nil {
			httpError(w, r, http.StatusInternalServerError)
			return
		}
	}

	tpl.Execute(w, map[string]string{"server": r.Host, "key": SERVER_CONFIG.ServerKey})
}

func nodeHandler(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	if key != SERVER_CONFIG.ServerKey {
		httpError(w, r, http.StatusForbidden)
		return
	}

	tpl, err := template.New("node").Parse(TPL_TEMPLATE["node.html"])
	if err != nil {
		httpError(w, r, http.StatusInternalServerError)
		return
	}

	if DEBUG {
		tpl, err = template.ParseFiles("template/node.html")
		if err != nil {
			httpError(w, r, http.StatusInternalServerError)
			return
		}
	}

	tpl.Execute(w, map[string]string{"server": r.Host, "key": SERVER_CONFIG.ServerKey, "version": Version()})
}

func robotsTXTHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "User-agent: *\r\nDisallow: /")
}

func apiNodeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if !isLogin(w, r) {
		result := `{"status": {"code": 0, "message": "login timeout"}}`
		fmt.Fprintf(w, result)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		result := `{"status": {"code": 0, "message": "data error"}}`
		fmt.Fprintf(w, result)
		return
	}

	data, err := simplejson.Loads(string(body))
	if err != nil {
		result := `{"status": {"code": 0, "message": "data invalid"}}`
		fmt.Fprintf(w, result)
		return
	}

	dataId, _ := data.Get("id").String()
	dataIdDir := SERVER_CONFIG.BaseDir + SERVER_CONFIG.DataDir + "/" + dataId[3:]
	if !FileExists(dataIdDir) {
		result := `{"status": {"code": 0, "message": "node id invalid"}}`
		fmt.Fprintf(w, result)
		return
	}

	err = os.RemoveAll(dataIdDir)
	if err != nil {
		result := `{"status": {"code": 0, "message": "` + err.Error() + `"}}`
		fmt.Fprintf(w, result)
		return
	}

	result := `{"status": {"code": 1, "message": "ok"}}`
	fmt.Fprintf(w, result)
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

func staticHandler(w http.ResponseWriter, r *http.Request) {
	n := strings.LastIndex(r.URL.Path, ".")
	if n == -1 {
		httpError(w, r, http.StatusNotFound)
		return
	}

	ext := r.URL.Path[n+1:]
	if ext == "css" {
		w.Header().Set("Content-Type", "text/css; charset=utf-8")
	} else if ext == "js" {
		w.Header().Set("Content-Type", "text/javascript; charset=utf-8")
	}

	if DEBUG {
		http.ServeFile(w, r, r.URL.Path[1:])
	} else {
		if test, ok := TPL_STATIC[r.URL.Path[8:]]; ok {
			fmt.Fprint(w, test)
		} else {
			httpError(w, r, http.StatusNotFound)
		}
	}
}

func pkgsHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, r.URL.Path[1:])
}

func getHTTPHeader(r *http.Request, name string) string {
	if line, ok := r.Header[name]; ok {
		return line[0]
	}

	return ""
}

// httpError returns a http error
func httpError(w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status)
	if status == http.StatusForbidden {
		fmt.Fprint(w, "<title>Forbidden</title><h1>Forbidden</h1>")
	} else if status == http.StatusNotFound {
		fmt.Fprint(w, "<title>Not Found</title><h1>Not Found</h1>")
	} else if status == http.StatusInternalServerError {
		fmt.Fprint(w, "<title>Internal Server Error</title><h1>Internal Server Error</h1>")
	}
}

// isLogin returns request has login
func isLogin(w http.ResponseWriter, r *http.Request) bool {
	cookie, err := r.Cookie("id")
	if err != nil || cookie.Value == "" {
		return false
	}

	value := Md5(SERVER_CONFIG.ServerKey, SERVER_CONFIG.PassWord)
	if value != cookie.Value {
		return false
	}

	return true
}

// isRobots returns is a robot request
func isRobots(w http.ResponseWriter, r *http.Request) bool {
	agent := strings.ToLower(getHTTPHeader(r, "User-Agent"))
	robots := []string{"bot", "spider", "archiver", "yahoo! slurp", "haosou"}
	for _, v := range robots {
		if strings.Contains(agent, v) {
			return true
		}
	}

	return false
}

//获取engine对象
func getEngine() *xorm.Engine {
	engine, err := xorm.NewEngine("mysql", "root:123456@tcp(127.0.0.1:3306)/data?parseTime=true")
	if err != nil {
		log.Println("生成engine对象失败", err)
		panic(err)
	}
	return engine
}
