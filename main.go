package main

import (
	"fmt"
	"log"
	"os/exec"
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Service struct {
	*Cache
	bot      *tgbotapi.BotAPI
	updates  tgbotapi.UpdatesChannel
	mu       sync.Mutex
	IOstatus bool
	Repositry
}
type Repositry struct {
	CSV string
}

type Cache struct {
	Users      map[int64]*User
	Datasets   map[int]*VoiceMessage
	SuperAdmin int64
	Admin      []int64
	Blocked    []int64
}

type User struct {
	ChatID       int64
	LastmsgID    int
	Status       bool
	Record       []int
	WaitingWords *int
}

type VoiceMessage struct {
	VoiceIndex int
	UserID     int64
	Confirmed  bool
}

func (s *Service) CreateUser(userID, chatID int64, msgID int) {

	s.mu.Lock()
	s.Users[userID] = &User{
		Status: false,
		ChatID: chatID,
	}
	s.Users[userID].LastmsgID = msgID
	s.mu.Unlock()
}

func main() {
	var service Service
	bot, err := tgbotapi.NewBotAPI("5005564686:AAGyPZX32onyXWCRGdkIq804LPmqBCgo3O0")
	if err != nil {
		log.Panic(err)
	}
	// bot.Debug = true
	// log.Printf("Authorized on account %s", bot.Self.UserName)

	UpdateConfi := tgbotapi.NewUpdate(0)
	UpdateConfi.Timeout = 60
	updates := bot.GetUpdatesChan(UpdateConfi)

	service.New(bot, updates)

	service.Start()
	bot.LogOut()
	time.Sleep(10 * time.Minute)
	service.UpdateUrSelf()

}

func (s *Service) New(bot *tgbotapi.BotAPI, updates tgbotapi.UpdatesChannel) *Service {
	s.bot = bot
	s.updates = updates
	s.Cache = &Cache{
		Users:      make(map[int64]*User),
		Datasets:   make(map[int]*VoiceMessage),
		SuperAdmin: 1034094796,
		Admin:      []int64{513685447, 1034094796},
		Blocked:    []int64{},
	}
	return s
}

func (s *Service) IsAdmin(userID int64) bool {
	for _, u := range s.Admin {
		if u == userID {
			return true
		}
	}
	return false
}

func (u *User) GetCallBack() (int64, int) {
	return u.ChatID, u.LastmsgID
}

func (u *User) UpdateMsg(msgID int) {
	u.LastmsgID = msgID
}

func (s *Service) NotificationBuilder(usersID int) {

}

func (s *Service) UpdateUserCache(userID, chatID int64, msgID int) {
	if _, found := s.Users[userID]; found {
		s.Users[userID].UpdateMsg(msgID)
		return
	}
	var NewUser *User = &User{
		Status: false,
		ChatID: chatID,
	}

	s.Users[userID] = NewUser
	s.Users[userID].LastmsgID = msgID
}

func (s *Service) Start() {
	for update := range s.updates {

		switch {
		case update.CallbackQuery != nil:
			switch {
			case update.CallbackQuery.Data == "8":
				return
			default:
				go s.CallbackQueryHandler(&update)
			}
		case update.InlineQuery != nil:
			go s.InlineQueryHandler(&update)
		case update.ChosenInlineResult != nil:
			go s.ChosenInlineResultHandler(&update)
		case update.Message != nil:
			go s.MessageHandler(&update)

		}

	}

}

func (s *Service) UpdateUrSelf() {
	cmd := exec.Command("git", "clone", "https://github.com/Zerihun-H/ServiceTestBot.git")
	cmd.Run()

	cmd = exec.Command("mv", "ServiceTestBot/serviceTest", ".")
	cmd.Run()

	cmd = exec.Command("rm", "-rf", "ServiceTestBot")
	cmd.Run()
	cmd = exec.Command("service", "restart", "serviceTest")
	cmd.Run()

	println("Start Update mode")
}

func (s *Service) Home(update *tgbotapi.Update) {
	if s.IsAdmin(update.Message.From.ID) {
		s.startDashboard(update.Message.Chat.ID)
		return
	}
	if u, found := s.Users[update.Message.From.ID]; found {
		s.DeleteOldMsg(update.Message.From.ID, update.Message.Chat.ID, update.Message.MessageID)
		u.UpdateMsg(s.startMenu(update.Message.Chat.ID))
		s.messageCleaner(update.Message.Chat.ID, update.Message.MessageID)
		return
	}
	s.RegisterUser(update)

}

func (s *Service) DeleteOldMsg(userID, chatID int64, msgID int) {
	if _, found := s.Users[userID]; found {
		s.messageCleaner(chatID, s.Users[userID].LastmsgID)
		return
	}
}

//Register can add user if not exits in cache

func (s *Service) RegisterUser(update *tgbotapi.Update) *User {

	var NewUser *User = &User{
		Status: false,
		ChatID: update.Message.Chat.ID,
	}

	s.Users[update.Message.From.ID] = NewUser
	s.Users[update.Message.From.ID].LastmsgID = s.startMenu(update.Message.Chat.ID)

	return NewUser
}

func (s *Service) startMenu(chatID int64) int {
	msg := tgbotapi.NewMessage(chatID, MainMessage, "", true)
	msg.ReplyMarkup = MainKeyBord
	rep, err := s.bot.Send(msg)
	if err != nil {
		fmt.Println(err)
	}
	return rep.MessageID
}

func (s *Service) startDashboard(chatID int64) int {
	msg := tgbotapi.NewMessage(chatID, "Welcome  Admin", "", true)
	msg.ReplyMarkup = AdminKeyBoard
	rep, err := s.bot.Send(msg)
	if err != nil {
		fmt.Println(err)
	}
	return rep.MessageID
}

func (s *Service) SendMessage(chatID int64, msgs string, PagePreview bool) (tgbotapi.Message, error) {
	msg := tgbotapi.NewMessage(chatID, msgs, "MarkdownV2", false)
	return s.bot.Send(msg)
}

func (s *Service) messageCleaner(chatID int64, msgIDs ...int) {
	for _, msgID := range msgIDs {
		deleteco := tgbotapi.NewDeleteMessage(chatID, msgID)
		s.bot.Send(deleteco)
	}
}
