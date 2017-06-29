package main

// Package telebot provides a handy wrapper for interactions
// with Telegram bots.
//
// Here is an example of helloworld bot implementation:
//
import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os/exec"
	"time"

	"github.com/tucnak/telebot"
)

func downloadFile(bot *telebot.Bot, fileID string) (ret string, err error) {
	url, err := bot.GetFileDirectURL(fileID)
	tmpFile, err := ioutil.TempFile("/home/ingemar/workspace/go/src/github.com/ingemar0720/speechTeleBot", "SpeechFile")
	if err != nil {
		fmt.Println("create temp file error")
		return "", err
	}

	//defer os.Remove(tmpFile.Name())
	var resp *http.Response
	resp, err = http.Get(url)
	fmt.Println(url)

	if err == nil {
		defer resp.Body.Close()

		if _, err = io.Copy(tmpFile, resp.Body); err == nil {
			fmt.Println("download file " + tmpFile.Name())
			return tmpFile.Name(), err
		} else {
			fmt.Println("download file fail")
			return "", nil
		}
	} else {
		return "", nil
	}
}

func transcodeToMP3(ogaFile string) (mp3File string) {
	params := []string{"-i", ogaFile, "./Speech.mp3"}
	cmd := exec.Command("ffmpeg", params...)
	if _, err := cmd.CombinedOutput(); err != nil {
		return ""
	} else {
		return "Speech.mp3"
	}
}

func main() {
	//bot, err := telebot.NewBot("SECRET_TOKEN")
	bot, err := telebot.NewBot("370622763:AAHWWgh36fYrOheTNWPYLg6oqCZ3UwL77-g")
	if err != nil {
		return
	}

	messages := make(chan telebot.Message)
	bot.Listen(messages, 1*time.Second)

	for message := range messages {
		if message.Text == "/hi" {
			bot.SendMessage(message.Chat,
				"Hello, "+message.Sender.FirstName+"!", nil)
			/*} else if message.Text == "/photo" {
			fmt.Println("here in photo")
			photoObj := new(telebot.Photo)
			photoObj.File, err = telebot.NewFile("/home/ingemar/workspace/go/src/github.com/ingemar0720/speechTeleBot/test.jpg")
			if err == nil {
				fmt.Println("here to send photo")
				photoObj.Caption = "testPhoto"
				err := bot.SendPhoto(message.Chat, photoObj, nil)
				if err != nil {
					fmt.Println(err)
				}
			} else {
				fmt.Println("new file fail")
				fmt.Println(err)
			}
			*/
		} else if message.Voice.Exists() {
			fmt.Println("voice file exist")
			if ogaFile, err := downloadFile(bot, message.Voice.FileID); err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("transcode to mp3 file now")
				var mp3file string
				if mp3file = transcodeToMP3(ogaFile); mp3file != "" {
					fmt.Println("transcode mp3 file success")
				}
				var c WITClient
				c.PostSpeech(mp3file)
			}
		} else {
			fmt.Println("receive unrecognized message")
		}
	}
}
