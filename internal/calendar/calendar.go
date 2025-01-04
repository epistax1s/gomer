package calendar

import (
	"fmt"
	"strconv"
	"time"

	"github.com/epistax1s/gomer/internal/database"

	callback "github.com/epistax1s/gomer/internal/callback"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	BTN_PREV        = "<"
	BTN_NEXT        = ">"
	BTN_PREV_PREFIX = "prev_"
	BTN_NEXT_PREFIX = "next_"
)

var daysOfWeek = [7]string{"Пн", "Вт", "Ср", "Чт", "Пт", "Сб", "Вс"}

func GenerateCalendar(year int, month time.Month) *tgbotapi.InlineKeyboardMarkup {
	var keyboard tgbotapi.InlineKeyboardMarkup

	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, generateMonthYearRow(year, month))
	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, generateDaysNamesRow())
	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, generateDaysInMonth(year, month)...)
	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, generateNavigationButtons(year, month))

	return &keyboard
}

func HandlePrevButton(year int, month time.Month) (*tgbotapi.InlineKeyboardMarkup, int, time.Month) {
	if month == time.January {
		year--
		month = time.December
	} else {
		month--
	}
	return GenerateCalendar(year, month), year, month
}

func HandleNextButton(year int, month time.Month) (*tgbotapi.InlineKeyboardMarkup, int, time.Month) {
	if month == time.December {
		year++
		month = time.January
	} else {
		month++
	}
	return GenerateCalendar(year, month), year, month
}

func HandleButtonData(callbackData string) (*database.Date, error) {
	layout := "2006-01-02"
	parsedTime, err := time.Parse(layout, callbackData)
	if err != nil {
		return nil, err
	}

	date := database.Date{Time: parsedTime}
	return &date, nil
}

func generateMonthYearRow(year int, month time.Month) []tgbotapi.InlineKeyboardButton {
	return tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("%s %v", month, year), "current_month"),
	)
}

func generateDaysNamesRow() []tgbotapi.InlineKeyboardButton {
	var row []tgbotapi.InlineKeyboardButton
	for _, day := range daysOfWeek {
		row = append(row, tgbotapi.NewInlineKeyboardButtonData(day, day))
	}
	return row
}

func generateDaysInMonth(year int, month time.Month) [][]tgbotapi.InlineKeyboardButton {
	firstDay := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	daysInMonth := firstDay.AddDate(0, 1, -1).Day()

	var rows [][]tgbotapi.InlineKeyboardButton
	var row []tgbotapi.InlineKeyboardButton

	weekdayOffset := int(firstDay.Weekday())
	if weekdayOffset == 0 {
		weekdayOffset = 7
	}
	for i := 1; i < weekdayOffset; i++ {
		row = append(row, tgbotapi.NewInlineKeyboardButtonData(" ", "empty"))
	}

	for day := 1; day <= daysInMonth; day++ {
		btnText := strconv.Itoa(day)
		if time.Now().Year() == year && time.Now().Month() == month && time.Now().Day() == day {
			btnText += "!"
		}
		callbackData := callback.NewCalendarDateCallback(fmt.Sprintf("%d-%02d-%02d", year, month, day))
		row = append(row, tgbotapi.NewInlineKeyboardButtonData(btnText, callbackData))

		if len(row) == 7 {
			rows = append(rows, row)
			row = nil
		}
	}

	if len(row) > 0 {
		for len(row) < 7 {
			row = append(row, tgbotapi.NewInlineKeyboardButtonData(" ", "empty"))
		}
		rows = append(rows, row)
	}

	return rows
}

func generateNavigationButtons(year int, month time.Month) []tgbotapi.InlineKeyboardButton {
	return tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(BTN_PREV, callback.NewCalendarPrevCallback(year, month)),
		tgbotapi.NewInlineKeyboardButtonData(BTN_NEXT, callback.NewCalendarNextCallback(year, month)),
	)
}
