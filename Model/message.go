package Model

import (
	"fmt"
	"hash/crc32"
	"strconv"
	"time"
)

type Message struct {
	MegID    string `gorm:"PRIMARY KEY"`
	Title    string `gorm:"NOT NULL" json:"title"`
	Content  string `gorm:"NOT NULL" json:"sendMessage"`
	SendTime string `gorm:"NOT NULL"`
	Read     bool   `gorm:"DEFAULT 0"`
	//被通知人Id
	NotifiedID string `gorm:"NOT NULL" json:"toId"`
	//通知人Id
	NotifierID string `gorm:"NOT NULL" json:"fromId"`
}

type SendedMsg struct {
	ToId        string `json:"toId"`
	SendTime    string `json:"sendTime"`
	SendMessage string `json:"sendMessage"`
	MessageId   string `json:"messageId"`
	Title       string `json:"title"`
	Key         string `json:"key"`
	HasRead     bool   `json:"hasRead"`
}

type RecvMsg struct {
	FromId      string `json:"fromId"`
	SendTime    string `json:"sendTime"`
	SendMessage string `json:"sendMessage"`
	MessageId   string `json:"messageId"`
	Title       string `json:"title"`
	HasRead     bool   `json:"hasRead"`
	Key         string `json:"key"`
}

func NewMessage(u1Id, u2Id string) *Message {
	return &Message{
		MegID:      strconv.Itoa(int(crc32.ChecksumIEEE([]byte(u1Id + u2Id + time.Now().String())))),
		SendTime:   GetSystemTime(),
		NotifiedID: u1Id,
		NotifierID: u2Id,
	}
}

func GetRecvMessage(u User) []Message {
	var m []Message
	GlobalConn.Where(&Message{NotifiedID: u.Id}).Find(&m)
	return m
}

func GetSendMessage(u User) []Message {
	var m []Message
	GlobalConn.Where(&Message{NotifierID: u.Id}).Find(&m)
	return m
}

func GetMessageById(MId string) *Message {
	var m Message
	result := GlobalConn.Where(&Message{MegID: MId}).First(&m)
	if result.Error != nil || result.RowsAffected == 0 {
		return nil
	}
	return &m
}

func (m *Message) CopySended() SendedMsg {
	var Sended SendedMsg
	Sended.MessageId = m.MegID
	Sended.Key = m.MegID
	Sended.SendMessage = m.Content
	Sended.Title = m.Title
	Sended.ToId = m.NotifiedID
	Sended.SendTime = m.SendTime
	Sended.HasRead = m.Read
	return Sended
}

func (m *Message) CopyRecv() RecvMsg {
	var Recv RecvMsg
	Recv.MessageId = m.MegID
	Recv.Key = m.MegID
	Recv.Title = m.Title
	Recv.SendMessage = m.Content
	Recv.SendTime = m.SendTime
	user := GetUserById(m.NotifierID)
	Recv.FromId = fmt.Sprintf("%s(%s)", user.Name, m.NotifierID)
	Recv.HasRead = m.Read
	return Recv
}

func (m *Message) GetContent(con string) {
	m.Content = con
}

func (m *Message) AddMessage() {
	GlobalConn.Create(m)
}

func (m *Message) ReadMessage() {
	m.Read = true
	GlobalConn.Model(&Message{}).Where(&Message{MegID: m.MegID}).Update("read", m.Read)
}

func (m *Message) DeleteMessage() {
	GlobalConn.Delete(m)
}
