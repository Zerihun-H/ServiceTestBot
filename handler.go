package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type CallbackQueryHandler func(*tgbotapi.Update)

//Hnadler Trigger Callback from msg or inline
func (s *Service) CallbackQueryHandler(update *tgbotapi.Update) {
	switch update.CallbackQuery.Data {
	case "1":
		go s.VoiceRequestHandler(update)
	case "-1":
		go s.AnswerCallbackQuery(s.GoEnd, update, "ሌላ")
	case "0":
		go s.CloseVoiceRequest(update)
	case "-2":
		go s.AnswerCallbackQuery(s.GoBack, update, "ወደኋላ")
	case "2":
		go s.AnswerCallbackQuery(s.GoNext, update, "ቀጣይ")
	case "3":
		go s.RestartMenu(update)
	case "-3":
		go s.VoiceRequestHandler(update)
	case "4":
		go s.Profile(update)
	case "5":
		go s.ConfirmDataSet(update)
	case "-5":
		go s.RejectDataset(update)
	case "6":
		go s.BlockUser(update)
		go s.RejectDataset(update)
	case "7":
		go s.MoveBack(update)
	}

}

func (s *Service) Requested(userID int64, msgID int) bool {
	var user *User
	var found bool

	if user, found = s.Users[userID]; !found {
		s.CreateUser(userID, 0, msgID)
		s.Request[userID] = s.Request[userID] + 1
		s.Users[userID].RequestTime = time.Now().UnixNano()
		return true
	}
	if s.Request[userID]+1 > 10 {
		s.Blocked[userID] = msgID
		s.MakeNotice(BlockNotice, userID, msgID)
		s.ReportToAdmin(fmt.Sprintf("User ID [%d](tg://user?id=%d) Blocked by LamabaBot", userID, userID))
		return false
	}

	timeDiff := time.Since(time.Unix(0, user.RequestTime))

	if timeDiff.Seconds() < 1.5 {
		s.Request[userID] = s.Request[userID] + 1
		s.Users[userID].RequestTime = time.Now().UnixNano()
		return false
	}
	s.Request[userID] = s.Request[userID] + 1
	s.Users[userID].RequestTime = time.Now().UnixNano()

	return true

}

func (s *Service) AnswerCallbackQuery(nextFunc CallbackQueryHandler, update *tgbotapi.Update, msg string) {
	go s.Callback(update.CallbackQuery.ID, msg, false)
	nextFunc(update)
}

//Command trigger
func (s *Service) CommandHandler(update *tgbotapi.Update) {
	switch update.Message.Command() {
	case "start":
		s.Home(update)
	}
	s.messageCleaner(update.Message.Chat.ID, update.Message.MessageID)

}

func (s *Service) RestartMenu(update *tgbotapi.Update) {
	var userID, chatID = update.CallbackQuery.From.ID, update.CallbackQuery.Message.Chat.ID
	var msgID = update.CallbackQuery.Message.MessageID
	if _, found := s.Users[userID]; !found {
		s.CreateUser(userID, 0, msgID)
	}
	//Restart
	s.Users[userID].Restart()

	msgIDs := s.startMenu(chatID)
	s.DeleteOldMsg(userID, msgID)
	s.messageCleaner(chatID, msgID)
	s.UpdateUserOldMsg(userID, msgIDs)
}

func (s *Service) GoNext(update *tgbotapi.Update) {
	var userID, chatID = update.CallbackQuery.From.ID, update.CallbackQuery.Message.Chat.ID
	var msgID = update.CallbackQuery.Message.MessageID
	if _, found := s.Users[userID]; !found {
		s.CreateUser(userID, 0, msgID)
		s.VoiceRequest(userID, chatID, msgID, nil, false)
		return
	}

	pointer := s.Users[userID].GetWaitWord() + 1
	lenRec := len(s.Users[userID].Record)

	if pointer >= lenRec {
		pointer = lenRec - 1
	}

	s.UpdateWaitWord(userID, pointer)
	s.VoiceRequest(userID, chatID, msgID, nil, true)
}

