package main

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/boltdb/bolt"
)

func TestRecord_isAllowed(t *testing.T) {
	r := &Record{}
	testEqual(t, "isAllowed() = %+v, want %+v", r.isAllowed(), false)

	r = &Record{Paused: true}
	testEqual(t, "isAllowed() = %+v, want %+v", r.isAllowed(), true)
}

func TestRecord_jsonEncode(t *testing.T) {
	r := &Record{}
	d, _ := r.jsonEncode()
	testEqual(t, "jsonEncode() = %+v, want %+v", string(d), "{\"Paused\":false}")

	r = &Record{Paused: true}
	d, _ = r.jsonEncode()
	testEqual(t, "jsonEncode() = %+v, want %+v", string(d), "{\"Paused\":true}")
}

func TestRecord_jsonDecode(t *testing.T) {
	var r *Record

	r, _ = r.jsonDecode([]byte("{\"Paused\":false}"))
	testEqual(t, "jsonDecode() = %+v, want %+v", *r, Record{})

	r, _ = r.jsonDecode([]byte("{\"Paused\":true}"))
	testEqual(t, "jsonDecode() = %+v, want %+v", *r, Record{Paused: true})
}

func TestDB_keyCount(t *testing.T) {
	var r *Record
	db.Reset()

	c, _ := db.keyCount()
	testEqual(t, "keyCount() = %+v, want %+v", c, 0)

	if err := db.put("test.test", r); err != nil {
		t.Errorf("failed to put: %+v", err)
	}
	c, _ = db.keyCount()
	testEqual(t, "keyCount() = %+v, want %+v", c, 1)
}

func TestDB_put_get(t *testing.T) {
	var r *Record
	db.Reset()

	if err := db.put("Nil.test", r); err != nil {
		t.Errorf("failed to put: %+v", err)
	}
	db.View(func(tx *bolt.Tx) error {
		v := tx.Bucket(blacklistKey).Get([]byte("nil.test"))
		testEqual(t, "put() = %+v, want %+v", string(v), "")
		return nil
	})
	r, _ = db.get("nil.Test")
	testEqual(t, "get() = %+v, want %+v", *r, Record{})

	if err := db.put("empty.test", &Record{}); err != nil {
		t.Errorf("failed to put: %+v", err)
	}
	db.View(func(tx *bolt.Tx) error {
		v := tx.Bucket(blacklistKey).Get([]byte("empty.test"))
		testEqual(t, "put() = %+v, want %+v", string(v), "{\"Paused\":false}")
		return nil
	})
	r, _ = db.get("empty.test")
	testEqual(t, "get() = %+v, want %+v", *r, Record{})

	if err := db.put("paused.test", &Record{Paused: true}); err != nil {
		t.Errorf("failed to put: %+v", err)
	}
	db.View(func(tx *bolt.Tx) error {
		v := tx.Bucket(blacklistKey).Get([]byte("paused.test"))
		testEqual(t, "put() = %+v, want %+v", string(v), "{\"Paused\":true}")
		return nil
	})
	r, _ = db.get("paused.test")
	testEqual(t, "get() = %+v, want %+v", *r, Record{Paused: true})
}

func TestDB_delete(t *testing.T) {
	db.Reset()

	if err := db.put("delete.test", &Record{}); err != nil {
		t.Errorf("failed to put: %+v", err)
	}
	r, _ := db.get("delete.test")
	testEqual(t, "get() = %+v, want %+v", *r, Record{})
	if err := db.delete("delete.test"); err != nil {
		t.Errorf("failed to delete: %+v", err)
	}
	_, err := db.get("delete.test")
	testEqual(t, "get() err = %+v, want %+v", err, errRecordNotFound)
}

func TestDB_find(t *testing.T) {
	var r *Record
	db.Reset()

	if err := db.put("one.abcd", r); err != nil {
		t.Errorf("failed to put: %+v", err)
	}
	if err := db.put("two.abcd", r); err != nil {
		t.Errorf("failed to put: %+v", err)
	}
	if err := db.put("one.asdf", r); err != nil {
		t.Errorf("failed to put: %+v", err)
	}
	if err := db.put("one.arst", r); err != nil {
		t.Errorf("failed to put: %+v", err)
	}

	rs := db.find("one")
	testEqual(t, "len(find('one')) = %+v, want %+v", len(rs), 3)
	rs = db.find("abcd")
	testEqual(t, "len(find('abcd')) = %+v, want %+v", len(rs), 2)
	rs = db.find("arst")
	testEqual(t, "len(find('arst')) = %+v, want %+v", len(rs), 1)
}

func TestDB_getPaused(t *testing.T) {
	var r *Record
	db.Reset()

	if err := db.put("Nil.test", r); err != nil {
		t.Errorf("failed to put: %+v", err)
	}
	if err := db.put("empty.test", &Record{}); err != nil {
		t.Errorf("failed to put: %+v", err)
	}
	if err := db.put("paused.test", &Record{Paused: true}); err != nil {
		t.Errorf("failed to put: %+v", err)
	}

	rs := db.getPaused()
	testEqual(t, "len(getPaused()) = %+v, want %+v", len(rs), 1)
	testEqual(t, "getPaused()[0] = %+v, want %+v", *rs["paused.test"], Record{Paused: true})
}

func TestDB_importBlacklist(t *testing.T) {
	db.Reset()

	f, err := ioutil.TempFile("", "nogo-import-")
	if err != nil {
		t.Errorf("failed to create TempFile: %+v", err)
	}
	f.WriteString(" # comment\n.\ntst\ntest\n.test\ntes.t\ntest.test\n")
	f.Sync()
	defer f.Close()
	defer os.Remove(f.Name())

	db.importBlacklist(f.Name())

	c, err := db.keyCount()
	if err != nil {
		t.Errorf("failed to get keyCount: %+v", err)
	}
	testEqual(t, "keyCount() = %+v, want %+v", c, 1)

	r, _ := db.get("test.test")
	testEqual(t, "get('test.test') = %+v, want %+v", *r, Record{})
}

func Test_parseRecord(t *testing.T) {
	testEqual(t, "parseRecord('# comment') = %+v, want %+v", parseRecord("# comment"), "")
	testEqual(t, "parseRecord(' ') = %+v, want %+v", parseRecord(" "), "")
	testEqual(t, "parseRecord('partial # comment') = %+v, want %+v", parseRecord("partial # comment"), "partial")
	testEqual(t, "parseRecord('127.0.0.1\tlocalhost # comment') = %+v, want %+v", parseRecord("127.0.0.1\tlocalhost # comment"), "localhost")
	testEqual(t, "parseRecord('127.0.0.1 localhost alias # comment') = %+v, want %+v", parseRecord("127.0.0.1 localhost alias # comment"), "localhost")
}

func TestDB_lineCount(t *testing.T) {
	c, err := lineCount(strings.NewReader("one\ntwo\nthree\n"))
	if err != nil {
		t.Errorf("failed to get lineCount: %+v", err)
	}
	testEqual(t, "lineCount('one\\ntwo\\nthree\\n') = %+v, want %+v", c, 3)
}
