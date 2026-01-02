package cache

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestCache(t *testing.T) {
	// Create temp directory for test database
	tmpDir, err := os.MkdirTemp("", "cache_test")
	if err != nil {
		t.Fatalf("create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	c, err := New(filepath.Join(tmpDir, "cache"))
	if err != nil {
		t.Fatalf("new cache: %v", err)
	}
	defer c.Close()

	tests := []struct {
		name       string
		provider   string
		model      string
		sourceLang string
		targetLang string
		text       string
		entry      Entry
	}{
		{
			name:       "basic translation",
			provider:   "openai",
			model:      "gpt-4",
			sourceLang: "en",
			targetLang: "zh",
			text:       "Hello, world!",
			entry: Entry{
				Text: "你好，世界！",
				Usage: Usage{
					PromptTokens:     10,
					CompletionTokens: 5,
					TotalTokens:      15,
				},
				CreatedAt: time.Now(),
			},
		},
		{
			name:       "different model",
			provider:   "gemini",
			model:      "gemini-2.0-flash",
			sourceLang: "zh",
			targetLang: "en",
			text:       "你好",
			entry: Entry{
				Text: "Hello",
				Usage: Usage{
					PromptTokens:     5,
					CompletionTokens: 2,
					TotalTokens:      7,
				},
				CreatedAt: time.Now(),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key := GenerateKey(tt.provider, tt.model, tt.sourceLang, tt.targetLang, tt.text)

			// Should miss initially
			_, found := c.Get(key)
			if found {
				t.Error("expected cache miss, got hit")
			}

			// Set and get
			if err := c.Set(key, &tt.entry, DefaultTTL); err != nil {
				t.Fatalf("set: %v", err)
			}

			got, found := c.Get(key)
			if !found {
				t.Fatal("expected cache hit, got miss")
			}

			if got.Text != tt.entry.Text {
				t.Errorf("text = %q, want %q", got.Text, tt.entry.Text)
			}
			if got.Usage.TotalTokens != tt.entry.Usage.TotalTokens {
				t.Errorf("total tokens = %d, want %d", got.Usage.TotalTokens, tt.entry.Usage.TotalTokens)
			}
		})
	}
}

func TestGenerateKey(t *testing.T) {
	// Same inputs should produce same key
	key1 := GenerateKey("openai", "gpt-4", "en", "zh", "Hello")
	key2 := GenerateKey("openai", "gpt-4", "en", "zh", "Hello")
	if key1 != key2 {
		t.Errorf("same inputs produced different keys: %s vs %s", key1, key2)
	}

	// Different inputs should produce different keys
	key3 := GenerateKey("openai", "gpt-4", "en", "zh", "World")
	if key1 == key3 {
		t.Error("different inputs produced same key")
	}

	// Different model should produce different key
	key4 := GenerateKey("openai", "gpt-3.5", "en", "zh", "Hello")
	if key1 == key4 {
		t.Error("different model produced same key")
	}
}

func TestGenerateKeyNormalization(t *testing.T) {
	tests := []struct {
		name  string
		text1 string
		text2 string
		want  bool // true = should match, false = should differ
	}{
		{
			name:  "leading/trailing whitespace",
			text1: "Hello",
			text2: "  Hello  ",
			want:  true,
		},
		{
			name:  "multiple spaces",
			text1: "Hello World",
			text2: "Hello    World",
			want:  true,
		},
		{
			name:  "different line endings",
			text1: "Hello\nWorld",
			text2: "Hello\r\nWorld",
			want:  true,
		},
		{
			name:  "tabs vs spaces",
			text1: "Hello World",
			text2: "Hello\tWorld",
			want:  true,
		},
		{
			name:  "unicode NFC vs NFD (é)",
			text1: "caf\u00e9",  // é as single codepoint (NFC)
			text2: "cafe\u0301", // e + combining accent (NFD)
			want:  true,
		},
		{
			name:  "actual different content",
			text1: "Hello",
			text2: "Goodbye",
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key1 := GenerateKey("p", "m", "en", "zh", tt.text1)
			key2 := GenerateKey("p", "m", "en", "zh", tt.text2)

			match := key1 == key2
			if match != tt.want {
				t.Errorf("keys match = %v, want %v\n  text1: %q\n  text2: %q", match, tt.want, tt.text1, tt.text2)
			}
		})
	}
}

func TestStats(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cache_stats_test")
	if err != nil {
		t.Fatalf("create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	c, err := New(filepath.Join(tmpDir, "cache"))
	if err != nil {
		t.Fatalf("new cache: %v", err)
	}
	defer c.Close()

	key := GenerateKey("test", "model", "en", "zh", "test")

	// Initial stats should be zero
	stats := c.Stats()
	if stats.Hits != 0 || stats.Misses != 0 {
		t.Errorf("initial stats: hits=%d, misses=%d, want 0,0", stats.Hits, stats.Misses)
	}

	// Miss
	c.Get(key)
	stats = c.Stats()
	if stats.Misses != 1 {
		t.Errorf("misses = %d, want 1", stats.Misses)
	}

	// Set and hit
	entry := &Entry{Text: "test", CreatedAt: time.Now()}
	c.Set(key, entry, DefaultTTL)
	c.Get(key)

	stats = c.Stats()
	if stats.Hits != 1 {
		t.Errorf("hits = %d, want 1", stats.Hits)
	}

	// Check hit rate
	rate := stats.HitRate()
	expected := 50.0 // 1 hit, 1 miss
	if rate != expected {
		t.Errorf("hit rate = %.2f%%, want %.2f%%", rate, expected)
	}
}
