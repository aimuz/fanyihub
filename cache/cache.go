// Package cache provides LLM response caching with BadgerDB.
package cache

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"sync/atomic"
	"time"

	"github.com/dgraph-io/badger/v4"
	"golang.org/x/text/unicode/norm"
)

// DefaultTTL is the default cache entry time-to-live.
const DefaultTTL = 7 * 24 * time.Hour // 7 days

// Entry represents a cached LLM response.
type Entry struct {
	Text      string    `json:"text"`
	Usage     Usage     `json:"usage"`
	CreatedAt time.Time `json:"created_at"`
}

// Usage mirrors types.Usage to avoid import cycle.
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// Stats holds cache statistics.
type Stats struct {
	Hits   uint64 `json:"hits"`
	Misses uint64 `json:"misses"`
}

// HitRate returns the cache hit rate as a percentage.
func (s Stats) HitRate() float64 {
	total := s.Hits + s.Misses
	if total == 0 {
		return 0
	}
	return float64(s.Hits) / float64(total) * 100
}

// Cache wraps BadgerDB for LLM response caching.
type Cache struct {
	db     *badger.DB
	hits   atomic.Uint64
	misses atomic.Uint64
}

// New creates a new cache at the given path.
func New(path string) (*Cache, error) {
	opts := badger.DefaultOptions(path)
	opts.Logger = nil // Disable BadgerDB internal logging

	db, err := badger.Open(opts)
	if err != nil {
		return nil, fmt.Errorf("open badger: %w", err)
	}

	c := &Cache{db: db}

	// Start background GC goroutine
	go c.runGC()

	return c, nil
}

// runGC periodically runs BadgerDB garbage collection.
func (c *Cache) runGC() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		if c.db.IsClosed() {
			return
		}
		_ = c.db.RunValueLogGC(0.5)
	}
}

// whitespaceRe matches one or more whitespace characters.
var whitespaceRe = regexp.MustCompile(`\s+`)

// GenerateKey creates a cache key from translation parameters.
// The text is normalized before hashing to improve cache hit rate.
func GenerateKey(provider, model, sourceLang, targetLang, text string) string {
	normalized := normalizeText(text)
	data := fmt.Sprintf("%s|%s|%s|%s|%s", provider, model, sourceLang, targetLang, normalized)
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

// normalizeText applies transformations to improve cache hit rate:
//   - Trims leading/trailing whitespace
//   - Collapses multiple whitespace characters to single space
//   - Normalizes Unicode to NFC form (composed characters)
//   - Normalizes line endings to \n
func normalizeText(s string) string {
	// Unicode NFC normalization (e.g., é vs e+́)
	s = norm.NFC.String(s)

	// Normalize line endings: \r\n -> \n, \r -> \n
	s = strings.ReplaceAll(s, "\r\n", "\n")
	s = strings.ReplaceAll(s, "\r", "\n")

	// Collapse multiple whitespace to single space
	s = whitespaceRe.ReplaceAllString(s, " ")

	// Trim leading/trailing whitespace
	s = strings.TrimSpace(s)

	return s
}

// Get retrieves an entry from the cache.
// Returns nil and false if not found.
func (c *Cache) Get(key string) (*Entry, bool) {
	var entry Entry

	err := c.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}

		return item.Value(func(val []byte) error {
			return json.Unmarshal(val, &entry)
		})
	})

	if err != nil {
		c.misses.Add(1)
		return nil, false
	}

	c.hits.Add(1)
	return &entry, true
}

// Set stores an entry in the cache with the given TTL.
func (c *Cache) Set(key string, entry *Entry, ttl time.Duration) error {
	if ttl == 0 {
		ttl = DefaultTTL
	}

	data, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("marshal entry: %w", err)
	}

	return c.db.Update(func(txn *badger.Txn) error {
		e := badger.NewEntry([]byte(key), data).WithTTL(ttl)
		return txn.SetEntry(e)
	})
}

// Stats returns current cache statistics.
func (c *Cache) Stats() Stats {
	return Stats{
		Hits:   c.hits.Load(),
		Misses: c.misses.Load(),
	}
}

// Close closes the cache database.
func (c *Cache) Close() error {
	if c.db != nil {
		return c.db.Close()
	}
	return nil
}
