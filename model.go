package main

import (
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type UpFunc func(*tgbotapi.Update)

var token string = "5005564686:AAEYc8u0SIeKE6M2IZH53TyTCs-AT2V7yhc"

type Pair struct {
	Key   int64
	Value float32
}

type RankList []Pair

func (p RankList) Len() int           { return len(p) }
func (p RankList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p RankList) Less(i, j int) bool { return p[i].Value > p[j].Value }

type Service struct {
	*Cache
	bot     *tgbotapi.BotAPI
	updates tgbotapi.UpdatesChannel
	mu      sync.RWMutex
}
type Cache struct {
	Users        map[int64]*User
	WordList     []string
	SuperAdmin   int64
	Admin        []int64
	VoiceVersion string
	Blocked      map[int64]int
	RankList
	HighestInvention int
	Request          map[int64]int
	Contribution     map[int64]int
}

type User struct {
	LastmsgID     int
	InvitedBy     int64
	Invited       []int64
	Record        []int
	Rejected      int
	Confirmed     int
	Datasets      []*VoiceMessage
	RecordPointer int
	PhoneNum      string
	Rank          int
	Avarage       float32
	RequestTime   int64
}

func (u *User) Restart() {
	u.Record = nil
	u.RecordPointer = 0
}

type VoiceMessage struct {
	MsgID      int
	VoiceIndex int
	Confirmed  bool
	FileID     string
}
