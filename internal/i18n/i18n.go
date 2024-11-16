package i18n

import (
	"embed"
	"encoding/json"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

//go:embed localization/*.json
var localesFS embed.FS
var localizer *i18n.Localizer

func InitLocalizer() {
	bundle := i18n.NewBundle(language.Russian)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)

	files := []string{"ru.json", "en.json"}
	for _, file := range files {
		data, err := localesFS.ReadFile("localization/" + file)
		if err != nil {
			panic(err)
		}
		bundle.MustParseMessageFileBytes(data, file)
	}

	localizer = i18n.NewLocalizer(bundle, "ru")
}

func Localize(messageID string) string {
	return localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: messageID})
}
