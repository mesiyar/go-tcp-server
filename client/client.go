package main

import (
	"fmt"
	"io"
	"log"
	"net"
)

type Client struct {
	serverIp   string
	serverPort int

	conn net.Conn

	C chan bool
}

func NewClient(serverIp string, serverPort int) *Client {
	client := &Client{
		serverIp:   serverIp,
		serverPort: serverPort,
		C:          make(chan bool),
	}

	cli, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, serverPort))
	if err != nil {
		log.Fatal("dail to server failed ", err)
		return nil
	}
	client.conn = cli
	go client.send()
	go client.handleMessage()
	return client

}

func (c *Client) send() {
	var body string
	for {
		_, err := fmt.Scanln(&body)
		if err != nil {
			log.Println("get message failed", err)
			continue
		}
		_, err = c.conn.Write([]byte(body))
		if err != nil && err != io.EOF {
			log.Println("write message failed", err)
			continue
		}
	}
}

func (c *Client) handleMessage() {
	for {
		buf := [512]byte{}
		n, err := c.conn.Read(buf[:])
		if err != nil && err != io.EOF {
			fmt.Println("recv failed, err:", err)
			return
		}
		if err == io.EOF {
			c.C <- true
			return
		}
		fmt.Println(string(buf[:n]))
	}
}

func main() {
	c := NewClient("127.0.0.1", 8099)
	for {
		select {
		case <-c.C:
			fmt.Println("client closed")
			return
		}
	}
}
