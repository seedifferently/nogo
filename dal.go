package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/boltdb/bolt"
)

var errRecordNotFound = errors.New("record not found")

// Record represents a hosts record
type Record struct {
	Paused bool `json:"paused"`
}

func (r *Record) isAllowed() bool {
	// A paused record is an allowed record
	return r.Paused
}

func (r *Record) jsonEncode() ([]byte, error) {
	data, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (r *Record) jsonDecode(data []byte) (*Record, error) {
	if err := json.Unmarshal(data, &r); err != nil {
		return nil, err
	}

	return r, nil
}

func (r *Record) Bind(req *http.Request) error {
	return nil
}

// APISetting represents an API setting.
type APISetting struct {
	Disabled bool `json:"disabled"`
}

func (s *APISetting) Bind(req *http.Request) error {
	return nil
}

func (db *DB) keyCount() (int, error) {
	var stats bolt.BucketStats

	if err := db.View(func(tx *bolt.Tx) error {
		// Get bucket stats
		stats = tx.Bucket(blacklistKey).Stats()

		return nil
	}); err != nil {
		return 0, err
	}

	return stats.KeyN, nil
}

func (db *DB) get(key string) (*Record, error) {
	var r *Record

	err := db.View(func(tx *bolt.Tx) error {
		var err error

		v := tx.Bucket(blacklistKey).Get([]byte(strings.ToLower(key)))
		if v == nil {
			return errRecordNotFound
		} else if len(v) == 0 {
			// Empty value (likely due to hosts import)
			r = &Record{}
			return nil
		}

		r, err = r.jsonDecode(v)
		return err
	})
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (db *DB) put(key string, r *Record) error {
	err := db.Update(func(tx *bolt.Tx) error {
		var err error
		var v []byte

		// Check for record data
		if r != nil {
			v, err = r.jsonEncode()
			if err != nil {
				return err
			}
		}

		return tx.Bucket(blacklistKey).Put([]byte(strings.ToLower(key)), v)
	})

	return err
}

func (db *DB) delete(key string) error {
	err := db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(blacklistKey).Delete([]byte(strings.ToLower(key)))
	})

	return err
}

func (db *DB) find(search string) map[string]*Record {
	var err error
	var recs = make(map[string]*Record)

	if search != "" {
		search = strings.ToLower(search)

		db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket(blacklistKey)

			b.ForEach(func(k, v []byte) error {
				var r *Record

				if v == nil {
					// Skip "sub-buckets"
					return nil
				}

				// Convert key to string for processing
				sk := string(k)

				if strings.Contains(sk, search) {
					if len(v) == 0 {
						// Empty value (likely due to hosts import)
						r = &Record{}
					} else {
						if r, err = r.jsonDecode(v); err != nil {
							// Log the decode error and continue
							log.Printf("Record.jsonDecode(%s) Error: %s\n", sk, err)
							r = &Record{}
						}
					}

					recs[sk] = r
				}

				return nil
			})

			return nil
		})
	}

	return recs
}

func (db *DB) getPaused() map[string]*Record {
	var err error
	var recs = make(map[string]*Record)

	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(blacklistKey)

		b.ForEach(func(k, v []byte) error {
			var r *Record

			if v == nil || len(v) == 0 {
				// Skip "sub-buckets" and empty values
				return nil
			}

			if r, err = r.jsonDecode(v); err != nil {
				// Log the decode error and continue
				log.Printf("Record.jsonDecode(%s) Error: %s\n", k, err)
				return nil
			}

			if r.Paused {
				recs[string(k)] = r
			}

			return nil
		})

		return nil
	})

	return recs
}

func (db *DB) importBlacklist(fname string) error {
	f, err := os.Open(fname)
	if err != nil {
		return err
	}
	defer f.Close()

	lines, err := lineCount(f)
	if err != nil {
		return err
	}

	// Rewind
	if _, err = f.Seek(0, 0); err != nil {
		return err
	}

	n := 0
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		n++
		fmt.Printf("\rProcessing %d of %d", n, lines)

		// Parse line and trim any dots
		r := strings.Trim(parseRecord(scanner.Text()), ".")

		// Ignore records that don't appear to be valid
		if !isValidDomainName(r) {
			continue
		}

		db.put(r, nil)
	}

	fmt.Print("\n")

	return scanner.Err()
}

func parseRecord(s string) string {
	// Ignore comments
	i := strings.IndexByte(s, '#')
	if i == 0 {
		return ""
	} else if i > 0 {
		s = s[:i]
	}

	sf := strings.Fields(s)

	if len(sf) < 1 {
		// empty
		return ""
	} else if len(sf) > 1 {
		// Return 2nd item if more than 1
		return sf[1]
	} else {
		// Return one and only item
		return sf[0]
	}
}

func lineCount(r io.Reader) (int, error) {
	buf := make([]byte, 32*1024)
	sep := []byte{'\n'}
	count := 0

	for {
		c, err := r.Read(buf)
		count += bytes.Count(buf[:c], sep)

		switch {
		case err == io.EOF:
			return count, nil
		case err != nil:
			return count, err
		}
	}
}
