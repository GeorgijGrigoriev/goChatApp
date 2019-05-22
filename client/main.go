package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	"golang.org/x/net/websocket"
)

type Message struct {
	Text string `json:"text"`
}

var (
	port = flag.String("port", "9000", "port used for ws connection")
)

func connect() (*websocket.Conn, error) {
	return websocket.Dial(fmt.Sprintf("ws://localhost:%s", *port), "", "http://"+localIP())
}

func localIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Println(err.Error())
	}

	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	localIP := localAddr.IP.String()
	return localIP
}

func main() {
	flag.Parse()

	ws, err := connect()

	if err != nil {
		log.Fatal(err)
	}

	defer ws.Close()

	var m Message

	go func() {
		for {
			err := websocket.JSON.Receive(ws, &m)
			if err != nil {
				fmt.Println("Ошибка доставки сообщения: ", err.Error())
				break
			}
			fmt.Println("Message: ", m)

		}
	}()
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := scanner.Text()
		if text == "" {
			continue
		}
		m := Message{
			Text: text,
		}
		err = websocket.JSON.Send(ws, m)
		if err != nil {
			fmt.Println("Ошибка доставки сообщения: ", err.Error())
			break
		}
	}
}
