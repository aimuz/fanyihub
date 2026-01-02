// Package langdetect provides language detection using lingua-go.
package langdetect

import (
	"github.com/pemistahl/lingua-go"
	// Language model imports for lingua
	_ "github.com/pemistahl/lingua-go/language-models/ar"
	_ "github.com/pemistahl/lingua-go/language-models/de"
	_ "github.com/pemistahl/lingua-go/language-models/en"
	_ "github.com/pemistahl/lingua-go/language-models/es"
	_ "github.com/pemistahl/lingua-go/language-models/fr"
	_ "github.com/pemistahl/lingua-go/language-models/it"
	_ "github.com/pemistahl/lingua-go/language-models/ja"
	_ "github.com/pemistahl/lingua-go/language-models/ko"
	_ "github.com/pemistahl/lingua-go/language-models/pt"
	_ "github.com/pemistahl/lingua-go/language-models/ru"
	_ "github.com/pemistahl/lingua-go/language-models/zh"
)

// langInfo holds language code and display name.
type langInfo struct {
	code string
	name string
}

// languageMap maps lingua.Language to our language info (table-driven).
var languageMap = map[lingua.Language]langInfo{
	lingua.Chinese:    {"zh", "中文"},
	lingua.English:    {"en", "英语"},
	lingua.Japanese:   {"ja", "日语"},
	lingua.Korean:     {"ko", "韩语"},
	lingua.French:     {"fr", "法语"},
	lingua.German:     {"de", "德语"},
	lingua.Spanish:    {"es", "西班牙语"},
	lingua.Russian:    {"ru", "俄语"},
	lingua.Italian:    {"it", "意大利语"},
	lingua.Portuguese: {"pt", "葡萄牙语"},
	lingua.Arabic:     {"ar", "阿拉伯语"},
}

// supportedLanguages extracts the list of supported languages from the map.
func supportedLanguages() []lingua.Language {
	langs := make([]lingua.Language, 0, len(languageMap))
	for lang := range languageMap {
		langs = append(langs, lang)
	}
	return langs
}

var detector lingua.LanguageDetector

func init() {
	detector = lingua.NewLanguageDetectorBuilder().
		FromLanguages(supportedLanguages()...).
		WithPreloadedLanguageModels().
		Build()
}

// Detect detects the language of the given text.
// Returns language code and display name.
// If detection fails, returns ("auto", "").
func Detect(text string) (code, name string) {
	if text == "" {
		return "auto", ""
	}

	lang, ok := detector.DetectLanguageOf(text)
	if !ok {
		return "auto", ""
	}

	info, ok := languageMap[lang]
	if !ok {
		return "auto", ""
	}

	return info.code, info.name
}
