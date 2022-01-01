package main

import (
	"fmt"
	"math/rand"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var WordList = []string{
	"ሰላም", "ላምባ", "ከ", "የ", "በ", "ለ", "ኛ", "ደህና", "አደርሽ", "ዋልሽ", "አመሸሽ", "ነሽ",
	"አይ", "ተዪው", "ይቅር", "አልፈልግም", "በቃኝ", "ይበቃኛል", "ቻው", "እሺ", "አዎ", "እፈልጋለሁ",
	"ቀጥይ", "ይደገም", "ድገሚ", "ድገሚልኝ", "ድገሚው", "ነው", "ላኪልኝ", "ላኪ", "ላኪው", "አንብቢልኝ",
	"አዲስ", "ያዢልኝ", "አለ", "አለኝ", "ንገሪኝ", "ሆነ", "ቀኑ", "ላይ", "ተልኮልኛል", "ነበረኝ", "ዘርዝሪሊኝ", "ዘርዝሪ",
	"አንብቢ", "ስንት", "መልዕክት", "ያዢ", "አስቀምጪ", "የዛሬው", "የዛሬ", "ሲደመር", "ሲቀነስ", "ሲካፈል", "ሲባዛ", "ደምሪ",
	"ቀንሺ", "አካፍዪ", "አባዢ", "ሰኞ", "ማክሰኞ", "ዕሮብ", "ሀሙስ", "አርብ", "ቅዳሜ", "እሁድ", "ነገ", "ትናት", "ዛሬ", "ድምጽ", "ቀንሺ",
	"አጥፊ", "ዝቅ", "አርጊ", "ጨምሪልኝ", "ስልክ", "ቁጥሩ", "ቁጥር", "ማንቂያ", "ደውል", "ሙይልኝ", "ሙይ", "ሰአት", "ንገሪኝ", "ቀን", "ቀኑን",
	"ዜሮ", "አንድ", "ሁለት", "ሶስት", "አራት", "አምስት", "ስድስት", "ሰባት", "ስምንት", "ዘጠኝ", "አስር", "አስራ", "ሃያ", "ሰላሳ", "አርባ",
	"ሃምሳ", "ስልሳ", "ሰባ", "ሰማንያ", "ዘጠና", "መቶ", "ሺ", "ሚሊየን", "የት", "ምን", "ማን", "እንዴት", "ለምን", "አይነቶች",
	"ክፍሎች", "ዝርዝሮች", "ዋጋ", "ያህል", "መች", "ትርጉም", "ትርጉሙ", "መቼ", "ማለት", "ምንድን", "ነኝ", "ያለሁበት",
	"ቦታ", "ያለሁት", "ያህል", "ይርቃል", "እዚህ", "ምን", "አዲስ", "ነገር", "ውሎ", "ዜና", "ነበር", "አየሩ", "ጃኬት",
	"ጃንጥላ", "መያዝ", "አለብኝ", "የአየር", "ሁኔታ", "መልበስ", "ደውይ", "ደውይልኝ", "ደውይለት",
	"ደውይላት", "ማስታወሻ", "ቴሌግራም", "ኢሜል", "መጽሐፍ", "ክፍል", "ታሪክ", "ሚስኮል", "ትዕዛዝ"}

func (s *Service) RandWord(userID, chatID int64, msgID int) (string, int) {

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
		s.VoiceRequestWithNewMessage(userID, chatID, msgID, withPrevious)
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
		text := fmt.Sprintf(VoiceRequestMessage, waitingWord, WordList[waitingWord])
		msg = tgbotapi.NewEditMessageTextAndMarkup(chatID, msgID, text, "MarkdownV2", false)

		s.UpdateWaitWord(userID, nextWord)
	case withPrevious == nil:
		word, randValue := s.RandWord(userID, chatID, msgID)
		word = fmt.Sprintf(VoiceRequestMessage, randValue, word)
		msg = tgbotapi.NewEditMessageTextAndMarkup(chatID, msgID, word, "MarkdownV2", false)

	case withPrevious != nil:
		waitingWord := s.Users[userID].Record[*withPrevious]
		word := WordList[waitingWord]
		word = fmt.Sprintf(VoiceRequestMessage, waitingWord, word)
		msg = tgbotapi.NewEditMessageTextAndMarkup(chatID, msgID, word, "MarkdownV2", false)
	}
	Keyboard := s.UserMenu(userID)
	msg.ReplyMarkup = &Keyboard
	s.bot.Send(msg)

}

func (s *Service) VoiceRequestWithNewMessage(userID, chatID int64, msgID int, withPrevious *int) {

	var msg tgbotapi.MessageConfig
	lenRecord := len(s.Users[userID].Record)
	recordPointer := s.Users[userID].RecordPointer

	println("Record Points :", recordPointer)
	println("Record Lens :", lenRecord)

	switch {
	case lenRecord > 0 && recordPointer < lenRecord-1 && withPrevious == nil:
		nextWord := s.Users[userID].RecordPointer + 1
		waitingWord := s.Users[userID].Record[nextWord]
		word := fmt.Sprintf(VoiceRequestMessage, waitingWord, WordList[waitingWord])
		msg = tgbotapi.NewMessage(chatID, word, "MarkdownV2", false)
		s.UpdateWaitWord(userID, nextWord)
	case withPrevious == nil:
		word, randValue := s.RandWord(userID, chatID, msgID)
		word = fmt.Sprintf(VoiceRequestMessage, randValue, word)
		msg = tgbotapi.NewMessage(chatID, word, "MarkdownV2", false)
	case withPrevious != nil:

		waitingWord := s.Users[userID].Record[*withPrevious]
		word := fmt.Sprintf(VoiceRequestMessage, waitingWord, WordList[waitingWord])
		msg = tgbotapi.NewMessage(chatID, word, "MarkdownV2", false)
	}

	msg.ReplyMarkup = s.UserMenu(userID)
	rep, _ := s.bot.Send(msg)

	s.DeleteOldMsg(userID, chatID, msgID)
	s.messageCleaner(chatID, msgID)
	//Update Last Message
	msgID = rep.MessageID
	s.UpdateUserOldMsg(userID, chatID, msgID)
}

func (s *Service) CloseVoiceRequest(update *tgbotapi.Update) {
	var userID, chatID = update.CallbackQuery.From.ID, update.CallbackQuery.Message.Chat.ID
	var msgID = update.CallbackQuery.Message.MessageID

	if _, found := s.Users[userID]; !found {
		s.CreateUser(userID, chatID, msgID)
	}

	msg := tgbotapi.NewMessage(chatID, ThanksMessage, "", true)
	msg.ReplyMarkup = EndeKeyBord
	rep, _ := s.bot.Send(msg)

	s.DeleteOldMsg(userID, chatID, msgID)
	s.messageCleaner(chatID, msgID)

	msgID = rep.MessageID
	s.UpdateUserOldMsg(userID, chatID, msgID)

}

func (s *Service) RestartMenu(update *tgbotapi.Update) {
	var userID, chatID = update.CallbackQuery.From.ID, update.CallbackQuery.Message.Chat.ID
	var msgID = update.CallbackQuery.Message.MessageID
	if _, found := s.Users[userID]; !found {
		s.CreateUser(userID, chatID, msgID)
	}
	//Restart
	s.Users[userID].Restart()

	msgIDs := s.startMenu(chatID)
	s.DeleteOldMsg(userID, chatID, msgID)
	s.messageCleaner(chatID, msgID)
	s.UpdateUserOldMsg(userID, chatID, msgIDs)
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
