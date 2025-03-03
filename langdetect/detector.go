package langdetect

import (
	"github.com/pemistahl/lingua-go"
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

var detector lingua.LanguageDetector

// 初始化语言检测器
func init() {
	// 支持的语言列表
	languages := []lingua.Language{
		lingua.Chinese,
		lingua.English,
		lingua.Japanese,
		lingua.Korean,
		lingua.French,
		lingua.German,
		lingua.Spanish,
		lingua.Russian,
		lingua.Italian,
		lingua.Portuguese,
		lingua.Arabic,
	}

	// 创建语言检测器
	detector = lingua.NewLanguageDetectorBuilder().
		FromLanguages(languages...).
		WithPreloadedLanguageModels().
		Build()
}

// DetectLanguage 检测文本的语言
func DetectLanguage(text string) (string, string) {
	if text == "" {
		return "auto", ""
	}

	// 检测语言
	language, exists := detector.DetectLanguageOf(text)
	if !exists {
		return "auto", ""
	}

	// 返回语言代码和语言名称
	return mapLanguageCode(language), mapLanguageName(language)
}

// 将 lingua 语言映射到我们的语言代码
func mapLanguageCode(language lingua.Language) string {
	switch language {
	case lingua.Chinese:
		return "zh"
	case lingua.English:
		return "en"
	case lingua.Japanese:
		return "ja"
	case lingua.Korean:
		return "ko"
	case lingua.French:
		return "fr"
	case lingua.German:
		return "de"
	case lingua.Spanish:
		return "es"
	case lingua.Russian:
		return "ru"
	case lingua.Italian:
		return "it"
	case lingua.Portuguese:
		return "pt"
	case lingua.Arabic:
		return "ar"
	default:
		return "auto"
	}
}

// 将 lingua 语言映射到语言名称
func mapLanguageName(language lingua.Language) string {
	switch language {
	case lingua.Chinese:
		return "中文"
	case lingua.English:
		return "英语"
	case lingua.Japanese:
		return "日语"
	case lingua.Korean:
		return "韩语"
	case lingua.French:
		return "法语"
	case lingua.German:
		return "德语"
	case lingua.Spanish:
		return "西班牙语"
	case lingua.Russian:
		return "俄语"
	case lingua.Italian:
		return "意大利语"
	case lingua.Portuguese:
		return "葡萄牙语"
	case lingua.Arabic:
		return "阿拉伯语"
	default:
		return ""
	}
}
