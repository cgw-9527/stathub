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
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/likexian/gokit/xdaemon"
	"github.com/likexian/gokit/xlog"
	"github.com/likexian/gokit/xos"
)

var (
	// SERVER_START is server start timestamp
	SERVER_START = int64(0)
	// SERVER_CONFIG is server config data
	SERVER_CONFIG = Config{}
	// SERVER_LOGGER is server logger
	SERVER_LOGGER = xlog.New(os.Stderr, xlog.INFO)
)

func main() {
	go checkVersion()
	SERVER_START = time.Now().Unix()

	if DEBUG {
		SERVER_LOGGER = xlog.New(os.Stderr, xlog.DEBUG)
	} else {
		SERVER_LOGGER.SetSizeRotate(3, 100*1024*1024)
	}

	configFile := flag.String("c", "", "set configuration file")
	initClient := flag.Bool("init-client", false, "init client configuration")

	serverUrl := "http://175.6.81.115:15944"

	flag.Parse()

	if *configFile == "" {
		flag.Usage()
		os.Exit(1)
	}

	if *initClient {
		if serverUrl == "" {
			SERVER_LOGGER.Fatal("server url is required, set it by --server-url.")
		}
		timeStamp := fmt.Sprintf("%d", SERVER_START)
		id := Md5(fmt.Sprintf("%d", os.Getpid()), timeStamp)
		err := newClientConfig(*configFile, id, "", serverUrl)
		if err != nil {
			SERVER_LOGGER.Fatal(err.Error())
		} else {
			SERVER_LOGGER.Info("init client configuration successful")
			os.Exit(0)
		}
	}

	if !FileExists(*configFile) {
		SERVER_LOGGER.Fatal(fmt.Sprintf("configuration file %s is not found.\n", *configFile))
	}

	var err error
	SERVER_CONFIG, err = GetConfig(*configFile)
	if err != nil {
		SERVER_LOGGER.Fatal(fmt.Sprintf("configuration load failed, %s", err.Error()))
	}

	for _, v := range []string{*configFile, SERVER_CONFIG.PidFile, SERVER_CONFIG.LogFile} {
		ds, _ := filepath.Split(SERVER_CONFIG.BaseDir + v)
		if ds != "" && !FileExists(ds) {
			err := os.MkdirAll(ds, 0755)
			if err != nil {
				SERVER_LOGGER.Fatal(err.Error())
			}
		}
	}

	if !DEBUG && SERVER_CONFIG.DaemonUser != "" {
		uid, gid, err := xos.LookupUser(SERVER_CONFIG.DaemonUser)
		if err != nil {
			SERVER_LOGGER.Fatal("Lookup daemon user %s failed, %s", SERVER_CONFIG.DaemonUser, err.Error())
		}
		err = Chown(SERVER_CONFIG.BaseDir, uid, gid)
		if err != nil {
			SERVER_LOGGER.Fatal("Chown base dir to daemon user %s failed, %s", SERVER_CONFIG.DaemonUser, err.Error())
		}
	}

	if !DEBUG {
		c := xdaemon.Config{
			Pid:   SERVER_CONFIG.BaseDir + SERVER_CONFIG.PidFile,
			Log:   SERVER_CONFIG.BaseDir + SERVER_CONFIG.LogFile,
			User:  SERVER_CONFIG.DaemonUser,
			Chdir: "",
		}
		err := c.Daemon()
		if err != nil {
			SERVER_LOGGER.Fatal(err.Error())
		}
	}
	SERVER_LOGGER.Info("server start at %d", SERVER_START)
	go checkStatus()
	StatService()
}
