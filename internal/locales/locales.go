package locales

import (
	"encoding/json"
	"net/url"
	"os"
)

type Locale = map[string]string
type Locales = map[string]Locale

func MustLocales(p string) Locales {
	t := make(Locales)
	langs := []string{"en", "ru", "de", "ja"}
	for _, lang := range langs {
		path, err := url.JoinPath(p, lang)
		if err != nil {
			panic(err)
		}
		data, err := os.ReadFile(path + ".json")
		if err != nil {
			panic(err)
		}
		var langLoc Locale
		json.Unmarshal(data, &langLoc)
		t[lang] = langLoc
	}
	return t
}
