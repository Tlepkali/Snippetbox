package main

import (
	"bytes"
	"net/http"
	"testing"
)

var config = &Config{
	Addr:      "4000",
	StaticDir: "./ui/static",
}

func TestSignupUser(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes(config))
	defer ts.Close()

	_, _, body := ts.get(t, "/user/signup")
	csrfToken := extractCSRFToken(t, body)

	t.Log(csrfToken)
}

func TestShowSnippet(t *testing.T) {
	app := newTestApplication(t)

	ts := newTestServer(t, app.routes(config))
	defer ts.Close()

	tests := []struct {
		name     string
		urlPath  string
		wantCode int
		wantBody []byte
	}{
		{"Valid ID", "/snippets/1", http.StatusOK, []byte("Fake Content")},
		{"Non-existent ID", "/snippets/2", http.StatusNotFound, nil},
		{"Negative ID", "/snippets/-1", http.StatusNotFound, nil},
		{"Decimal ID", "/snippets/1.23", http.StatusNotFound, nil},
		{"String ID", "/snippets/foo", http.StatusNotFound, nil},
		{"Empty ID", "/snippets/", http.StatusNotFound, nil},
		{"Redirect", "/snippets", http.StatusSeeOther, nil},
		{"Trailing slash", "/snippets/1/", http.StatusNotFound, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, _, body := ts.get(t, tt.urlPath)

			if code != tt.wantCode {
				t.Errorf("want %d; got %d", tt.wantCode, code)
			}

			if !bytes.Contains(body, tt.wantBody) {
				t.Errorf("want body to contain %q", tt.wantBody)
			}
		})
	}
}

func TestPing(t *testing.T) {
	app := newTestApplication(t)

	ts := newTestServer(t, app.routes(config))
	defer ts.Close()

	code, _, body := ts.get(t, "/ping")

	if code != http.StatusOK {
		t.Errorf("want %d; got %d", http.StatusOK, code)
	}

	if string(body) != "OK" {
		t.Errorf("want body to equal %q", "OK")
	}
}