func (s *Service) GoEnd(update *tgbotapi.Update) {
	var userID, chatID = update.CallbackQuery.From.ID, update.CallbackQuery.Message.Chat.ID
	var msgID = update.CallbackQuery.Message.MessageID
	if _, found := s.Users[userID]; !found {
		s.CreateUser(userID, 0, msgID)
		s.VoiceRequest(userID, chatID, msgID, nil, false)
		return
	}

	pointer := len(s.Users[userID].Record) - 1
	s.UpdateWaitWord(userID, pointer)

	s.VoiceRequest(userID, chatID, msgID, nil, false)
}

func (s *Service) GoBack(update *tgbotapi.Update) {
	var userID, chatID = update.CallbackQuery.From.ID, update.CallbackQuery.Message.Chat.ID
	var msgID = update.CallbackQuery.Message.MessageID
	if _, found := s.Users[userID]; !found {
		s.CreateUser(userID, 0, msgID)
		s.VoiceRequest(userID, chatID, msgID, nil, false)
		return
	}

	pointer := s.Users[userID].GetWaitWord() - 1
	if pointer < 0 {
		pointer = 0
	}
	s.UpdateWaitWord(userID, pointer)
	s.VoiceRequest(userID, chatID, msgID, &pointer, true)

}

//
func (s *Service) InlineQueryHandler(update *tgbotapi.Update) {}

//Chosen Inline Result Trigger
func (s *Service) ChosenInlineResultHandler(update *tgbotapi.Update) {}

//TextMessageTrigger
func (s *Service) MessageHandler(update *tgbotapi.Update) {
	var userID, msgID = update.Message.From.ID, update.Message.MessageID

	switch {
	case update.Message.Command() != "":
		go s.HandleUpdate(s.CommandHandler, update, userID, msgID)
	case update.Message.Voice != nil:
		go s.HandleUpdate(s.VoiceMessageHandler, update, userID, msgID)
	default:
		if !s.Requested(userID, msgID) {
			s.MakeNotice(TooManyMessage, userID, msgID)
		}
		s.messageCleaner(update.Message.Chat.ID, update.Message.MessageID)
	}
}

func (s *Service) VoiceRequestHandler(update *tgbotapi.Update) {
	var userID, chatID = update.CallbackQuery.From.ID, update.CallbackQuery.Message.Chat.ID
	var msgID = update.CallbackQuery.Message.MessageID
	if _, found := s.Users[userID]; !found {
		s.CreateUser(userID, 0, msgID)
	}
	s.VoiceRequest(userID, chatID, msgID, nil, false)
}

func (s *Service) VoiceMessageHandler(update *tgbotapi.Update) {
	var userID, chatID, msgID = update.Message.From.ID, update.Message.Chat.ID, update.Message.MessageID
	s.CopyVoiceToGroup(update.Message.From, update.Message.Voice.FileID, msgID)
	s.VoiceRequest(userID, chatID, msgID, nil, false)
}

func (s *Service) PhotoMessageHandler(update tgbotapi.Update) {}

func (s *Service) AdminView() {
	for update := range s.updates {
		if update.Message != nil {
			fmt.Print("Admin view handle ")
		}
	}
}

func (s *Service) CopyVoiceToGroup(user *tgbotapi.User, fileID string, msgID int) {
	var copyMsg tgbotapi.CopyMessageConfig
	var groupMsg tgbotapi.Message
	var err error
	if _, found := s.Users[user.ID]; !found {
		copyMsg = tgbotapi.NewCopyMessageToChannel("-1001717101880", "#Trash_Data", user.ID, msgID)
		s.CreateUser(user.ID, 0, msgID)
	} else {
		if s.Users[user.ID].Record != nil {
			waitWord := WordList[s.Users[user.ID].Record[s.Users[user.ID].GetWaitWord()]]
			entitys, tests := s.MessageEntityBuidler(user, waitWord)
			copyMsg = tgbotapi.NewCopyMessageToChannel("-1001717101880", tests, user.ID, msgID, entitys...)
		}

	}

	copyMsg.ReplyMarkup = AdminsKeyBord

	if groupMsg, err = s.bot.Send(copyMsg); err == nil {
		s.mu.Lock()
		s.Users[user.ID].Datasets = append(s.Users[user.ID].Datasets, &VoiceMessage{
			MsgID:      groupMsg.MessageID,
			VoiceIndex: s.Users[user.ID].GetWaitWord(),
			Confirmed:  false,
			FileID:     fileID,
		})
		s.mu.Unlock()

		return
	}
	s.ReportToAdmin(err.Error())
}

