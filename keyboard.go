package main

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type keyboard struct {
}

func (k *keyboard) Maker(func(string)) {

}

func (k *keyboard) Editor() tgbotapi.InlineKeyboardMarkup {

	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📈 Grade Summary ", "1"),
			tgbotapi.NewInlineKeyboardButtonData("🔔 Notification ", "2"),
		))
}

var MainKeyBord = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("ጀመር", "0"),
	),
)

var VoiceMarkupButton = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("በቃኝ", "1"),
		tgbotapi.NewInlineKeyboardButtonData("ሌላ", "2"),
	),
)

var EndeKeyBord = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("በዲስ አስጀምር", "3"),
	),
)

var AdminsKeyBord = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("አሳልፍ", "4"),
		tgbotapi.NewInlineKeyboardButtonData("አግድ", "5"),
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
