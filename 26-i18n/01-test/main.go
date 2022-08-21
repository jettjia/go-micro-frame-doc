package main

import (
	"fmt"

	"golang.org/x/text/language"

	"github.com/nicksnyder/go-i18n/v2/i18n"
)

func main() {
	bundle := i18n.NewBundle(language.English)
	localizer := i18n.NewLocalizer(bundle, "en")
	catsMessage := &i18n.Message{
		ID:    "Cats",
		One:   "I have {{.PluralCount}} cat.",
		Other: "I have {{.PluralCount}} cats.",
	}
	fmt.Println(localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: catsMessage,
		PluralCount:    1,
	}))
	fmt.Println(localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: catsMessage,
		PluralCount:    2,
	}))
}
