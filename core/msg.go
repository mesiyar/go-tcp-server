package core

type Message struct {
	Type string `json:"type"`
	To   string `json:"to"`
	Body string `json:"body"`
}

const MsgTypeOnlineList = "onlineList"
const MsgTypeChat = "chat"
const MsgTypeRename = "rename"

type SendMsg struct {
	Type string `json:"type"`
	From string `json:"from"`
	Body string `json:"body"`
}
