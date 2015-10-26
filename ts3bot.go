package main

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/toqueteos/ts3"
	"strings"
	"time"
)

type tomlConfig struct {
	Tsuser   string
	Tspasswd string
	Tsurl    string
}

func main() {
	var config tomlConfig
	if _, err := toml.DecodeFile("config.toml", &config); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Connecting to TS3 Server at: %s \n", config.Tsurl)
	conn, err := ts3.Dial(config.Tsurl, true)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Connected!")
	defer conn.Close()

	bot(conn, config.Tsuser, config.Tspasswd)
}

func bot(conn *ts3.Conn, user string, passwd string) {
	defer conn.Cmd("quit")
	connectionCommand := "login " + user + " " + passwd
	var cmds = []string{"version", connectionCommand, "use 1"}

	for _, s := range cmds {
		r, _ := conn.Cmd(s)
		fmt.Println("response:  ", r)
		time.Sleep(500 * time.Millisecond)
	}
	r, _ := conn.Cmd("clientlist")
	fmt.Println("Response: ", r)
	playerLine := strings.Split(r, "|")
	for pl := range playerLine {
		parts := strings.Split(playerLine[pl], " ")
		for i := range parts {
			if strings.Contains(parts[i], "client_nickname") {
				user := strings.Split(parts[i], "=")[1]
				fmt.Println("Found a user ", user)
			}
		}
	}
}
