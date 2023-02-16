package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strings"
)

type User struct {
	Name string
	Addr string
	C    chan string
	conn net.Conn

	server *Server
}

//NewUser
func NewUser(conn net.Conn, s *Server) *User {
	name := conn.RemoteAddr().String()
	addr := name

	user := &User{
		Name:   name,
		Addr:   addr,
		conn:   conn,
		C:      make(chan string),
		server: s,
	}

	go user.ListenMessage()

	return user
}

// ListenMessage 监听消息
func (u *User) ListenMessage() {
	for {
		msg := <-u.C
		_, err := u.conn.Write([]byte(msg + "\r\n"))
		if err != nil {
			log.Printf("send to user %s failed error %v", u.Name, err)
		}

	}
}

// Online
func (u *User) Online() {
	u.server.MapLock.Lock()
	defer u.server.MapLock.Unlock()
	u.server.OnlineMap[u.Name] = u

	u.server.BroadCast(u, "online !")
}

// Offline
func (u *User) Offline() {
	u.server.MapLock.Lock()
	defer u.server.MapLock.Unlock()
	delete(u.server.OnlineMap, u.Name)

	u.server.BroadCast(u, "offline !")
}

// HandleMsg
func (u *User) HandleMsg(msg []byte) {
	index := bytes.IndexByte(msg, 0)
	msg = msg[:index]
	m := &Message{}
	err := json.Unmarshal(msg, m)
	if err != nil {
		log.Printf("Unmarshal message error %v", err)
		return
	}

	if m.Type == MsgTypeOnlineList { // 在线列表
		u.getOnlineList()
	} else if m.Type == MsgTypeChat { // 聊天
		u.handleChat(m)
	} else if m.Type == MsgTypeRename { // 改名
		u.rename(m.Body)
	} else {
		u.server.BroadCast(u, m.Body)
	}
}

// SendMsg send msg to current connection
func (u *User) SendMsg(msg string) {
	_, err := u.conn.Write([]byte(msg))
	if err != nil {
		fmt.Println("send failed")
	}
}

// rename
func (u *User) rename(name string) {
	u.server.MapLock.Lock()
	defer u.server.MapLock.Unlock()
	if _, ok := u.server.OnlineMap[name]; ok {
		u.SendMsg("change name failed " + name + " already exists")
		return
	}
	delete(u.server.OnlineMap, u.Name)
	u.Name = name
	u.server.OnlineMap[name] = u
	u.SendMsg("change name to " + name + " successfully")
}

// sendOnlineList
func (u *User) getOnlineList() {
	u.server.MapLock.Lock()
	defer u.server.MapLock.Unlock()
	ol := make([]string, len(u.server.OnlineMap))
	for i := range u.server.OnlineMap {
		ol = append(ol, u.server.OnlineMap[i].Name)
	}
	u.SendMsg(strings.Join(ol, "\n"))
}

// handleChat 处理聊天信息
func (u *User) handleChat(m *Message) {
	u.server.MapLock.Lock()
	defer u.server.MapLock.Unlock()
	if m.To == "" {
		u.server.BroadCast(u, m.Body)
		return
	}
	_, ok := u.server.OnlineMap[m.To]
	if !ok {
		u.SendMsg("user " + m.To + "is not online !")
		return
	}
	msg :=
	u.server.OnlineMap[m.To].SendMsg(m.Body)
}
