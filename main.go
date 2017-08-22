package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"

	"golang.org/x/crypto/ssh/terminal"

	"github.com/gertd/pdbq/helper"
)

// Token - Puppet RBAC token
type Token struct {
	Token string `json:"token"`
}

func main() {

	var token Token
	var hostname string
	var username string
	var password string

	flag.StringVar(&hostname, "hostname", "", "puppet hostname")
	flag.StringVar(&username, "username", "", "puppet username")
	flag.StringVar(&password, "password", "", "puppet password")

	flag.Parse()

	if len(password) == 0 {
		fmt.Println("puppet password")
		buf, err := terminal.ReadPassword(0)
		if err != nil {
			log.Fatalln(err)
		}
		password = string(buf)
	}

	{
		const portnumber = 4433
		const endpoint = `rbac-api/v1/auth/token`

		url := fmt.Sprintf("https://%s:%d/%s", hostname, portnumber, endpoint)

		var jsonStr = []byte(fmt.Sprintf(`{"login": "%s", "password": "%s"}`, username, password))

		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
		if err != nil {
			log.Fatalln("err ", err)
		}
		req.Header.Set("Content-Type", "application/json")

		tr := &http.Transport{
			Dial: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).Dial,
			TLSHandshakeTimeout: 10 * time.Second,
			TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
			DisableCompression:  true,
		}
		client := &http.Client{Transport: tr}

		resp, err := client.Do(req)
		if err != nil {
			log.Fatalln("err ", err)
		}
		if resp != nil {
			defer resp.Body.Close()
		}

		buf, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln("err ", err)
		}

		err = json.Unmarshal(buf, &token)
		if err != nil {
			log.Fatalln("err ", err)
		}

	}
	{
		const portnumber = 8081
		const endpoint = `pdb/query/v4/inventory`

		url := fmt.Sprintf("https://%s:%d/%s", hostname, portnumber, endpoint)

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			log.Fatalln("err ", err)
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Authentication", token.Token)

		tr := &http.Transport{
			Dial: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).Dial,
			TLSHandshakeTimeout: 10 * time.Second,
			TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
			DisableCompression:  true,
		}
		client := &http.Client{Transport: tr}

		resp, err := client.Do(req)
		if err != nil {
			log.Fatalln("err ", err)
		}
		if resp != nil {
			defer resp.Body.Close()
		}

		buf, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln("err ", err)
		}

		var result interface{}
		err = json.Unmarshal(buf, &result)
		if err != nil {
			log.Fatalln("err ", err)
		}

		fmt.Println(helper.PrettyPrintJSON(result))
	}
}