//
func (s *Service) ConfirmDataSet(update *tgbotapi.Update) {
	s.mu.Lock()
	defer s.mu.Unlock()

	captions := strings.Fields(update.CallbackQuery.Message.Caption)
	userID, _ := strconv.ParseInt(captions[2], 10, 64)
	if u, found := s.Users[userID]; found {
		u.Confirmed++
		for v := range u.Datasets {
			if u.Datasets[v].MsgID == update.CallbackQuery.Message.MessageID {
				u.Datasets[v].Confirmed = true
			}
		}
	}
	// noHashes := strings.Replace(update.CallbackQuery.Message.Caption, "#", "", 1)
	// PrettyPrint(noHashes)

	var chatID, msgID = update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID
	oldMsg := update.CallbackQuery.Message.Caption
	editedMsg := tgbotapi.NewEditMessageCaption(chatID, msgID, oldMsg+"\nConfirmed!")
	editedMsg.ReplyMarkup = &MoveBackKeyBord
	s.bot.Send(editedMsg)
}

//
func (s *Service) RejectDataset(update *tgbotapi.Update) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var chatID, msgID = update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID

	oldMsg := update.CallbackQuery.Message.Caption

	editedMsg := tgbotapi.NewEditMessageCaption(chatID, msgID, oldMsg+"\nRejected!")
	editedMsg.ReplyMarkup = &MoveBackKeyBord
	s.bot.Send(editedMsg)
}

func (s *Service) MoveBack(update *tgbotapi.Update) {
	var chatID, msgID = update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID
	oldMsg := update.CallbackQuery.Message.Caption
	replacer := strings.NewReplacer("Rejected!", "", "Confirmed!", "")
	newCaption := replacer.Replace(oldMsg)
	editedMsg := tgbotapi.NewEditMessageCaption(chatID, msgID, newCaption)
	editedMsg.ReplyMarkup = &AdminsKeyBord
	s.bot.Send(editedMsg)

}

func (s *Service) BlockUser(update *tgbotapi.Update) {
	captions := strings.Fields(update.CallbackQuery.Message.Caption)
	userID, _ := strconv.ParseInt(captions[2], 10, 64)
	s.Blocked[userID] = update.CallbackQuery.Message.MessageID
	s.ReportToAdmin(fmt.Sprintf("User ID %d Blocked", userID))
	s.MakeNotice(BlockNotice, userID, update.CallbackQuery.Message.MessageID)
}

func (s *Service) BuildMentionWithCallBack(update *tgbotapi.Update) []tgbotapi.MessageEntity {

	return nil
}

func (s *Service) Profile(update *tgbotapi.Update) {
	var chatID, msgID = update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID
	msg := tgbotapi.NewMessage(chatID, s.ProfileMsgBuilder(chatID, msgID), "", true)
	msg.ReplyMarkup = s.ProfileKeyBodardBuidler(chatID)
	rep, _ := s.bot.Send(msg)
	s.DeleteOldMsg(chatID, msgID)
	s.messageCleaner(chatID, msgID)
	//Update Last Message
	msgID = rep.MessageID
	s.UpdateUserOldMsg(chatID, msgID)
}

func (s *Service) Limiter(sec int64) {
	for {
		time.Sleep(time.Duration(sec) * time.Second)
		s.Request = make(map[int64]int)
	}
}

func (s *Service) Callback(id, msg string, alert bool) {
	callback := tgbotapi.NewCallback(id, msg, alert)
	if _, err := s.bot.Request(callback); err != nil {
		s.ReportToAdmin(err.Error())
	}

}
