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
	"bytes"
	"net/http"
	"time"
)

// httpSend send data to stat api
func httpSend(server, stat string) (err error) {
	surl := server + "/receiveStat"

	request, err := http.NewRequest("POST", surl, bytes.NewBuffer([]byte(stat)))
	if err != nil {
		return
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: time.Duration(30 * time.Second),
	}

	response, err := client.Do(request)
	if err != nil {
		return
	}

	defer response.Body.Close()

	return nil
}
