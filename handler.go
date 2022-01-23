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
		go s.GoEnd(update)
	case "0":
		go s.CloseVoiceRequest(update)
	case "-2":
		go s.AnswerCallbackQuery(s.GoBack, update, "ወደኋላ")
	case "2":
		go s.AnswerCallbackQuery(s.GoNext, update, "ቀጣይ")
	case "3":
		go s.RestartMenu(update)
	case "-3":
		go s.GoEnd(update)
	case "4":
		go s.Profile(update)
	case "5":
		go s.ConfirmDataSet(update, false)
	case "-5":
		go s.RejectDataset(update, false)
	case "6":
		go s.BlockUser(update)
	case "7":
		go s.MoveBack(update, false)
	case "14":
		go s.ConfirmDataSet(update, true)
	case "-14":
		go s.RejectDataset(update, true)
	case "00":
		go s.BlockUser(update)
	}

}

func (s *Service) CloseVoiceRequest(update *tgbotapi.Update) {
	var userID, chatID = update.CallbackQuery.From.ID, update.CallbackQuery.Message.Chat.ID
	var msgID = update.CallbackQuery.Message.MessageID

	if _, found := s.Users[userID]; !found {
		s.CreateUser(userID, 0, msgID)
	}

	msg := tgbotapi.NewMessage(chatID, ThanksMessage, "", true)
	msg.ReplyMarkup = EndeKeyBord
	rep, _ := s.bot.Send(msg)
	go s.Users[userID].UpdateWaitWord(0)

	s.messageCleaner(userID, msgID)

	msgID = rep.MessageID
	s.UpdateUserOldMsg(userID, msgID)
	s.UpdateUserMenuState(userID, UserMenuPage)
}

func (s *Service) Requested(userID int64, msgID int) bool {
	var user *User
	var found bool

	if user, found = s.Users[userID]; !found {
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
	case "enabled":
		s.VerificationEnabled(update.Message.From.ID)
	case "disabled":
		s.verificationDisabled(update.Message.From.ID)
	}
	s.messageCleaner(update.Message.Chat.ID, update.Message.MessageID)

}
func (s *Service) VerificationEnabled(userID int64) {
	s.mu.Lock()
	if u, found := s.Users[userID]; found {
		u.VerifiedSetting = true
		println("Enabled")
	}
	s.mu.Unlock()

}
func (s *Service) verificationDisabled(userID int64) {
	s.mu.Lock()

	if u, found := s.Users[userID]; found {
		u.VerifiedSetting = false
		println("Disabled")
	}

	s.mu.Unlock()
}

func (s *Service) RestartMenu(update *tgbotapi.Update) {
	var userID, msgID = update.CallbackQuery.From.ID, update.CallbackQuery.Message.MessageID
	if _, found := s.Users[userID]; !found {
		s.CreateUser(userID, 0, msgID)
	}
	//Restart
	s.Users[userID].Restart()

	msgIDs := s.startMenu(userID)

	s.messageCleaner(userID, msgID)

	s.UpdateUserOldMsg(userID, msgIDs)
}

func (s *Service) GoNext(update *tgbotapi.Update) {
	var userID, chatID = update.CallbackQuery.From.ID, update.CallbackQuery.Message.Chat.ID
	var msgID = update.CallbackQuery.Message.MessageID
	if _, found := s.Users[userID]; !found {
		s.CreateUser(userID, 0, msgID)
		s.VoiceRequest(userID, chatID, msgID, nil, false, false)
		return
	}

	pointer := s.Users[userID].GetWaitWord() + 1
	lenRec := len(s.Users[userID].Record)

	if pointer >= lenRec {
		pointer = lenRec - 1
	}

	s.UpdateWaitWord(userID, pointer)
	s.VoiceRequest(userID, chatID, msgID, nil, true, false)
}

func (s *Service) GoEnd(update *tgbotapi.Update) {
	var userID, chatID = update.CallbackQuery.From.ID, update.CallbackQuery.Message.Chat.ID
	var msgID = update.CallbackQuery.Message.MessageID
	if _, found := s.Users[userID]; !found {
		s.CreateUser(userID, 0, msgID)
		s.VoiceRequest(userID, chatID, msgID, nil, false, false)
		return
	}

	pointer := len(s.Users[userID].Record) - 1
	s.UpdateWaitWord(userID, pointer)

	s.VoiceRequest(userID, chatID, msgID, nil, false, false)
}

