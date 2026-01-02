package langdetect

import "testing"

func TestDetect(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantCode string
		wantName string
	}{
		{"empty", "", "auto", ""},
		{"chinese", "你好世界", "zh", "中文"},
		{"english", "Hello World", "en", "英语"},
		{"japanese", "こんにちは", "ja", "日语"},
		{"korean", "안녕하세요", "ko", "韩语"},
		{"french", "Bonjour le monde", "fr", "法语"},
		{"german", "Guten Morgen", "de", "德语"},
		{"spanish", "Hola mundo", "es", "西班牙语"},
		{"russian", "Привет мир", "ru", "俄语"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, name := Detect(tt.input)
			if code != tt.wantCode {
				t.Errorf("Detect(%q) code = %q, want %q", tt.input, code, tt.wantCode)
			}
			if name != tt.wantName {
				t.Errorf("Detect(%q) name = %q, want %q", tt.input, name, tt.wantName)
			}
		})
	}
}
