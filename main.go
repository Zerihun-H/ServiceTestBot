package main

import (
	"encoding/gob"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type UpFunc func(*tgbotapi.Update)

var token string = "5005564686:AAGyPZX32onyXWCRGdkIq804LPmqBCgo3O0"

type Service struct {
	*Cache
	bot     *tgbotapi.BotAPI
	updates tgbotapi.UpdatesChannel
	mu      sync.RWMutex
}
type Cache struct {
	Users      map[int64]*User
	WordList   []string
	SuperAdmin int64
	Admin      []int64
	Blocked    map[int64]int
}

type User struct {
	LastmsgID     int
	InvitedBy     int64
	Invited       []int64
	Record        []int
	Rejected      int32
	Confirmed     int32
	Datasets      []*VoiceMessage
	RecordPointer int
	PhoneNum      string
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

func (s *Service) CreateUser(userID, inviterID int64, msgID int) {

	s.mu.Lock()
	s.Users[userID] = &User{
		InvitedBy: inviterID,
	}
	s.Users[userID].LastmsgID = msgID
	s.mu.Unlock()
}

func main() {

	var service Service
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatal(err)
	}

	// bot.Debug = true
	// log.Printf("Authorized on account %s", bot.Self.UserName)

	UpdateConfi := tgbotapi.NewUpdate(0)
	UpdateConfi.Timeout = 60
	updates := bot.GetUpdatesChan(UpdateConfi)

	service.New(bot, updates)

	go service.Doctor(60)
	service.Start()
	bot.LogOut()
	time.Sleep(10 * time.Minute)
	service.UpdateUrSelf()
}

func (s *Service) Start() {
	for update := range s.updates {
		switch {
		case update.CallbackQuery != nil:
			if update.CallbackQuery.Data == "8" {
				return
			}
			go s.HandleUpdate(s.CallbackQueryHandler, &update, update.CallbackQuery.From.ID, update.CallbackQuery.Message.MessageID)
		case update.InlineQuery != nil:
			go s.HandleUpdate(s.InlineQueryHandler, &update, update.InlineQuery.From.ID, 0)
		case update.ChosenInlineResult != nil:
			go s.HandleUpdate(s.ChosenInlineResultHandler, &update, update.ChosenInlineResult.From.ID, 0)
		case update.Message != nil:
			go s.MessageHandler(&update)
		}
	}
}

func (s *Service) AskCache() {
	var FileID string
	s.ReportToAdmin("please send to me a last updated Cache")
	for update := range s.updates {
		if update.Message != nil {
			if s.SuperAdmin == update.Message.From.ID {
				if update.Message.Text == "!" || update.Message.Command() == "start" {
					s.messageCleaner(update.Message.Chat.ID, update.Message.MessageID)
					return
				}
				if update.Message.Document != nil {
					FileID = update.Message.Document.FileID
					break
				}
			}
			s.messageCleaner(update.Message.Chat.ID, update.Message.MessageID)
		}
	}
	if fileURL, err := s.bot.GetFileDirectURL(FileID); err == nil {
		if err = s.DownloadFile("cache.gob", fileURL); err == nil {
			s.LoadBackup("cache")
			s.ReportToAdmin("successfully Updated Cache")
			return
		}
		s.ReportToAdmin(err.Error())
	}
}

func (s *Service) New(bot *tgbotapi.BotAPI, updates tgbotapi.UpdatesChannel) *Service {
	s.bot = bot
	s.updates = updates
	s.Cache = &Cache{
		Users:      make(map[int64]*User),
		SuperAdmin: 1034094796,
		//513685447
		WordList: WordList,
		Admin:    []int64{1034094796},
		Blocked:  make(map[int64]int),
	}

	switch {
	case s.Exists("cache.gob"):
		s.LoadBackup("cache")
	default:
		s.AskCache()
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

func (u *User) GetCallBack() int {
	return u.LastmsgID
}

func (u *User) UpdateMsg(msgID int) {
	u.LastmsgID = msgID
}

func (s *Service) NotificationBuilder(usersID int) {}

func (s *Service) UpdateUserOldMsg(userID int64, msgID int) {
	if _, found := s.Users[userID]; !found {
		s.CreateUser(userID, 0, msgID)
	}

	s.Users[userID].UpdateMsg(msgID)
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

}

func (s *Service) Home(update *tgbotapi.Update) {
	var userID, msgID = update.Message.From.ID, update.Message.MessageID
	if s.IsAdmin(userID) {
		s.startDashboard(userID)
		return
	}
	if s.IsBlock(userID) {
		s.MakeNotice(BlockNotice, userID, msgID)
		return
	}
	if u, found := s.Users[update.Message.From.ID]; found {
		s.DeleteOldMsg(userID, msgID)
		u.UpdateMsg(s.startMenu(userID))
		s.messageCleaner(userID, msgID)
		return
	}

	var inviter int64

	messages := strings.Split(update.Message.Text, " ")
	if len(messages) > 1 {
		inviter, _ = strconv.ParseInt(messages[1], 10, 64)
	}
	if inviter != 0 && inviter != userID {
		s.AddInvitation(inviter, userID)
	}

	s.CreateUser(userID, inviter, msgID)
	msgIDs := s.startMenu(userID)
	s.UpdateUserOldMsg(userID, msgIDs)

}

func (s *Service) DeleteOldMsg(userID int64, msgID int) {
	if _, found := s.Users[userID]; found {
		s.messageCleaner(userID, s.Users[userID].LastmsgID)
		return
	}
}

func (s *Service) startMenu(usrID int64) int {
	msg := tgbotapi.NewMessage(usrID, MainMessage, "MarkdownV2", false)
	msg.ReplyMarkup = MainKeyBord
	rep, _ := s.bot.Send(msg)
	return rep.MessageID
}

func (s *Service) startDashboard(usrID int64) int {
	msg := tgbotapi.NewMessage(usrID, "Welcome  Admin", "", true)
	msg.ReplyMarkup = AdminKeyBoard
	rep, err := s.bot.Send(msg)
	if err != nil {
		s.ReportToAdmin(err.Error())
	}
	return rep.MessageID
}

func (s *Service) SendMessage(chatID int64, msgs, parse string, PagePreview bool) (tgbotapi.Message, error) {
	msg := tgbotapi.NewMessage(chatID, msgs, parse, PagePreview)
	return s.bot.Send(msg)
}

func (s *Service) messageCleaner(chatID int64, msgIDs ...int) {
	for _, msgID := range msgIDs {
		deleteco := tgbotapi.NewDeleteMessage(chatID, msgID)
		s.bot.Send(deleteco)
	}
}

func (s *Service) UserMention(name, username string, userID int64) string {
	msgText := fmt.Sprintf("[[Mention user](tg://user?id=%d)](https://t.me/%s)", userID, username)

	// msgText := fmt.Sprintf("[[%s](tg://user?id=%d)](https://t.me/%s)-%d", name, userID, username, userID)
	return msgText
}

func (s *Service) MessageEntityBuidler(usr *tgbotapi.User, word string, userName ...string) ([]tgbotapi.MessageEntity, string) {
	var messageEntity []tgbotapi.MessageEntity
	var name string
	var username string

	word = "#" + word
	userID := strconv.FormatInt(usr.ID, 10)

	switch {
	case len(userName) > 0:
		username = userName[0]
	default:
		username = "@" + usr.UserName
	}

	switch {
	case len([]rune(usr.FirstName)) < 2 || strings.Contains(usr.FirstName, "ㅤ"):
		name = "Anonymous"
	default:
		name = usr.FirstName
	}

	wordLen, nameLen, _, _ := len([]rune(word)), len([]rune(name)), len(userID), len([]rune(username))

	hashEntity := tgbotapi.MessageEntity{
		Type:   "hashtag",
		Offset: 0,
		Length: wordLen,
	}
	nameEntity := tgbotapi.MessageEntity{
		Type:   "text_mention",
		Offset: wordLen + 1,
		Length: nameLen,
		User:   usr,
	}

	messageEntity = append(messageEntity, hashEntity, nameEntity)

	var text string
	switch {
	case username != "":
		text = word + "\n" + name + "\n" + userID + "\n" + username
	default:
		text = word + "\n" + name + "\n" + userID
	}

	return messageEntity, text
}

func (s *Service) IsBlock(usrID int64) bool {
	if _, found := s.Blocked[usrID]; found {
		return true
	}
	return false
}

func (s *Service) HandleUpdate(nextFunc UpFunc, update *tgbotapi.Update, userID int64, msgID int) {
	switch {
	case s.IsBlock(userID):
		go s.MakeNotice(BlockNotice, userID, msgID)
	default:
		go nextFunc(update)
	}
}

func (s *Service) AddInvitation(inviter, userID int64) {
	if _, found := s.Users[inviter]; !found {
		s.CreateUser(inviter, 0, 0)
	}
	s.Users[inviter].Invited = append(s.Users[inviter].Invited, userID)
}

func (s *Service) CountInvitation(inviter int64) int {
	if _, found := s.Users[inviter]; !found {
		s.CreateUser(inviter, 0, 0)
		return 0
	}

	return len(s.Users[inviter].Invited)
}

func (s *Service) MakeNotice(msgText string, userID int64, msgID int) {
	msg := tgbotapi.NewMessage(userID, msgText, "MarkdownV2", true)
	rep, err := s.bot.Send(msg)
	s.ReportToAdmin(fmt.Sprintf("User ID %d Blocked", userID))
	if err != nil {
		s.ReportToAdmin(err.Error())
	}

	if msgID != 0 {
		s.DeleteOldMsg(userID, msgID)
	}

	s.UpdateUserOldMsg(userID, rep.MessageID)
	s.messageCleaner(userID, msgID)
}

func (s *Service) BackupCache(filename string) {
	var file *os.File
	var err error
	if file, err = os.Create(fmt.Sprintf("%s.gob", filename)); err != nil {
		s.ReportToAdmin(err.Error())
	}
	defer file.Close()
	encoder := gob.NewEncoder(file)
	encoder.Encode(s.Cache)
}

func (s *Service) LoadBackup(filename string) {
	var file *os.File
	var err error
	if file, err = os.Open(fmt.Sprintf("%s.gob", filename)); err != nil {
		s.ReportToAdmin(err.Error())
	}
	defer file.Close()
	decoder := gob.NewDecoder(file)
	decoder.Decode(s.Cache)
}

func (s *Service) Doctor(sec int64) {
	for {
		time.Sleep(time.Duration(sec) * time.Second)
		s.BackupCache("cache")
		fileBackUp := tgbotapi.NewDocumentChannel("-1001777959481", tgbotapi.FilePath("cache.gob"))
		if _, err := s.bot.Send(fileBackUp); err != nil {
			s.ReportToAdmin(err.Error())
		}
	}
}

func (s *Service) ReportToAdmin(msgText string) {
	if _, err := s.SendMessage(s.SuperAdmin, msgText, "", false); err != nil {
		fmt.Printf("errors %s", err.Error())
	}

}

func (s *Service) DownloadFile(filepath string, url string) (err error) {

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) Exists(name string) bool {
	_, err := os.Stat(name)
	if err == nil {
		return true
	}
	if errors.Is(err, os.ErrNotExist) {
		return false
	}
	return false
}
