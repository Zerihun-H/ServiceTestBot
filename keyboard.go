package main

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	separator    = "⦁"
	continueIcon = "᠁"
	opener       = "⎾"
	closer       = "⏌"
)

type Builder struct {
	Inc      int
	Position bool
	Start    int
	Size     int
}

func NewBuilder(start, Size int, position bool) *Builder {
	return &Builder{
		Inc:      0,
		Position: position,
		Start:    start,
		Size:     Size,
	}
}

func (b *Builder) Execute() {
	switch {
	case b.Position:
		b.Start++
		b.Inc++
	default:
		b.Start--
		b.Inc++
	}
}

var MainKeyBord = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("ጀምር", "1"),
	),
)

var MenuButton = tgbotapi.NewInlineKeyboardRow(
	tgbotapi.NewInlineKeyboardButtonData("⨴ በቃኝ", "0"),
	tgbotapi.NewInlineKeyboardButtonData("ሌላ ⨮", "-1"),
)

var pointAtFirst = tgbotapi.NewInlineKeyboardRow(
	// tgbotapi.NewInlineKeyboardButtonData("ㅤ", "-"),
	tgbotapi.NewInlineKeyboardButtonData("ቀጣይ ⫸", "2"),
)

var pointAtEnd = tgbotapi.NewInlineKeyboardRow(
	tgbotapi.NewInlineKeyboardButtonData("⫷ ወደኋላ", "-2"),
	// tgbotapi.NewInlineKeyboardButtonData("ㅤ", "-"),
)

var pointAtMiddle = tgbotapi.NewInlineKeyboardRow(
	tgbotapi.NewInlineKeyboardButtonData("⫷ ወደኋላ", "-2"),
	tgbotapi.NewInlineKeyboardButtonData("ቀጣይ ⫸", "2"),
)

var EndeKeyBord = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("↻ በአዲስ", "3"),
		tgbotapi.NewInlineKeyboardButtonData("በቀድሞ ⟴", "-3"),
	), tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("𐙞 እኔ", "4"),
	),
)

var AdminsKeyBord = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("⊕ አሳልፍ", "5"),
		tgbotapi.NewInlineKeyboardButtonData("አሳቀር ⊝", "-5"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("⊗ ላኪውን አግድ ⊗", "6"),
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
		tgbotapi.NewInlineKeyboardButtonData("Update Voice Path", "10"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Update Words", "10"),
		tgbotapi.NewInlineKeyboardButtonData("Update Caches", "11"),
	),

	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Block", "10"),
		tgbotapi.NewInlineKeyboardButtonData("leaderboard", "11"),
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

func (s *Service) PageViewBuilder(userID int64) string {
	lenRecord := len(s.Users[userID].Record)
	recordPointer := s.Users[userID].RecordPointer
	var data, icon, data2, icon2 string
	var backLimit, forwardLimit int

	switch {
	case recordPointer == 0 && lenRecord > 1:
		data, icon = NewBuilder(recordPointer, lenRecord, true).String()
		return fmt.Sprintf("%s1%s%s%s", opener, closer, data, icon)
	case recordPointer > 0 && recordPointer == lenRecord-1:
		data, icon = NewBuilder(recordPointer, lenRecord, false).String()

		return fmt.Sprintf("%s%s%s%d%s", icon, data, opener, recordPointer+1, closer)
	case recordPointer > 0 && recordPointer < lenRecord-1:
		backLimit = lenRecord - 1 - recordPointer
		backLimit = 6 - backLimit
		if backLimit < 3 {
			backLimit = 3
		}
		data, icon = NewBuilder(recordPointer, lenRecord, false).String(backLimit)
		forwardLimit = 6 - recordPointer
		if forwardLimit < 3 {
			forwardLimit = 3
		}
		// limit =
		data2, icon2 = NewBuilder(recordPointer, lenRecord, true).String(forwardLimit)
		return fmt.Sprintf("%s%s%s%d%s%s%s", icon, data, opener, recordPointer+1, closer, data2, icon2)
	}

	return ""
}

func (b *Builder) String(limits ...int) (string, string) {
	var data string
	var limit int

	switch {
	case len(limits) == 0:
		limit = 6
	default:
		limit = limits[0]
	}

	for {

		b.Execute()
		switch {
		case b.Position:
			data = data + separator + fmt.Sprint(b.Start+1)
			if b.Start == b.Size-1 {
				return data, " "
			}
		default:
			data = fmt.Sprint(b.Start+1) + separator + data
			if b.Start == 0 {
				return data, " "
			}
		}

		if b.Inc == limit {
			switch {
			case b.Position:
				return data, continueIcon
			default:
				return data, continueIcon
			}
		}

	}
}

func (s *Service) InlineKeyboardMarkupBuilder(rows ...[]tgbotapi.InlineKeyboardButton) tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

func (s *Service) ProfileKeyBodardBuidler(userID int64) tgbotapi.InlineKeyboardMarkup {
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
