package translation

import (
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var matcher = language.NewMatcher(
	[]language.Tag{
		language.MustParse("en-GB"),
		language.MustParse("ru-RU"),
		language.MustParse("uk-UA"),
	},
)

func NewTranslator(locale string) *message.Printer {
	tag := GetSupportedLocale(locale)

	return message.NewPrinter(tag)
}

func GetSupportedLocale(locale string) language.Tag {
	tag, _ := language.MatchStrings(matcher, locale)

	return tag
}
