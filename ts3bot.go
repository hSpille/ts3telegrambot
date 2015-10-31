package main

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/Syfaro/telegram-bot-api"
	"github.com/toqueteos/ts3"
	"log"
	"strings"
	"time"
)

type tomlConfig struct {
	Tsuser          string
	Tspasswd        string
	Tsurl           string
	Telegrammapikey string
	Telegrammchatid int
}

func main() {
	var config tomlConfig
	if _, err := toml.DecodeFile("config.toml", &config); err != nil {
		fmt.Println(err)
		return
	}
	tgBot, err := tgbotapi.NewBotAPI(config.Telegrammapikey)
	if err != nil {
		log.Panic(err)
	}
	tgBot.Debug = false
	log.Printf("Telegramm Bot authorized on account %s", tgBot.Self.UserName)
	log.Printf("Connecting to TS3 Server at: %s \n", config.Tsurl)
	tsConn, err := ts3.Dial(config.Tsurl, true)
	if err != nil {
		panic(err)
	}

	defer tsConn.Cmd("quit")
	defer tsConn.Close()
	var oldState []string
	var newState []string
	for {
		onlineUsers := tsBot(tsConn, config.Tsuser, config.Tspasswd)
		for _, onlineUser := range onlineUsers {
			if !contains(oldState, onlineUser) {
				msg := tgbotapi.NewMessage(config.Telegrammchatid, fmt.Sprintf("%v im Teamspeak!", onlineUser))
				tgBot.SendMessage(msg)
				time.Sleep(50 * time.Millisecond)
				newState = append(newState, onlineUser)
			}
		}
		for _, onlineUser := range oldState {
			if !contains(onlineUsers, onlineUser) {
				msg := tgbotapi.NewMessage(config.Telegrammchatid, fmt.Sprintf("%v hat Teamspeak verlassen ", onlineUser))
				tgBot.SendMessage(msg)
				time.Sleep(50 * time.Millisecond)
			}

		}
		oldState = onlineUsers
		newState = newState[:0]
		time.Sleep(60000 * time.Millisecond)
	}

}

func contains(slice []string, elem string) bool {
	for _, e := range slice {
		if e == elem {
			return true
		}
	}
	return false
}

func tsBot(conn *ts3.Conn, user string, passwd string) []string {
	connectionCommand := "login " + user + " " + passwd
	var cmds = []string{"version", connectionCommand, "use 1"}

	for _, s := range cmds {
		r, _ := conn.Cmd(s)
		fmt.Println("response:  ", r)
		time.Sleep(500 * time.Millisecond)
	}
	r, _ := conn.Cmd("clientlist")
	log.Println("Response: ", r)
	playerLine := strings.Split(r, "|")
	var toReturn []string
	for pl := range playerLine {
		if strings.Contains(playerLine[pl], "client_type=1") {
			log.Println("Skipping " + playerLine[pl])
			continue
		}
		parts := strings.Split(playerLine[pl], " ")
		for i := range parts {
			if strings.Contains(parts[i], "client_nickname") {
				user := strings.Split(parts[i], "=")[1]
				log.Println("Found a user ", user)
				toReturn = append(toReturn, user)
			}
		}
	}
	return toReturn
}
