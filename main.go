package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

const botAPI = "https://api.telegram.org/bot"

var botToken string

func init() {
	if len(os.Args) == 2 {
		botToken = os.Args[1]
	} else {
		fmt.Println("Укажите токен первым аргументом")
		os.Exit(2)
	}
}
func main() {
	botURL := botAPI + botToken
	offset := 0
	for {
		updates, err := getUpdates(botURL, offset)
		if err != nil {
			log.Printf("Ошибка при получении апдейтов: %v\n", err.Error())
		}
		for _, update := range updates {
			err = respond(botURL, update)
			offset = update.UpdateID + 1
		}
		fmt.Println(updates)
	}
}

// запрос обновления
func getUpdates(botURL string, offset int) ([]Update, error) {
	resp, err := http.Get(botURL + "/getUpdates" + "?offset=" + strconv.Itoa(offset))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var restResponse RestResponse
	err = json.Unmarshal(body, &restResponse)
	if err != nil {
		return nil, err
	}
	return restResponse.Result, nil
}

func respond(botURL string, update Update) error {
	var botMessage BotMessage
	botMessage.ChatID = update.Message.Chat.ChatID
	botMessage.Text = update.Message.Text
	botMessage.Date = time.Now().Unix()
	buf, err := json.Marshal(botMessage)
	if err != nil {
		return err
	}
	_, err = http.Post(botURL+"/sendMessage", "application/json", bytes.NewBuffer(buf))
	if err != nil {
		return err
	}
	return nil
}
