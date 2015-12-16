package tsstatus

import (
	"github.com/toqueteos/ts3"
	"log"
	"time"
)

type Event struct {
	Typ  string
	User string
}

type TsStatus struct {
	conn   *ts3.Conn
	users  map[int]string
	events chan Event
}

func NewStatus(host, user, pwd string) *TsStatus {
	conn, err := ts3.Dial(host, true)
	if err != nil {
		panic(err)
	}
	s := &TsStatus{
		conn:   conn,
		users:  make(map[int]string),
		events: make(chan Event),
	}

	s.start(conn, user, pwd)

	return s
}

func (s *TsStatus) start(conn *ts3.Conn, user, passwd string) {
	connectionCommand := "login " + user + " " + passwd
	r, _ := conn.Cmd(connectionCommand)
	log.Println("Login Response: ", r)
	time.Sleep(500 * time.Millisecond)
	r, _ = conn.Cmd("use 1")
	time.Sleep(500 * time.Millisecond)
	r, _ = conn.Cmd("servernotifyregister event=server ")
	conn.NotifyFunc(func(x, y string) {
		log.Println("Notify: " + x)
		event := Event{
			Typ:  x,
			User: y,
		}
		s.events <- event
	})
}

func (s *TsStatus) GetChan() chan Event {
	return s.events
}
