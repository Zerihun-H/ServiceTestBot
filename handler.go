package main

import (
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

//Hnadler Trigger Callback from msg or inline
func (s *Service) CallbackQueryHandler(update *tgbotapi.Update) {
	switch update.CallbackQuery.Data {
	case "1", "-1":
		go s.VoiceRequestHandler(update)
	case "0":
		go s.CloseVoiceRequest(update)
	case "3":
		go s.RestartMenu(update)
	case "-3":
		go s.VoiceRequestHandler(update)
	case "4":
		go s.ConfirmDataSet(update)
	case "-4":
		go s.RejectDataset(update)
	case "6":
		go s.BlockUser(update)
	case "7":
		go s.MoveBack(update)
	case "-2":
		go s.GoBack(update)
	case "2":
		go s.GoNext(update)
	}

}

func (s *Service) GoNext(update *tgbotapi.Update) {
	var userID, chatID = update.CallbackQuery.From.ID, update.CallbackQuery.Message.Chat.ID
	var msgID = update.CallbackQuery.Message.MessageID
	if _, found := s.Users[userID]; !found {
		s.CreateUser(userID, chatID, msgID)
		s.VoiceRequest(userID, chatID, msgID, nil, false)
		return
	}
	pointer := s.Users[userID].GetWaitWord() + 1
	s.UpdateWaitWord(userID, pointer)
	s.VoiceRequest(userID, chatID, msgID, &pointer, true)
	PrettyPrint(s.Users[userID])
}

func (s *Service) GoBack(update *tgbotapi.Update) {
	var userID, chatID = update.CallbackQuery.From.ID, update.CallbackQuery.Message.Chat.ID
	var msgID = update.CallbackQuery.Message.MessageID
	if _, found := s.Users[userID]; !found {
		s.CreateUser(userID, chatID, msgID)
		s.VoiceRequest(userID, chatID, msgID, nil, false)
		return
	}

	pointer := s.Users[userID].GetWaitWord() - 1
	s.UpdateWaitWord(userID, pointer)

	s.VoiceRequest(userID, chatID, msgID, &pointer, true)

	PrettyPrint(s.Users[userID])

}

//
func (s *Service) InlineQueryHandler(update *tgbotapi.Update) {}

//Chosen Inline Result Trigger
func (s *Service) ChosenInlineResultHandler(update *tgbotapi.Update) {}

//TextMessageTrigger
func (s *Service) MessageHandler(update *tgbotapi.Update) {
	switch {
	case update.Message.Command() != "":
		go s.CommandHandler(update)
	case update.Message.Voice != nil:
		go s.VoiceMessageHandler(update)
	default:
		s.messageCleaner(update.Message.Chat.ID, update.Message.MessageID)

	}
	if _, found := s.Users[update.Message.From.ID]; found {
		PrettyPrint(s.Users[update.Message.From.ID])
	}

}

func (s *Service) VoiceRequestHandler(update *tgbotapi.Update) {
	var userID, chatID = update.CallbackQuery.From.ID, update.CallbackQuery.Message.Chat.ID
	var msgID = update.CallbackQuery.Message.MessageID
	if _, found := s.Users[userID]; !found {
		s.CreateUser(userID, chatID, msgID)
	}
	s.VoiceRequest(userID, chatID, msgID, nil, false)
}

func (s *Service) VoiceMessageHandler(update *tgbotapi.Update) {
	userID, chatID, msgID := update.Message.From.ID, update.Message.Chat.ID, update.Message.MessageID

	s.CopyVoiceToGroup(userID, chatID, msgID)
	s.VoiceRequest(userID, chatID, msgID, nil, false)
}

//Command trigger
func (s *Service) CommandHandler(update *tgbotapi.Update) {
	switch update.Message.Command() {
	case "start":
		s.Home(update)
	}
	s.messageCleaner(update.Message.Chat.ID, update.Message.MessageID)

}

func (s *Service) PhotoMessageHandler(update tgbotapi.Update) {}

func (s *Service) AdminView() {
	for update := range s.updates {
		if update.Message != nil {
			fmt.Print("Admin view handle ")
		}
	}
}

func (s *Service) CopyVoiceToGroup(userID, chatID int64, msgID int) {
	var copyMsg tgbotapi.CopyMessageConfig

	if _, found := s.Users[userID]; !found {

		copyMsg = tgbotapi.NewCopyMessageToChannel("-1001717101880", "#Trash_Data", chatID, msgID)
		s.CreateUser(userID, chatID, msgID)
	} else {

		waitWord := WordList[s.Users[userID].Record[s.Users[userID].GetWaitWord()]]
		copyMsg = tgbotapi.NewCopyMessageToChannel("-1001717101880", "#"+waitWord, chatID, msgID)
	}
	copyMsg.ReplyMarkup = AdminsKeyBord
	s.bot.Send(copyMsg)
	// println(s.bot.GetFileDirectURL(message.Voice.FileID))
}

//
func (s *Service) ConfirmDataSet(update *tgbotapi.Update) {

	noHashes := strings.Replace(update.CallbackQuery.Message.Caption, "#", "", 1)
	PrettyPrint(noHashes)

	var chatID, msgID = update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID
	oldMsg := update.CallbackQuery.Message.Caption
	editedMsg := tgbotapi.NewEditMessageCaption(chatID, msgID, oldMsg+"\nConfirmed!")
	editedMsg.ReplyMarkup = &MoveBackKeyBord
	s.bot.Send(editedMsg)
}

//
func (s *Service) RejectDataset(update *tgbotapi.Update) {
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

}