func (s *Service) GoBack(update *tgbotapi.Update) {
	var userID, chatID = update.CallbackQuery.From.ID, update.CallbackQuery.Message.Chat.ID
	var msgID = update.CallbackQuery.Message.MessageID
	if _, found := s.Users[userID]; !found {
		s.CreateUser(userID, 0, msgID)
		s.VoiceRequest(userID, chatID, msgID, nil, false, false)
		return
	}

	pointer := s.Users[userID].GetWaitWord() - 1
	if pointer < 0 {
		pointer = 0
	}
	s.UpdateWaitWord(userID, pointer)
	s.VoiceRequest(userID, chatID, msgID, &pointer, true, false)

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
	s.VoiceRequest(userID, chatID, msgID, nil, false, false)
}

func (s *Service) VoiceMessageHandler(update *tgbotapi.Update) {
	var userID, chatID, msgID = update.Message.From.ID, update.Message.Chat.ID, update.Message.MessageID

	if !s.CopyVoiceToGroup(update.Message.From, update.Message.Voice.FileID, msgID) {
		s.VoiceRequest(userID, chatID, msgID, nil, false, true)
		return
	}

	s.messageCleaner(userID, msgID)
	s.DeleteOldMsg(userID)
}

func (s *Service) PhotoMessageHandler(update tgbotapi.Update) {}

func (s *Service) AdminView() {
	// for update := range s.updates {
	// 	if update.Message != nil {

	// 	}
	// }
}

func (s *Service) CopyVoiceToGroup(user *tgbotapi.User, fileID string, msgID int) bool {
	var copyMsg tgbotapi.CopyMessageConfig
	var groupMsg tgbotapi.Message
	var err error
	u, found := s.Users[user.ID]
	switch {
	case !found:
		copyMsg = tgbotapi.NewCopyMessageToChannel(DatasetGroup, "#Trash_Data", user.ID, msgID)
		s.CreateUser(user.ID, 0, msgID)
		return false
	case u.MenuState != WaitingVoice:
		copyMsg = tgbotapi.NewCopyMessageToChannel(DatasetGroup, "#Trash_Data", user.ID, msgID)
		copyMsg.ReplyMarkup = BlockBtn
		s.bot.Send(copyMsg)
		return false
	default:

		if u.Record != nil {
			waitWord := WordList[u.Record[u.GetWaitWord()]]
			entitys, tests := s.MessageEntityBuidler(user, waitWord)
			copyMsg = tgbotapi.NewCopyMessageToChannel(DatasetGroup, tests, user.ID, msgID, entitys...)
			copyMsg.ReplyMarkup = AdminKeyBoard
			go s.UpdateTreeRecord(user.ID)
			if groupMsg, err = s.bot.Send(copyMsg); err == nil {

				s.mu.Lock()
				u.Datasets = append(u.Datasets, &VoiceMessage{
					MsgID:      groupMsg.MessageID,
					VoiceIndex: u.GetWaitWord(),
					Confirmed:  false,
					FileID:     fileID,
				})
				s.mu.Unlock()

				if u.VerifiedSetting {
					copyMsgUser := tgbotapi.NewCopyMessage(user.ID, "#"+waitWord+"\n"+fmt.Sprint(groupMsg.MessageID), user.ID, msgID)
					copyMsgUser.ReplyMarkup = UserVoiceMenu
					if _, err := s.bot.Send(copyMsgUser); err != nil {
						s.ReportToAdmin(err.Error())
						return false
					}
					u.MenuState = WaitingVerification
					return true
				}
				return false
			}
			s.ReportToAdmin(err.Error())
			return false
		}
		return false
	}
}

// 	copyMsg = tgbotapi.NewCopyMessageToChannel(BackupGroup, tests, user.ID, msgID, entitys...)
// 	copyMsg.ReplyMarkup = AdminBackUpKeyBord
// 	if groupMsg, err = s.bot.Send(copyMsg); err == nil {
// 		u.BackUpDatasets = append(u.BackUpDatasets, &VoiceMessage{
// 			MsgID:      groupMsg.MessageID,
// 			VoiceIndex: u.GetWaitWord(),
// 			Confirmed:  false,
// 			FileID:     fileID,
// 		})

// 	}

