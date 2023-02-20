package core

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"sync"
	"time"
)

type Server struct {
	Ip   string
	Port int

	OnlineMap map[string]*User
	Message   chan string

	MapLock sync.RWMutex
}

// NewServer
func NewServer(ip string, port int) *Server {
	return &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}
}

func (s *Server) handleConn(conn net.Conn) {
	user := NewUser(conn, s)

	user.Online()

	isLive := make(chan bool)
	isOffline := make(chan bool)
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			fmt.Println("buff len: ", n)
			if n == 0 { // 用户下线
				fmt.Println("offline")
				user.Offline()
				isOffline <- true
				return
			}

			if err != nil && err != io.EOF {
				fmt.Println("Conn read err ", err)
			}

			log.Println("Conn read msg ", string(buf))
			user.HandleMsg(buf)
			isLive <- true
		}
	}()
	// {"type":"rename","to":"","body":"eddie"}
	// {"type":"chat","to":"eddie1","body":"hello"}
	// {"type":"onlineList","to":"","body":" "}{"type":"onlineList"}
	for {
		select {
		case <-isOffline:
            return
		case <-isLive:
		case <-time.After(time.Second * 5): // 超时强踢
			user.SendMsg("timeout!")
			s.MapLock.Lock()
			user.conn.Close()
			close(user.C)
			delete(s.OnlineMap, user.Name)
			s.MapLock.Unlock()
			return
		}
	}
}

// BroadCast 消息广播
func (s *Server) BroadCast(user *User, msg string) {
	sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg
	s.Message <- sendMsg
}

// HandleMessage 推送消息
func (s *Server) HandleMessage() {
	for {
		msg := <-s.Message

		for _, user := range s.OnlineMap {
			user.C <- msg
		}
	}

}

// bind start
func (s *Server) Start() {
	lister, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.Ip, s.Port))
	if err != nil {
		log.Println("failed to listen", err)
		os.Exit(1)
	}
	log.Println("server start at ", lister.Addr().String())

	defer lister.Close()

	go s.HandleMessage()

	for {
		conn, err := lister.Accept()
		if err != nil {
			log.Println("failed to accept", err)
			continue
		}

		go s.handleConn(conn)
	}

}
