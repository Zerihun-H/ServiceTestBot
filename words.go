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

func (s *Service) RandWord(userID, chatID int64, msgID int) string {
	if _, found := s.Users[userID]; !found {
		s.CreateUser(userID, chatID, msgID)
		return WordList[0]
	}
	if len(s.Users[userID].Record) == len(WordList) {
		return ""
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

	s.UpdateWaitWord(userID, randValue)
	return WordList[randValue]
}

func (s *Service) VoiceRequest(update *tgbotapi.Update) {
	var userID, chatID int64
	var msgID int

	switch {
	case update.Message != nil:
		userID, chatID, msgID = update.Message.From.ID, update.Message.Chat.ID, update.Message.MessageID
	case update.CallbackQuery != nil:
		userID, chatID, msgID = update.CallbackQuery.Message.From.ID, update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID
	}

	msg := tgbotapi.NewMessage(chatID, fmt.Sprintf(VoiceRequestMessage, s.RandWord(userID, chatID, msgID)), "MarkdownV2", false)
	msg.ReplyMarkup = VoiceMarkupButton
	rep, _ := s.bot.Send(msg)

	s.DeleteOldMsg(userID, chatID, msgID)
	s.messageCleaner(chatID, msgID)
	//Update Last Message
	msgID = rep.MessageID
	s.UpdateUserCache(userID, chatID, msgID)
}

func (s *Service) CloseVoiceRequest(update *tgbotapi.Update) {
	var userID, chatID = update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.Chat.ID
	var msgID = update.CallbackQuery.Message.MessageID

	msg := tgbotapi.NewMessage(chatID, ThanksMessage, "", true)
	msg.ReplyMarkup = EndeKeyBord
	rep, _ := s.bot.Send(msg)
	s.DeleteOldMsg(userID, chatID, msgID)
	msgID = rep.MessageID
	s.UpdateUserCache(userID, chatID, msgID)
}

func (s *Service) BackToMenu(update *tgbotapi.Update) {
	var userID, chatID = update.CallbackQuery.Message.From.ID, update.CallbackQuery.Message.Chat.ID
	var msgID = update.CallbackQuery.Message.MessageID

	msgIDs := s.startMenu(chatID)
	s.DeleteOldMsg(userID, chatID, msgID)
	s.UpdateUserCache(userID, chatID, msgIDs)
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

func (s *Service) StopWaitingUsers(userID int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Users[userID].StopWaiting()
}

func (s *Service) GetUserWaitWord(userID, chatID int64, msgID int) int {
	if _, found := s.Users[userID]; !found {
		s.CreateUser(userID, chatID, msgID)
		return 0
	}
	return s.Users[userID].GetWaitWord()
}

func (u *User) StopWaiting() {
	u.WaitingWords = nil
}

func (u *User) UpdateRec(index int) {
	u.Record = append(u.Record, index)
}

func (u *User) UpdateWaitWord(index int) {
	u.WaitingWords = &index
}

func (u *User) GetWaitWord() int {

	if u.WaitingWords == nil {
		return 0
	}
	return *u.WaitingWords
}

func (s *Service) SerchWordIndex(word string) int {
	for i, w := range WordList {
		if word == w {
			return i
		}
	}
	return 0
}
