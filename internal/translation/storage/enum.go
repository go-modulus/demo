package storage

import (
	"fmt"
	"golang.org/x/text/language"
)

func LocaleFromTag(tag language.Tag) (Locale, error) {
	base, _ := tag.Base()

	switch base.String() {
	case "en":
		return LocaleEn, nil
	case "id":
		return LocaleID, nil
	}
	return Locale(base.String()), fmt.Errorf("unsupported locale: %s", base.String())
}
