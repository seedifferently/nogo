package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/pressly/chi"
)

func Test_basicAuth(t *testing.T) {
	db.Reset()

	// Unauthorized (no Authorization header)
	r := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	cr := chi.NewRouter()
	cr.Use(basicAuth("test"))
	cr.Get("/", rootIndexHandler)
	cr.ServeHTTP(w, r)
	testEqual(t, "Response code = %+v, want %+v", w.Code, 401)

	// Unauthorized (wrong password)
	r = httptest.NewRequest("GET", "/", nil)
	r.Header.Set("Authorization", "Basic YWRtaW46YWRtaW4=")
	w = httptest.NewRecorder()
	cr.ServeHTTP(w, r)
	testEqual(t, "Response code = %+v, want %+v", w.Code, 401)

	// Authorized
	r = httptest.NewRequest("GET", "/", nil)
	r.Header.Set("Authorization", "Basic YWRtaW46dGVzdA==")
	w = httptest.NewRecorder()
	cr.ServeHTTP(w, r)
	testEqual(t, "Response code = %+v, want %+v", w.Code, 200)
	testEqual(t, "Body contains '0 total records' = %+v, want %+v", strings.Contains(w.Body.String(), "0 total records"), true)
}

func Test_rootIndexHandler(t *testing.T) {
	db.Reset()

	// No records
	r := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	rootIndexHandler(w, r)
	testEqual(t, "Response code = %+v, want %+v", w.Code, 200)
	testEqual(t, "Body contains '0 total records' = %+v, want %+v", strings.Contains(w.Body.String(), "0 total records"), true)

	// A record
	if err := db.put("test.test", &Record{Paused: true}); err != nil {
		t.Errorf("failed to put: %+v", err)
	}
	r = httptest.NewRequest("GET", "/", nil)
	w = httptest.NewRecorder()
	rootIndexHandler(w, r)
	testEqual(t, "Response code = %+v, want %+v", w.Code, 200)
	testEqual(t, "Content-Type header = %+v, want %+v", w.Header().Get("Content-Type"), "text/html; charset=utf-8")
	testEqual(t, "Body contains '1 total records' = %+v, want %+v", strings.Contains(w.Body.String(), "1 total records"), true)

	// Search record
	r = httptest.NewRequest("GET", "/?q=te", nil)
	w = httptest.NewRecorder()
	rootIndexHandler(w, r)
	testEqual(t, "Response code = %+v, want %+v", w.Code, 422)
	r = httptest.NewRequest("GET", "/?q=tes", nil)
	w = httptest.NewRecorder()
	rootIndexHandler(w, r)
	testEqual(t, "Response code = %+v, want %+v", w.Code, 200)
	testEqual(t, "Content-Type header = %+v, want %+v", w.Header().Get("Content-Type"), "text/html; charset=utf-8")
	testEqual(t, "Body contains 'Found 1 of 1 total records' = %+v, want %+v", strings.Contains(w.Body.String(), "Found <span id=\"data-count\">1</span> of <span id=\"total-count\">1</span> total records"), true)
	testEqual(t, "Body contains 'test.test' = %+v, want %+v", strings.Contains(w.Body.String(), "<div class=\"column key\">test.test</div>"), true)

	// List paused records
	r = httptest.NewRequest("GET", "/?p=1", nil)
	w = httptest.NewRecorder()
	rootIndexHandler(w, r)
	testEqual(t, "Response code = %+v, want %+v", w.Code, 200)
	testEqual(t, "Content-Type header = %+v, want %+v", w.Header().Get("Content-Type"), "text/html; charset=utf-8")
	testEqual(t, "Body contains '1 of 1 total records' = %+v, want %+v", strings.Contains(w.Body.String(), "<span id=\"data-count\">1</span> of <span id=\"total-count\">1</span> total records"), true)
	testEqual(t, "Body contains 'test.test' = %+v, want %+v", strings.Contains(w.Body.String(), "<div class=\"column key\">test.test</div>"), true)
}