// 	return false
// }
//
func (s *Service) ConfirmDataSet(update *tgbotapi.Update, isUserCcmmand bool) {
	var msgBody, noHashesWord, userName, firstName string
	var moveKeyBord *tgbotapi.InlineKeyboardMarkup
	var userID int64
	var msgID int

	switch {
	case isUserCcmmand:
		captions := strings.Fields(update.CallbackQuery.Message.Caption)
		noHashesWord = strings.Replace(captions[0], "#", "", 1)
		firstName = update.CallbackQuery.From.FirstName
		userName = update.CallbackQuery.From.UserName
		userID = update.CallbackQuery.Message.Chat.ID
		msgID, _ = strconv.Atoi(captions[1])
		msgBody = "\nConfirmByUser"
		moveKeyBord = &AdminKeyBoard
		go s.VoiceRequest(userID, userID, update.CallbackQuery.Message.MessageID, nil, false, true)

	default:
		captions := strings.Fields(update.CallbackQuery.Message.Caption)
		noHashesWord = strings.Replace(captions[0], "#", "", 1)
		userID, _ = strconv.ParseInt(captions[2], 10, 64)
		msgID = update.CallbackQuery.Message.MessageID
		firstName = captions[1]
		userName = captions[3]
		msgBody = "\nConfirmed!"
		moveKeyBord = &MoveBackKeyBord
		go s.ConfirmUserDataSet(userID, msgID)

		if len(userName) > 2 {
			userName = strings.Replace(captions[3], "@", "", -1)
		}
	}

	entitys, texts := s.MessageEntityBuidler(&tgbotapi.User{
		ID:        userID,
		IsBot:     false,
		FirstName: firstName,
		UserName:  userName,
	}, noHashesWord)

	// PrettyPrint(noHashes)

	editedMsg := tgbotapi.NewEditMessageCaption(-1001717101880, msgID, texts+msgBody, entitys...)
	editedMsg.ReplyMarkup = moveKeyBord
	s.bot.Send(editedMsg)

}

//
func (s *Service) RejectDataset(update *tgbotapi.Update, isUserCcmmand bool) {
	var msgBody, noHashesWord, userName, firstName string
	var userID int64
	var msgID int

	switch {
	case isUserCcmmand:
		captions := strings.Fields(update.CallbackQuery.Message.Caption)
		noHashesWord = strings.Replace(captions[0], "#", "", 1)
		firstName = update.CallbackQuery.From.FirstName
		userName = update.CallbackQuery.From.UserName
		userID = update.CallbackQuery.Message.Chat.ID
		msgID, _ = strconv.Atoi(captions[1])
		msgBody = "\nRejectedByUser"

		go s.VoiceRequest(userID, userID, update.CallbackQuery.Message.MessageID, nil, false, true)

	default:
		captions := strings.Fields(update.CallbackQuery.Message.Caption)
		noHashesWord = strings.Replace(captions[0], "#", "", 1)
		userID, _ = strconv.ParseInt(captions[2], 10, 64)
		msgID = update.CallbackQuery.Message.MessageID
		firstName = captions[1]
		msgBody = "\nRejected!"
		userName = captions[3]

		go s.RejectUserDataset(userID, msgID)

		if len(userName) > 2 {
			userName = strings.Replace(captions[3], "@", "", -1)
		}
	}

	entitys, texts := s.MessageEntityBuidler(&tgbotapi.User{
		ID:        userID,
		IsBot:     false,
		FirstName: firstName,
		UserName:  userName,
	}, noHashesWord)

	editedMsg := tgbotapi.NewEditMessageCaption(-1001717101880, msgID, texts+msgBody, entitys...)
	editedMsg.ReplyMarkup = &MoveBackKeyBord
	s.bot.Send(editedMsg)
}

func (s *Service) MoveBack(update *tgbotapi.Update, backUp bool) {

	var chatID, msgID = update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID
	// oldCaption := update.CallbackQuery.Message.Caption

	captions := strings.Fields(update.CallbackQuery.Message.Caption)
	userID, _ := strconv.ParseInt(captions[2], 10, 64)
	noHashesWord := strings.Replace(captions[0], "#", "", 1)
	userName := captions[3]

	switch {
	case len(userName) > 2:
		userName = strings.Replace(captions[3], "@", "", -1)
	default:
		userName = ""
	}

	entitys, texts := s.MessageEntityBuidler(&tgbotapi.User{
		ID:        userID,
		IsBot:     false,
		FirstName: captions[1],
		UserName:  userName,
	}, noHashesWord)

	//Whe Manual Work :)

	// replacer := strings.NewReplacer("Rejected!", "", "Confirmed!", "")
	// newCaption := replacer.Replace(oldCaption)

	editedMsg := tgbotapi.NewEditMessageCaption(chatID, msgID, texts, entitys...)
	editedMsg.ReplyMarkup = &AdminKeyBoard
	s.bot.Send(editedMsg)

}

func (s *Service) BlockUser(update *tgbotapi.Update) {
	if !s.IsAdmin(update.CallbackQuery.From.ID) {
		return
	}
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
