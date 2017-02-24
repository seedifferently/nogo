package main

import (
	"math/rand"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"testing"
	"time"

	"github.com/boltdb/bolt"
)

func TestMain(m *testing.M) {
	db = MustOpenDB()
	exitVal := m.Run()
	db.MustClose()

	os.Exit(exitVal)
}

func randInt() int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Intn(10000)
}

func tempfilePath(prefix string) string {
	var name string
	dir := os.TempDir()
	conflict := true

	for i := 0; i < 10000; i++ {
		name = filepath.Join(dir, prefix+strconv.Itoa(randInt()))

		if _, err := os.Stat(name); os.IsNotExist(err) {
			conflict = false
			break
		}
	}

	if conflict {
		panic("couldn't find a suitable tempfile path")
	}

	return name
}

func MustOpenDB() *DB {
	bdb, err := bolt.Open(tempfilePath("nogo-db-"), 0666, nil)
	if err != nil {
		panic(err)
	}

	return &DB{bdb}
}

func (db *DB) Reset() {
	db.Update(func(tx *bolt.Tx) error {
		// Delete bucket
		tx.DeleteBucket(blacklistKey)
		return nil
	})

	if err := db.Update(func(tx *bolt.Tx) error {
		// Create bucket
		_, err := tx.CreateBucket(blacklistKey)
		return err
	}); err != nil {
		panic(err)
	}
}

func (db *DB) MustClose() {
	defer os.Remove(db.Path())

	if err := db.Close(); err != nil {
		panic(err)
	}
}

func testEqual(t *testing.T, msg string, args ...interface{}) bool {
	if !reflect.DeepEqual(args[len(args)-2], args[len(args)-1]) {
		t.Errorf(msg, args...)
		return false
	}
	return true
}
