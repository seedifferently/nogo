package main

import (
	"encoding/base64"
	"html/template"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/boltdb/bolt"
	"github.com/pressly/chi"
	"github.com/pressly/chi/render"
)

// H represents a map[string]interface{}
type H map[string]interface{}

// HTTP basic auth middleware
func basicAuth(password string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)

			s := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
			if len(s) != 2 {
				http.Error(w, http.StatusText(401), 401)
				return
			}

			b, err := base64.StdEncoding.DecodeString(s[1])
			if err != nil {
				log.Printf("base64.StdEncoding.DecodeString() Error: %s\n", err)
				http.Error(w, http.StatusText(401), 401)
				return
			}

			pair := strings.SplitN(string(b), ":", 2)
			if len(pair) != 2 {
				log.Printf("strings.SplitN() Error: %s\n", err)
				http.Error(w, http.StatusText(401), 401)
				return
			}

			if pair[0] != "admin" || pair[1] != password {
				http.Error(w, http.StatusText(401), 401)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// GET / (root index)
func rootIndexHandler(w http.ResponseWriter, r *http.Request) {
	var data map[string]*Record
	q := r.FormValue("q")
	p := r.FormValue("p")

	if q != "" { // GET /?q=query
		if len(q) >= 3 {
			data = db.find(q)
		} else {
			http.Error(w, http.StatusText(422), 422)
			return
		}
	} else if p == "1" { // GET /?p=1
		data = db.getPaused()
	}

	totalCount, err := db.keyCount()
	if err != nil {
		log.Printf("db.keyCount() Error: %s\n", err)
		http.Error(w, http.StatusText(500), 500)
		return
	}

	tmpl, err := template.New("index").Parse(indexTmpl)
	if err != nil {
		log.Printf("template.ParseFiles() Error: %s\n", err)
		http.Error(w, http.StatusText(500), 500)
		return
	}

	if err = tmpl.Execute(w, H{"data": data, "isDisabled": isDisabled, "totalCount": totalCount, "q": q, "p": p}); err != nil {
		log.Printf("tmpl.Execute() Error: %s\n", err)
		http.Error(w, http.StatusText(500), 500)
	}
}

// POST /records/
func recordsCreateHandler(w http.ResponseWriter, r *http.Request) {
	var rec *Record

	key := r.FormValue("key")
	if len(key) < 4 {
		http.Error(w, http.StatusText(422), 422)
		return
	}

	p := r.FormValue("paused")
	if p == "1" {
		rec = &Record{Paused: true}
	}

	// Save
	if err := db.put(key, rec); err != nil {
		log.Printf("db.put(%s) Error: %s\n", key, err)
		http.Error(w, http.StatusText(500), 500)
		return
	}

	// Redirect to key view
	http.Redirect(w, r, strings.Join([]string{"/records/", key}, ""), 302)
}

// GET /records/:key
func recordsReadHandler(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "key")

	rec, err := db.get(key)
	if err == errRecordNotFound {
		http.Error(w, http.StatusText(404), 404)
		return
	} else if err != nil {
		log.Printf("db.get(%s) Error: %s\n", key, err)
		http.Error(w, http.StatusText(500), 500)
		return
	}
	data := map[string]*Record{key: rec}

	totalCount, err := db.keyCount()
	if err != nil {
		log.Printf("db.keyCount() Error: %s\n", err)
		http.Error(w, http.StatusText(500), 500)
		return
	}

	tmpl, err := template.New("index").Parse(indexTmpl)
	if err != nil {
		log.Printf("template.ParseFiles() Error: %s\n", err)
		http.Error(w, http.StatusText(500), 500)
		return
	}

	if err = tmpl.Execute(w, H{"data": data, "isDisabled": isDisabled, "totalCount": totalCount}); err != nil {
		log.Printf("tmpl.Execute() Error: %s\n", err)
		http.Error(w, http.StatusText(500), 500)
	}
}

