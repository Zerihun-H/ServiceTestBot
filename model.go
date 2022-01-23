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

type AdminMenuState uint8
type MenuState uint8

const (
	Admin           = 1034094796
	DatasetGroup    = "-1001717101880"
	DatasetGroupInt = -1001717101880
	BackupGroup     = "-1001615301597"
)

const (
	HomePage AdminMenuState = iota
	BlockPage
	UnBlockPage
	ChangeAdmin
)

const (
	UserHomePage MenuState = iota
	UserMenuPage
	WaitingVoice
	WaitingVerification
)

type RankList []*Pair

func (p RankList) Len() int           { return len(p) }
func (p RankList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p RankList) Less(i, j int) bool { return p[i].Value > p[j].Value }

type Service struct {
	*Cache
	bot     *tgbotapi.BotAPI
	updates tgbotapi.UpdatesChannel
	mu      sync.RWMutex
}
type Cache2 struct {
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

type Cache struct {
	Users        map[int64]*User
	WordList     []string
	Admin        int64
	VoiceVersion string
	Blocked      map[int64]int
	RankList
	HighestInvitation     float32
	HighestTreeRecord     float32
	HighestTreeConfirmed  float32
	HighestTreeInvitation float32
	Request               map[int64]int
	AdminMenuState
	ColleagueContribution map[int64]int
}
type User struct {
	LastmsgID       int
	InvitedBy       int64
	Invited         []int64
	Record          []int
	Confirmed       int
	Datasets        []*VoiceMessage
	TreeConfirmed   float32
	TreeRecord      float32
	TreeInvitation  float32
	RecordPointer   int
	PhoneNum        string
	Rank            int
	Avarage         float32
	RequestTime     int64
	VerifiedSetting bool
	MenuState
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
