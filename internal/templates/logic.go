package templates

import (
	"blog/internal/locales"
	"fmt"
	"html/template"
)

func Functions(name string, path string) (*template.Template, error) {
	funcs := template.FuncMap{
		"add": func(a int, b int) int { return a + b },
		"sub": func(a int, b int) int { return a - b },
		"i18n": func(locale locales.Locale, str string) string {
			translated, ok := locale[str]
			if !ok {
				return "translation error"
			}
			return translated
		},
	}
	t := template.New(name)
	t.Funcs(funcs)
	t, err := t.ParseGlob(path + "/*.html")
	if err != nil {
		return nil, fmt.Errorf("template parse glob: %w", err)
	}
	return t, nil
}
