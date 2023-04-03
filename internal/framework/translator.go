package framework

import (
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var matcher = language.NewMatcher(
	[]language.Tag{
		language.English,
	},
)

func NewTranslator(locale string) *message.Printer {
	tag, _ := language.MatchStrings(matcher, locale)

	return message.NewPrinter(tag)
}
