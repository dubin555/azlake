package azcat

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"github.com/dgraph-io/badger/v4"
)

var ErrNotFound = errors.New("not found")

// KVStore wraps BadgerDB for simple key-value operations
type KVStore struct {
	db *badger.DB
}

// OpenKV opens or creates a BadgerDB at the given path
func OpenKV(dataDir string) (*KVStore, error) {
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, err
	}
	opts := badger.DefaultOptions(dataDir).
		WithLogger(nil) // suppress badger's verbose logging
	db, err := badger.Open(opts)
	if err != nil {
		return nil, err
	}
	return &KVStore{db: db}, nil
}

func (s *KVStore) Close() error {
	return s.db.Close()
}

// Get retrieves a value by key
func (s *KVStore) Get(key string) ([]byte, error) {
	var val []byte
	err := s.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			if errors.Is(err, badger.ErrKeyNotFound) {
				return ErrNotFound
			}
			return err
		}
		val, err = item.ValueCopy(nil)
		return err
	})
	return val, err
}

// Set stores a key-value pair
func (s *KVStore) Set(key string, value []byte) error {
	return s.db.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(key), value)
	})
}

// Delete removes a key
func (s *KVStore) Delete(key string) error {
	return s.db.Update(func(txn *badger.Txn) error {
		return txn.Delete([]byte(key))
	})
}

// Scan iterates over keys with the given prefix, calling fn for each.
// Keys are returned in lexicographic order. Pass after="" to start from beginning.
// Returns up to limit results.
func (s *KVStore) Scan(prefix string, after string, limit int, fn func(key string, value []byte) error) error {
	return s.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.Prefix = []byte(prefix)
		it := txn.NewIterator(opts)
		defer it.Close()

		start := []byte(prefix)
		if after != "" {
			start = []byte(after + "\x00") // seek past 'after'
		}

		count := 0
		for it.Seek(start); it.Valid(); it.Next() {
			if limit > 0 && count >= limit {
				break
			}
			item := it.Item()
			k := string(item.Key())
			v, err := item.ValueCopy(nil)
			if err != nil {
				return err
			}
			if err := fn(k, v); err != nil {
				return err
			}
			count++
		}
		return nil
	})
}

// DeletePrefix removes all keys with the given prefix
func (s *KVStore) DeletePrefix(prefix string) error {
	return s.db.Update(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.Prefix = []byte(prefix)
		opts.PrefetchValues = false
		it := txn.NewIterator(opts)
		defer it.Close()

		var keys [][]byte
		for it.Seek([]byte(prefix)); it.Valid(); it.Next() {
			keys = append(keys, it.Item().KeyCopy(nil))
		}
		for _, k := range keys {
			if err := txn.Delete(k); err != nil {
				return err
			}
		}
		return nil
	})
}

// JSON helpers

func (s *KVStore) GetJSON(key string, v interface{}) error {
	data, err := s.Get(key)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}

func (s *KVStore) SetJSON(key string, v interface{}) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return s.Set(key, data)
}

// DefaultDataDir returns ~/.azlake/data
func DefaultDataDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".azlake", "data")
}

// DefaultObjectsDir returns ~/.azlake/objects
func DefaultObjectsDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".azlake", "objects")
}
