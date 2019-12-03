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
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/go-xorm/xorm"
)

// Round returns round number
func Round(data float64, precision int) (result float64) {
	pow := math.Pow(10, float64(precision))
	digit := pow * data
	_, div := math.Modf(digit)

	if div >= 0.5 {
		result = math.Ceil(digit)
	} else {
		result = math.Floor(digit)
	}
	result = result / pow

	return
}

// FileExists returns file is exists
func FileExists(fname string) bool {
	_, err := os.Stat(fname)
	return !os.IsNotExist(err)
}

func float64ToString(in float64) string {
	s1 := strconv.FormatFloat(in, 'f', -1, 64)
	return s1
}
func stringToFloat64(in string) float64 {
	defer func() {
		if err := recover(); err != nil {
			log.Println("[[Recovery] panic recovered:", err)
		}
	}()
	count, err := strconv.ParseFloat(in, 64)
	if err != nil {
		log.Fatalln("Failure of String Conversion.")
		panic(err)
	}
	return count
}

//获取engine对象
func getEngine() *xorm.Engine {
	engine, err := xorm.NewEngine("mysql", "root:123456@tcp(127.0.0.1:3306)/data?parseTime=true")
	if err != nil {
		log.Println("生成engine对象失败", err)
		engine, _ = xorm.NewEngine("mysql", "root:123456@tcp(127.0.0.1:3306)/data?parseTime=true")
	}
	engine.SetMaxOpenConns(5)
	return engine
}

// Chown do recurse chown go file or folder
func Chown(fname string, uid, gid int) (err error) {
	isDir, err := IsDir(fname)
	if err != nil {
		return
	}

	err = os.Chown(fname, uid, gid)
	if err != nil || !isDir {
		return
	}

	if !strings.HasSuffix(fname, "/") {
		fname += "/"
	}

	fs, err := ioutil.ReadDir(fname)
	if err != nil {
		return
	}

	for _, f := range fs {
		err = Chown(fname+f.Name(), uid, gid)
		if err != nil {
			return
		}
	}

	return
}

// IsDir returns if path is a dir
func IsDir(fname string) (bool, error) {
	f, err := os.Stat(fname)
	if err != nil {
		return false, err
	}

	return f.Mode().IsDir(), nil
}

// Md5 returns hex md5 of string
func Md5(str, key string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(str+key)))
}