func Test_recordsCreateHandler(t *testing.T) {
	db.Reset()

	// Invalid
	r := &http.Request{
		Method: "POST",
		URL:    &url.URL{Path: "/records/"},
		Form:   url.Values{"key": {"tst"}, "paused": {"1"}},
	}
	w := httptest.NewRecorder()
	recordsCreateHandler(w, r)
	testEqual(t, "Response code = %+v, want %+v", w.Code, 422)

	// Valid
	r = &http.Request{
		Method: "POST",
		URL:    &url.URL{Path: "/records/"},
		Form:   url.Values{"key": {"test.test"}, "paused": {"1"}},
	}
	w = httptest.NewRecorder()
	recordsCreateHandler(w, r)
	testEqual(t, "Response code = %+v, want %+v", w.Code, 302)
	testEqual(t, "Location header = %+v, want %+v", w.Header().Get("Location"), "/records/test.test")
	// verify record created in db
	rec, _ := db.get("test.test")
	testEqual(t, "get() = %+v, want %+v", *rec, Record{Paused: true})
}

func Test_recordsReadHandler(t *testing.T) {
	db.Reset()

	// No records
	r := httptest.NewRequest("GET", "/records/test.test", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Set("key", "test.test")
	w := httptest.NewRecorder()
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
	recordsReadHandler(w, r)
	testEqual(t, "Response code = %+v, want %+v", w.Code, 404)

	// A record
	if err := db.put("test.test", &Record{Paused: true}); err != nil {
		t.Errorf("failed to put: %+v", err)
	}
	r = httptest.NewRequest("GET", "/records/test.test", nil)
	rctx = chi.NewRouteContext()
	rctx.URLParams.Set("key", "test.test")
	w = httptest.NewRecorder()
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
	recordsReadHandler(w, r)
	testEqual(t, "Response code = %+v, want %+v", w.Code, 200)
	testEqual(t, "Content-Type header = %+v, want %+v", w.Header().Get("Content-Type"), "text/html; charset=utf-8")
	testEqual(t, "Body contains '1 of 1 total records' = %+v, want %+v", strings.Contains(w.Body.String(), "<span id=\"data-count\">1</span> of <span id=\"total-count\">1</span> total records"), true)
	testEqual(t, "Body contains 'test.test' = %+v, want %+v", strings.Contains(w.Body.String(), "<div class=\"column key\">test.test</div>"), true)
}

func Test_exportHostsHandler(t *testing.T) {
	db.Reset()

	if err := db.put("test.test", &Record{}); err != nil {
		t.Errorf("failed to put: %+v", err)
	}
	r := httptest.NewRequest("GET", "/export/hosts.txt", nil)
	w := httptest.NewRecorder()
	exportHostsHandler(w, r)
	testEqual(t, "Response code = %+v, want %+v", w.Code, 200)
	testEqual(t, "Content-Type header = %+v, want %+v", w.Header().Get("Content-Type"), "text/plain; charset=utf-8")
	testEqual(t, "Content-Disposition header = %+v, want %+v", w.Header().Get("Content-Disposition"), "attachment; filename=hosts.txt")
	testEqual(t, "Body contains '127.0.0.1 localhost' = %+v, want %+v", strings.Contains(w.Body.String(), "127.0.0.1 localhost\n"), true)
	testEqual(t, "Body contains '0.0.0.0 test.test' = %+v, want %+v", strings.Contains(w.Body.String(), "0.0.0.0 test.test\n"), true)
}

func Test_apiRecordsIndexHandler(t *testing.T) {
	db.Reset()
	if err := db.put("test.test", &Record{Paused: true}); err != nil {
		t.Errorf("failed to put: %+v", err)
	}

	// Bare request
	r := httptest.NewRequest("GET", "/api/records/", nil)
	w := httptest.NewRecorder()
	apiRecordsIndexHandler(w, r)
	testEqual(t, "Response code = %+v, want %+v", w.Code, 422)

	// Search record
	r = httptest.NewRequest("GET", "/api/records/?q=te", nil)
	w = httptest.NewRecorder()
	apiRecordsIndexHandler(w, r)
	testEqual(t, "Response code = %+v, want %+v", w.Code, 422)
	r = httptest.NewRequest("GET", "/api/records/?q=tes", nil)
	w = httptest.NewRecorder()
	apiRecordsIndexHandler(w, r)
	testEqual(t, "Response code = %+v, want %+v", w.Code, 200)
	testEqual(t, "Content-Type header = %+v, want %+v", w.Header().Get("Content-Type"), "application/json; charset=utf-8")
	testEqual(t, "Body = %+v, want %+v", w.Body.String(), "{\"data\":{\"test.test\":{\"paused\":true}}}\n")

	// List paused records
	r = httptest.NewRequest("GET", "/api/records/?p=1", nil)
	w = httptest.NewRecorder()
	apiRecordsIndexHandler(w, r)
	testEqual(t, "Response code = %+v, want %+v", w.Code, 200)
	testEqual(t, "Content-Type header = %+v, want %+v", w.Header().Get("Content-Type"), "application/json; charset=utf-8")
	testEqual(t, "Body = %+v, want %+v", w.Body.String(), "{\"data\":{\"test.test\":{\"paused\":true}}}\n")
}

func Test_apiRecordsReadHandler(t *testing.T) {
	db.Reset()

	// No records
	r := httptest.NewRequest("GET", "/api/records/test.test", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Set("key", "test.test")
	w := httptest.NewRecorder()
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
	apiRecordsReadHandler(w, r)
	testEqual(t, "Response code = %+v, want %+v", w.Code, 404)

	// A record
	if err := db.put("test.test", &Record{Paused: true}); err != nil {
		t.Errorf("failed to put: %+v", err)
	}
	r = httptest.NewRequest("GET", "/api/records/test.test", nil)
	rctx = chi.NewRouteContext()
	rctx.URLParams.Set("key", "test.test")
	w = httptest.NewRecorder()
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
	apiRecordsReadHandler(w, r)
	testEqual(t, "Response code = %+v, want %+v", w.Code, 200)
	testEqual(t, "Content-Type header = %+v, want %+v", w.Header().Get("Content-Type"), "application/json; charset=utf-8")
	testEqual(t, "Body = %+v, want %+v", w.Body.String(), "{\"data\":{\"test.test\":{\"paused\":true}}}\n")
}

func Test_apiRecordsUpdateHandler(t *testing.T) {
	db.Reset()

	// Invalid
	r := httptest.NewRequest("GET", "/api/records/tst", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Set("key", "tst")
	w := httptest.NewRecorder()
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
	apiRecordsUpdateHandler(w, r)
	testEqual(t, "Response code = %+v, want %+v", w.Code, 422)

	// Create
	r = httptest.NewRequest("GET", "/api/records/unpaused.test", nil)
	rctx = chi.NewRouteContext()
	rctx.URLParams.Set("key", "unpaused.test")
	w = httptest.NewRecorder()
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
	apiRecordsUpdateHandler(w, r)
	testEqual(t, "Response code = %+v, want %+v", w.Code, 200)
	testEqual(t, "Content-Type header = %+v, want %+v", w.Header().Get("Content-Type"), "application/json; charset=utf-8")
	testEqual(t, "Body = %+v, want %+v", w.Body.String(), "{\"data\":{\"unpaused.test\":{\"paused\":false}}}\n")
	// verify record created in db
	rec, _ := db.get("unpaused.test")
	testEqual(t, "get() = %+v, want %+v", *rec, Record{Paused: false})

	// Create paused
	r = httptest.NewRequest("GET", "/api/records/paused.test", strings.NewReader("{\"paused\":true}"))
	rctx = chi.NewRouteContext()
	rctx.URLParams.Set("key", "paused.test")
	w = httptest.NewRecorder()
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
	apiRecordsUpdateHandler(w, r)
	testEqual(t, "Response code = %+v, want %+v", w.Code, 200)
	testEqual(t, "Content-Type header = %+v, want %+v", w.Header().Get("Content-Type"), "application/json; charset=utf-8")
	testEqual(t, "Body = %+v, want %+v", w.Body.String(), "{\"data\":{\"paused.test\":{\"paused\":true}}}\n")
	// verify record created in db
	rec, _ = db.get("paused.test")
	testEqual(t, "get() = %+v, want %+v", *rec, Record{Paused: true})

	// Update
	r = httptest.NewRequest("GET", "/api/records/unpaused.test", strings.NewReader("{\"paused\":true}"))
	rctx = chi.NewRouteContext()
	rctx.URLParams.Set("key", "unpaused.test")
	w = httptest.NewRecorder()
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
	apiRecordsUpdateHandler(w, r)
	testEqual(t, "Response code = %+v, want %+v", w.Code, 200)
	testEqual(t, "Body = %+v, want %+v", w.Body.String(), "{\"data\":{\"unpaused.test\":{\"paused\":true}}}\n")
	// verify record updated in db
	rec, _ = db.get("unpaused.test")
	testEqual(t, "get() = %+v, want %+v", *rec, Record{Paused: true})
}

func Test_apiRecordsDeleteHandler(t *testing.T) {
	db.Reset()
	if err := db.put("test.test", &Record{}); err != nil {
		t.Errorf("failed to put: %+v", err)
	}

	r := httptest.NewRequest("DELETE", "/api/records/test.test", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Set("key", "test.test")
	w := httptest.NewRecorder()
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
	apiRecordsDeleteHandler(w, r)
	testEqual(t, "Response code = %+v, want %+v", w.Code, 204)
	// verify record deleted in db
	_, err := db.get("test.test")
	testEqual(t, "get() err = %+v, want %+v", err, errRecordNotFound)
}

func Test_apiSettingsUpdateHandler(t *testing.T) {
	db.Reset()

	// Disable
	testEqual(t, "isDisabled = %+v, want %+v", isDisabled, false)
	r := httptest.NewRequest("GET", "/api/settings/", strings.NewReader("{\"disabled\":true}"))
	w := httptest.NewRecorder()
	apiSettingsUpdateHandler(w, r)
	testEqual(t, "Response code = %+v, want %+v", w.Code, 200)
	testEqual(t, "Body = %+v, want %+v", w.Body.String(), "{\"data\":{\"disabled\":true}}\n")
	testEqual(t, "isDisabled = %+v, want %+v", isDisabled, true)

	// Enable
	testEqual(t, "isDisabled = %+v, want %+v", isDisabled, true)
	r = httptest.NewRequest("GET", "/api/settings/", strings.NewReader("{\"disabled\":false}"))
	w = httptest.NewRecorder()
	apiSettingsUpdateHandler(w, r)
	testEqual(t, "Response code = %+v, want %+v", w.Code, 200)
	testEqual(t, "Body = %+v, want %+v", w.Body.String(), "{\"data\":{\"disabled\":false}}\n")
	testEqual(t, "isDisabled = %+v, want %+v", isDisabled, false)
}

func Test_cssHandler(t *testing.T) {
	r := httptest.NewRequest("GET", "/css/nogo.css", nil)
	w := httptest.NewRecorder()
	cssHandler(w, r)
	testEqual(t, "Response code = %+v, want %+v", w.Code, 200)
	testEqual(t, "Cache-Control header = %+v, want %+v", w.Header().Get("Cache-Control"), "public, max-age=31536000")
	testEqual(t, "Content-Type header = %+v, want %+v", w.Header().Get("Content-Type"), "text/css")
	testEqual(t, "Body contains 'Roboto' = %+v, want %+v", strings.Contains(w.Body.String(), "Roboto"), true)
}
