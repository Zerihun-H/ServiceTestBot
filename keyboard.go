package main

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var MainKeyBord = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("ጀመር", "1"),
	),
)

var MenuButton = tgbotapi.NewInlineKeyboardRow(
	tgbotapi.NewInlineKeyboardButtonData("⨴ በቃኝ", "0"),
	tgbotapi.NewInlineKeyboardButtonData("ሌላ ↺", "-1"),
)

var pointAtFirst = tgbotapi.NewInlineKeyboardRow(
	tgbotapi.NewInlineKeyboardButtonData("▓█▓", "-"),
	tgbotapi.NewInlineKeyboardButtonData("ቀጣይ ⫸", "2"),
)

var pointAtEnd = tgbotapi.NewInlineKeyboardRow(
	tgbotapi.NewInlineKeyboardButtonData("⫷ ተመለስ", "-2"),
	tgbotapi.NewInlineKeyboardButtonData("▓█▓", "-"),
)

var pointAtMiddle = tgbotapi.NewInlineKeyboardRow(
	tgbotapi.NewInlineKeyboardButtonData("⫷ ተመለስ", "-2"),
	tgbotapi.NewInlineKeyboardButtonData("ቀጣይ ⫸", "2"),
)

var EndeKeyBord = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("በዲስ", "3"),
		tgbotapi.NewInlineKeyboardButtonData("በቀድሞ", "-3"),
	), tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("𐙞 እኔ", "4"),
	),
)

var AdminsKeyBord = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("አሳልፍ", "5"),
		tgbotapi.NewInlineKeyboardButtonData("አግድ", "-5"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("ላኪውን አግድ", "6"),
	),
)

//
var MoveBackKeyBord = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("መልስ ➥", "7"),
	))

var AdminKeyBoard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Self Update", "8"),
		tgbotapi.NewInlineKeyboardButtonData("Get Cache ", "9"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Update Words", "10"),
		tgbotapi.NewInlineKeyboardButtonData("Update Caches", "11"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Analytics ", "12"),
	),
)

func (s *Service) UserMenu(userID int64) tgbotapi.InlineKeyboardMarkup {

	lenRecord := len(s.Users[userID].Record)
	recordPointer := s.Users[userID].RecordPointer

	switch {
	case recordPointer == 0 && lenRecord > 1:
		return s.InlineKeyboardMarkupBuilder(MenuButton, pointAtFirst)
	case recordPointer > 0 && recordPointer == lenRecord-1:
		return s.InlineKeyboardMarkupBuilder(MenuButton, pointAtEnd)
	case recordPointer > 0 && recordPointer < lenRecord-1:
		return s.InlineKeyboardMarkupBuilder(MenuButton, pointAtMiddle)
	}

	return s.InlineKeyboardMarkupBuilder(MenuButton)
}

func (s *Service) InlineKeyboardMarkupBuilder(rows ...[]tgbotapi.InlineKeyboardButton) tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

func (s *Service) ProfileKeyBodardBuidler(userID int64) tgbotapi.InlineKeyboardMarkup {
	// botName := @Lambas_bot
	url := fmt.Sprintf("http://t.me/%s?start=%d", s.bot.Self.UserName, userID)
	referral := fmt.Sprintf("https://telegram.me/share/url?url=%s&text=%s", url, "ኑ ላምባን እናስትምር")

	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL("ጓደኛ ጋብዝ", referral),
		), tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("ተመለስ", "0"),
		),
	)
}
