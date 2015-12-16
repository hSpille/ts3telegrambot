package main

import (
	"encoding/json"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/Syfaro/telegram-bot-api"
	"github.com/toqueteos/ts3"
	"log"
	"os"
	"strings"
	"time"
	"ts3bot/tsstatus"
)

type tomlConfig struct {
	Tsuser          string
	Tspasswd        string
	Tsurl           string
	Telegrammapikey string
	Telegrammchatid int
	Dbfilename      string
}

var config tomlConfig

func main() {
	if _, err := toml.DecodeFile("config.toml", &config); err != nil {
		fmt.Println(err)
		return
	}
	tgBot, err := tgbotapi.NewBotAPI(config.Telegrammapikey)
	if err != nil {
		log.Panic(err)
	}
	tgBot.Debug = true
	log.Printf("Telegramm Bot authorized on account %s", tgBot.Self.UserName)
	log.Printf("Connecting to TS3 Server at: %s \n", config.Tsurl)
	tsConn, err := ts3.Dial(config.Tsurl, true)
	if err != nil {
		panic(err)
	}
	//logins := make(map[string]time.Time)

	defer tsConn.Cmd("quit")
	defer tsConn.Close()
	//var oldState []string
	var state []string
	status := tsstatus.NewStatus(config.Tsurl, config.Tsuser, config.Tspasswd)
	c := status.GetChan()
	for event := range c {
		log.Println("Event-Typ:", event.Typ)
		if strings.Compare("notifycliententerview", event.Typ) == 0 {
			log.Println("Player joined")
			parts := strings.Split(event.User, " ")
			for _, part := range parts {
				if strings.Contains(part, "client_nickname=") {
					log.Println("Part:" + part)
					user := strings.Split(part, "=")[1]
					log.Println("Found a user ", user)
					state = append(state, user)
				}
			}
		}
		if strings.Compare("notifyclientleftview", event.Typ) == 0 {
			log.Println("Player joined")
			parts := strings.Split(event.User, " ")
			for _, part := range parts {
				if strings.Contains(part, "client_nickname=") {
					log.Println("Part:" + part)
					user := strings.Split(part, "=")[1]
					log.Println("Found a user ", user)
					//remove player!
				}
			}
		}
		log.Println("Online currently: ", state)
	}
	// for {
	// 	onlineUsers := tsBot(tsConn, config.Tsuser, config.Tspasswd)
	// 	for _, onlineUser := range onlineUsers {
	// 		if !contains(oldState, onlineUser) {
	// 			msg := tgbotapi.NewMessage(config.Telegrammchatid, fmt.Sprintf("%v im Teamspeak!", onlineUser))
	// 			logins[onlineUser] = time.Now()
	// 			tgBot.Send(msg)
	// 			time.Sleep(50 * time.Millisecond)
	// 			newState = append(newState, onlineUser)
	// 		}
	// 	}
	// 	for _, onlineUser := range oldState {
	// 		if !contains(onlineUsers, onlineUser) {
	// 			duration := time.Since(logins[onlineUser])
	// 			storeTime(onlineUser, duration)
	// 			msg := tgbotapi.NewMessage(config.Telegrammchatid, fmt.Sprintf("%v hat Teamspeak verlassen nach %v", onlineUser, duration))
	// 			tgBot.Send(msg)
	// 			time.Sleep(50 * time.Millisecond)
	// 		}

	// 	}
	// 	oldState = onlineUsers
	// 	newState = newState[:0]
	// 	time.Sleep(60000 * time.Millisecond)
	// }

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

func storeTime(user string, duration time.Duration) {
	totalTimes := readFile()

	var err error

	if existsFile("db.json") {
		err = os.Remove("db.json")
		if err != nil {
			log.Panic(err)
		}
	}

	file, err := os.OpenFile("db.json", os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		log.Panic(err)
	}
	enc := json.NewEncoder(file)
	timeInDb, ok := totalTimes[user]
	if !ok {
		timeInDb = time.Time{}
	}
	totalTimes[user] = timeInDb.Add(duration)
	err = enc.Encode(totalTimes)
	if err != nil {
		log.Panic(err)
	}
}

func existsFile(filename string) bool {
	_, err := os.Stat(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
		log.Panic(err)
	}
	return true
}

/*
 * Liefert leere map, wenn keine Datei gefunden.
 */
func readFile() map[string]time.Time {
	totalTimes := make(map[string]time.Time)
	file, err := os.Open("db.json")
	if err != nil {
		if os.IsNotExist(err) {
			return totalTimes
		}
		log.Panic(err)
	}
	defer file.Close()
	dec := json.NewDecoder(file)
	err = dec.Decode(&totalTimes)
	if err != nil {
		log.Println("DB Decode fehlerhaft: ", err)
	}
	return totalTimes
}
