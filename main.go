// Copyright 2022 Carlos Rabelo.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/proxy"
)

type Response struct {
	Client  string `json:"client"`
	Ip      string `json:"ip"`
	Name    string `json:"name"`
	Port    int    `json:"port"`
	Region  string `json:"region"`
	Server  string `json:"server"`
	Success bool   `json:"success"`
}

var USERNAME = "solracolebar"
var DIFFICULTY = ""
var MINING_KEY = "None"
var RIG_IDENTIFIER = "pc"

var MINER_BANNER = "Ministo"
var MINER_VERSION = "0.1"

func main() {

	//

	single_miner_id := rand.Intn(2811)

	//

	response, err := http.Get("https://server.duinocoin.com/getPool")

	if err != nil {
		log.Println("No response from request getpool")
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)

	var result Response

	if err := json.Unmarshal(body, &result); err != nil {
		log.Println("Can't resolve JSON")
	}

	//

	u, err := url.Parse("socks5://localhost:9050")

	if nil != err {
		log.Fatalf("Parse: %v", err)
	}

	//

	d, err := proxy.FromURL(u, proxy.Direct)

	if nil != err {
		log.Fatalf("Proxy: %v", err)
	}

	//

	conn, _ := d.Dial("tcp", result.Ip+":"+fmt.Sprint(result.Port))

	defer conn.Close()

	//

	buffer := make([]byte, 100)

	_, err = conn.Read(buffer)

	if err != nil {
		log.Println(err)
		log.Fatal("Error getting the version")
	}

	buffer = bytes.Trim(buffer, "\x00")

	//

	version := strings.TrimSpace(string(buffer))

	log.Printf("Server is on version: %v", version)

	//

	for {

		//

		_, err = conn.Write([]byte(fmt.Sprintf("%v,%v,%v,%v", "JOB", USERNAME, DIFFICULTY, MINING_KEY)))

		if err != nil {
			log.Fatal("Error requesting job")
		}

		//

		buffer = make([]byte, 1024)

		_, err = conn.Read(buffer)

		if err != nil {
			log.Println(err)
			log.Fatal("Error getting the job")
		}

		buffer = bytes.Trim(buffer, "\x00")

		//

		job := strings.Split(strings.TrimSpace(string(buffer)), ",")

		hash := job[0]
		goal := job[1]

		diff, _ := strconv.Atoi(job[2])

		hashingStartTime := time.Now()

		for result := 0; result <= diff*100; result++ {

			//

			h := sha1.New()

			h.Write([]byte(hash + strconv.Itoa(result)))

			nh := hex.EncodeToString(h.Sum(nil))

			//

			if goal == nh {

				//

				timeDifference := time.Since(hashingStartTime).Seconds()

				hashrate := int(float64(result) / timeDifference)

				khashrate := hashrate / 1000

				//

				_, err = conn.Write([]byte(fmt.Sprintf("%v,%v,%v %v,%v,,%v", result, hashrate, MINER_BANNER, MINER_VERSION, RIG_IDENTIFIER, single_miner_id)))

				if err != nil {
					log.Fatal("Error sending then job")
				}

				//

				buffer = make([]byte, 1024)

				_, err = conn.Read(buffer)

				if err != nil {
					log.Println(err)
					log.Fatal("Error getting the feedback")
				}

				buffer = bytes.Trim(buffer, "\x00")

				//

				feedback := strings.TrimSpace(string(buffer))

				if feedback == "GOOD" {

					log.Printf("Accepted share %v Hashrate %v kH/s Difficulty %v", result, khashrate, diff)

					break

				} else if feedback == "BAD" {

					log.Printf("Rejected share %v Hashrate %v kH/s Difficulty %v", result, khashrate, diff)

					break

				} else {

					log.Fatal(fmt.Sprintf("Feedback: %v", feedback))

					break

				}

			}

		}

	}

}
