package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
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
			if update.Message.HasPrefix(update.Message.Text, "https://play.golang.org/p/") {
				err = respond(botURL, update)
			}
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
	result, err := parsePG(update.Message.Text)
	if err != nil {
		return err
	}

	var botMessage BotMessage
	botMessage.ChatId = update.Message.Chat.ChatId
	botMessage.Text = result
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

func parsePG(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return "", fmt.Errorf("getting %s: %s", url, resp.StatusCode)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("parsing %s as HTML: %v\n", url, err)
	}
	doc := string(body)
	re := regexp.MustCompile(`<textarea.*?>((.|\n)*)</textarea>`)
	code := re.FindStringSubmatch(doc)[1]
	return code, nil
}
