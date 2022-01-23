package main

import (
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (s *Service) CreateUser(userID, inviterID int64, msgID int) {

	s.mu.Lock()
	s.Users[userID] = &User{
		InvitedBy:       inviterID,
		MenuState:       UserHomePage,
		VerifiedSetting: true,
	}
	s.Users[userID].LastmsgID = msgID
	s.mu.Unlock()

	if inviterID != 0 {
		go s.UpdateTreeInvitation(userID)
	}
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
	go service.Dispatcher(15)
	go service.Leaderboard(60)

	// service.MAkeAllWaiting()
	// service.RemoveAllBlock()
	// service.UpdateTreeSchema()

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

func (s *Service) MAkeAllWaiting() {
	for i := range s.Users {
		s.Users[i].VerifiedSetting = true
		s.Users[i].MenuState = WaitingVoice
	}
}

func (s *Service) RemoveAllBlock() {
	s.Blocked = make(map[int64]int)
}

func (s *Service) HandleUpdate(nextFunc UpFunc, update *tgbotapi.Update, userID int64, msgID int) {
	var requested = !s.Requested(userID, msgID)
	switch {
	case s.IsBlock(userID):
		go s.MakeNotice(BlockNotice, userID, msgID)
	case requested && update.CallbackQuery != nil:
		go s.Callback(update.CallbackQuery.ID, TooManyMessage, true)
	case requested:
		go s.MakeNotice(TooManyMessage, userID, msgID)
	default:
		go nextFunc(update)
	}
}

func (s *Service) New(bot *tgbotapi.BotAPI, updates tgbotapi.UpdatesChannel) *Service {
	s.bot = bot
	s.updates = updates
	s.Cache = &Cache{
		Users:                 make(map[int64]*User),
		WordList:              WordList,
		Admin:                 Admin,
		VoiceVersion:          "",
		Blocked:               make(map[int64]int),
		RankList:              []*Pair{},
		Request:               make(map[int64]int),
		ColleagueContribution: make(map[int64]int),
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
	return s.Admin == userID
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
		s.DeleteOldMsg(userID)
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

func (s *Service) DeleteOldMsg(userID int64) {
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
	msg.ReplyMarkup = AdminMenuKeyBoard
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

func (s *Service) AddInvitation(inviter, userID int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, found := s.Users[inviter]; !found {
		s.CreateUser(inviter, 0, 0)
	}
	s.Users[inviter].Invited = append(s.Users[inviter].Invited, userID)
	lenInvited := float32(len(s.Users[inviter].Invited))
	if lenInvited > s.HighestInvitation {
		s.HighestInvitation = lenInvited
	}
}

func (s *Service) CountInvitation(inviter int64) int {
	if _, found := s.Users[inviter]; !found {
		s.CreateUser(inviter, 0, 0)
		return 0
	}

	return len(s.Users[inviter].Invited)
}

func (s *Service) ConfirmUserDataSet(userID int64, msgID int) {
	if user, found := s.Users[userID]; found {
		for v := range user.Datasets {
			if user.Datasets[v].MsgID == msgID {
				if !user.Datasets[v].Confirmed {
					s.mu.Lock()
					user.Datasets[v].Confirmed = true
					user.Confirmed++
					s.mu.Unlock()
					s.UpdateTreeConfirmed(user.InvitedBy)
				}
			}
		}
	}
}

func (s *Service) RejectUserDataset(userID int64, msgID int) {
	if u, found := s.Users[userID]; found {
		for v := range u.Datasets {
			if u.Datasets[v].MsgID == msgID {
				if u.Datasets[v].Confirmed {
					s.mu.Lock()
					u.Confirmed--
					u.Datasets[v].Confirmed = false
					s.mu.Unlock()
				}
			}
		}
	}
}