// GET /export/hosts.txt
func exportHostsHandler(w http.ResponseWriter, r *http.Request) {
	var bol = []byte("0.0.0.0 ") // Beginning of each line
	var eol = []byte("\n")       // End of each line

	w.Header().Set("Content-Disposition", "attachment; filename=hosts.txt")
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	// Begin with localhost records
	if _, err := w.Write([]byte("# Exported from nogo (http://nogo.curia.solutions/)\n127.0.0.1 localhost\n127.0.0.1 localhost.localdomain\n127.0.0.1 local\n255.255.255.255 broadcasthost\n::1 localhost\nfe80::1%lo0 localhost\n\n")); err != nil {
		log.Printf("Write Error: %s\n", err)
		http.Error(w, http.StatusText(500), 500)
		return
	}

	// Export each host from the db
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(blacklistKey)

		return b.ForEach(func(k, v []byte) error {
			if v == nil {
				// Skip "sub-buckets"
				return nil
			}

			line := bol
			line = append(line, k...)
			line = append(line, eol...)
			if _, err := w.Write(line); err != nil {
				return err
			}

			return nil
		})
	})
	if err != nil {
		log.Printf("Export Error: %s\n", err)
		http.Error(w, http.StatusText(500), 500)
		return
	}
}

// GET /api/records/
func apiRecordsIndexHandler(w http.ResponseWriter, r *http.Request) {
	var data map[string]*Record

	q := r.FormValue("q")
	p := r.FormValue("p")
	if len(q) >= 3 { // GET /api/records/?q=query
		data = db.find(q)
	} else if p == "1" { // GET /api/records/?p=1
		data = db.getPaused()
	} else {
		http.Error(w, http.StatusText(422), 422)
		return
	}

	render.JSON(w, r, H{"data": data})
}

// GET /api/records/:key
func apiRecordsReadHandler(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "key")

	data, err := db.get(key)
	if err == errRecordNotFound {
		http.Error(w, http.StatusText(404), 404)
		return
	} else if err != nil {
		log.Printf("db.get(%s) Error: %s\n", key, err)
		http.Error(w, http.StatusText(500), 500)
		return
	}

	render.JSON(w, r, H{"data": H{key: data}})
}

// PUT /api/records/:key
func apiRecordsUpdateHandler(w http.ResponseWriter, r *http.Request) {
	var data Record

	key := chi.URLParam(r, "key")
	if len(key) < 4 {
		http.Error(w, http.StatusText(422), 422)
		return
	}

	// Bind
	if err := render.Bind(r.Body, &data); err != nil && err != io.EOF {
		log.Printf("render.Bind() Error: %s\n", err)
		http.Error(w, http.StatusText(400), 400)
		return
	}

	// Save
	if err := db.put(key, &data); err != nil {
		log.Printf("db.put(%s) Error: %s\n", key, err)
		http.Error(w, http.StatusText(500), 500)
		return
	}

	render.JSON(w, r, H{"data": H{key: data}})
}

// DELETE /api/records/:key
func apiRecordsDeleteHandler(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "key")

	// Delete
	if err := db.delete(key); err != nil {
		log.Printf("db.delete(%s) Error: %s\n", key, err)
		http.Error(w, http.StatusText(500), 500)
		return
	}

	render.NoContent(w, r)
}

// PUT /api/settings/
func apiSettingsUpdateHandler(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Disabled bool `json:"disabled"`
	}

	// Bind
	if err := render.Bind(r.Body, &data); err != nil {
		log.Printf("render.Bind() Error: %s\n", err)
		http.Error(w, http.StatusText(400), 400)
		return
	}

	// Update disabled toggle
	isDisabled = data.Disabled

	render.JSON(w, r, H{"data": data})
}

// GET /css/nogo.css
func cssHandler(w http.ResponseWriter, r *http.Request) {
	var data = []byte(nogoCSS)

	w.Header().Set("Cache-Control", "public, max-age=31536000")
	w.Header().Set("Content-Length", strconv.Itoa(len(data)))
	w.Header().Set("Content-Type", "text/css")

	w.Write(data)
}
