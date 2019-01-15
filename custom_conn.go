package whatsapp

import (
	"sync"
	"time"
)

/*
Creates a new connection with a given timeout. The websocket connection to the WhatsAppWeb servers get´s established.
The goroutine for handling incoming messages is started
*/
func NewCustomConn(timeout time.Duration, longClientName, shortClientName string) (*Conn, error) {
	wac := &Conn{
		wsConn:        nil, // will be set in connect()
		wsConnMutex:   sync.RWMutex{},
		listener:      make(map[string]chan string),
		listenerMutex: sync.RWMutex{},
		writeChan:     make(chan wsMsg),
		handler:       make([]Handler, 0),
		msgCount:      0,
		msgTimeout:    timeout,
		Store:         newStore(),

		longClientName:  longClientName,
		shortClientName: shortClientName,
	}

	if err := wac.connect(); err != nil {
		return nil, err
	}

	go wac.readPump()
	go wac.writePump()
	go wac.keepAlive(20000, 90000)

	return wac, nil
}
