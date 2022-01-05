package main

import (
	"fmt"
	"math/rand"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (s *Service) RandWord(userID int64, msgID int) (string, int) {

	if len(s.Users[userID].Record) == len(WordList) {
		return "", 0
	}

	var randValue int = 5
	newSource := rand.NewSource(time.Now().UnixNano())
	newRand := rand.New(newSource)

	for {
		randValue = newRand.Intn(len(WordList) - 1)
		if !s.Users[userID].InUserRecord(randValue) {
			break
		}
	}

	s.UpdateUserRec(userID, randValue)
	s.UpdateWaitWord(userID, len(s.Users[userID].Record)-1)

	return WordList[randValue], randValue
}

func (s *Service) VoiceRequest(userID, chatID int64, msgID int, withPrevious *int, edited bool) {
	if !edited {
		s.VoiceRequestWithNewMessage(userID, msgID, withPrevious)
		return
	}
	s.VoiceRequestWithEditMessage(userID, chatID, msgID, withPrevious)
}

func (s *Service) VoiceRequestWithEditMessage(userID, chatID int64, msgID int, withPrevious *int) {

	var msg tgbotapi.EditMessageTextConfig

	lenRecord := len(s.Users[userID].Record)
	recordPointer := s.Users[userID].RecordPointer

	switch {
	case lenRecord > 0 && recordPointer < lenRecord-1 && withPrevious == nil:
		nextWord := s.Users[userID].RecordPointer + 1
		waitingWord := s.Users[userID].Record[nextWord]
		text := fmt.Sprintf(VoiceRequestMessage, waitingWord, WordList[waitingWord], s.PageViewBuilder(userID))
		msg = tgbotapi.NewEditMessageTextAndMarkup(chatID, msgID, text, "MarkdownV2", false)

		s.UpdateWaitWord(userID, nextWord)
	case withPrevious == nil:
		word, randValue := s.RandWord(userID, msgID)
		word = fmt.Sprintf(VoiceRequestMessage, randValue, word, s.PageViewBuilder(userID))
		msg = tgbotapi.NewEditMessageTextAndMarkup(chatID, msgID, word, "MarkdownV2", false)

	case withPrevious != nil:
		if lenRecord == 0 {
			word, randValue := s.RandWord(userID, msgID)
			word = fmt.Sprintf(VoiceRequestMessage, randValue, word, s.PageViewBuilder(userID))
			msg = tgbotapi.NewEditMessageTextAndMarkup(chatID, msgID, word, "MarkdownV2", false)
		} else {
			waitingWord := s.Users[userID].Record[*withPrevious]
			word := WordList[waitingWord]
			word = fmt.Sprintf(VoiceRequestMessage, waitingWord, word, s.PageViewBuilder(userID))
			msg = tgbotapi.NewEditMessageTextAndMarkup(chatID, msgID, word, "MarkdownV2", false)
		}
	}
	Keyboard := s.UserMenu(userID)
	msg.ReplyMarkup = &Keyboard

	if _, err := s.bot.Send(msg); err != nil {
		s.ReportToAdmin(err.Error())
	}

}

func (s *Service) VoiceRequestWithNewMessage(userID int64, msgID int, withPrevious *int) {

	var msg tgbotapi.MessageConfig
	var rep tgbotapi.Message
	var err error

	lenRecord := len(s.Users[userID].Record)
	recordPointer := s.Users[userID].RecordPointer

	switch {
	case lenRecord > 0 && recordPointer < lenRecord-1 && withPrevious == nil:
		nextWord := recordPointer + 1
		waitingWord := s.Users[userID].Record[nextWord]
		word := fmt.Sprintf(VoiceRequestMessage, waitingWord, WordList[waitingWord], s.PageViewBuilder(userID))
		msg = tgbotapi.NewMessage(userID, word, "MarkdownV2", false)
		s.UpdateWaitWord(userID, nextWord)
	case withPrevious == nil:
		word, randValue := s.RandWord(userID, msgID)
		word = fmt.Sprintf(VoiceRequestMessage, randValue, word, s.PageViewBuilder(userID))
		msg = tgbotapi.NewMessage(userID, word, "MarkdownV2", false)
	case withPrevious != nil:
		if lenRecord == 0 {
			word, randValue := s.RandWord(userID, msgID)
			word = fmt.Sprintf(VoiceRequestMessage, randValue, word, s.PageViewBuilder(userID))
			msg = tgbotapi.NewMessage(userID, word, "MarkdownV2", false)

		} else {
			waitingWord := s.Users[userID].Record[*withPrevious]
			word := fmt.Sprintf(VoiceRequestMessage, waitingWord, WordList[waitingWord], s.PageViewBuilder(userID))
			msg = tgbotapi.NewMessage(userID, word, "MarkdownV2", false)
		}
	}

	msg.ReplyMarkup = s.UserMenu(userID)

	if rep, err = s.bot.Send(msg); err != nil {
		s.ReportToAdmin(err.Error())
		return
	}

	s.DeleteOldMsg(userID, msgID)
	s.messageCleaner(userID, msgID)
	//Update Last Message
	msgID = rep.MessageID
	s.UpdateUserOldMsg(userID, msgID)
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
	s.DeleteOldMsg(userID, msgID)
	s.messageCleaner(chatID, msgID)

	msgID = rep.MessageID
	s.UpdateUserOldMsg(userID, msgID)

}

func (u *User) InUserRecord(randIndex int) bool {
	for _, b := range u.Record {
		if b == randIndex {
			return true
		}
	}
	return false
}

func (s *Service) UpdateUserRec(userID int64, index int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Users[userID].UpdateRec(index)

}
func (s *Service) UpdateWaitWord(userID int64, index int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Users[userID].UpdateWaitWord(index)

}

func (s *Service) GetUserWaitWord(userID, chatID int64, msgID int) int {
	return s.Users[userID].GetWaitWord()
}

func (u *User) UpdateRec(index int) {
	u.Record = append(u.Record, index)
}

func (u *User) UpdateWaitWord(index int) {
	if index < 0 || index > len(u.Record)-1 {
		index = 0
		return
	}

	u.RecordPointer = index
}

func (u *User) GetWaitWord() int {
	return u.RecordPointer
}

func (s *Service) SerchWordIndex(word string) int {
	for i, w := range WordList {
		if word == w {
			return i
		}
	}
	return 0
}
