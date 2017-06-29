package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
)

const SERVER_TOKEN = "3NUQNVD6HZXSPNADGFBTXOVTF4Y4ZV6N"

type WITClient struct {
	Server_token string
	WClient      *http.Client
}

func NewClient() WITClient {
	fmt.Println("Construct new WITClient now")
	client := &http.Client{}
	var witClient WITClient
	witClient.Server_token = SERVER_TOKEN
	witClient.WClient = client

	return witClient
}

func (c WITClient) PostSpeech(filePath string) {
	if c.WClient == nil {
		c = NewClient()
	}

	var data []byte
	var err error
	fmt.Println("uploading file " + filePath)
	if data, err = ioutil.ReadFile("./" + filePath); err == nil {

		t := time.Now()
		year, month, date := t.Date()
		s_year := strconv.Itoa(year)
		s_month := strconv.Itoa(int(month))
		s_date := strconv.Itoa(date)
		version := s_year + s_month + s_date
		speechURL := "https://api.wit.ai/speech?v=" + version

		if req, err := http.NewRequest("POST", speechURL, bytes.NewBuffer(data)); err == nil {
			req.Header.Set("Authorization", "Bearer 3NUQNVD6HZXSPNADGFBTXOVTF4Y4ZV6N")
			req.Header.Set("Accept", fmt.Sprintf("application/vnd.wit.%s+json", version))
			req.Header.Set("Content-Type", "audio/mpeg3")

			if resp, err := c.WClient.Do(req); err == nil {
				defer resp.Body.Close()
				defer os.Remove("./" + filePath)
				defer os.Remove("./" + "SpeechFile*")
				res, _ := ioutil.ReadAll(resp.Body)

				fmt.Println("HTTP response: " + string(res))
			} else {
				fmt.Println(err)
			}
		} else {
			fmt.Println(err)
		}

	} else {
		fmt.Println(err)
	}
	//c.WClient.post()
	//	resp, err := http.Post
}
