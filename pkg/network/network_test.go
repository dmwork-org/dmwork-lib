package network

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestPostForWWWFormForBytres_BodyClose(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	}))
	defer server.Close()

	// Call the function
	result, err := PostForWWWFormForBytres(server.URL, map[string]string{"key": "value"}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if string(result) != `{"status":"ok"}` {
		t.Errorf("unexpected result: %s", result)
	}
}

func TestPostForWWWFormForBytres_NonOKStatus(t *testing.T) {
	// Create a test server that returns 500
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`error`))
	}))
	defer server.Close()

	// Call the function - should return error but still close body
	_, err := PostForWWWFormForBytres(server.URL, map[string]string{"key": "value"}, nil)
	if err == nil {
		t.Fatal("expected error for non-200 status")
	}
}

func TestPostForWWWFormForAll_BodyClose(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	}))
	defer server.Close()

	// Call the function
	result, err := PostForWWWFormForAll(server.URL, strings.NewReader("key=value"), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if string(result) != `{"status":"ok"}` {
		t.Errorf("unexpected result: %s", result)
	}
}

func TestPostForWWWFormForAll_NonOKStatus(t *testing.T) {
	// Create a test server that returns 404
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`not found`))
	}))
	defer server.Close()

	// Call the function - should return error but still close body
	_, err := PostForWWWFormForAll(server.URL, strings.NewReader("key=value"), nil)
	if err == nil {
		t.Fatal("expected error for non-200 status")
	}
}

func TestPostForWWWFormReXML_BodyClose(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`<xml>ok</xml>`))
	}))
	defer server.Close()

	// Call the function
	result, err := PostForWWWFormReXML(server.URL, map[string]string{"key": "value"}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if string(result) != `<xml>ok</xml>` {
		t.Errorf("unexpected result: %s", result)
	}
}

func TestPostForWWWFormForBytres_URLEncoding(t *testing.T) {
	// Test that special characters are properly URL-encoded
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			t.Fatalf("failed to parse form: %v", err)
		}
		// Verify special characters are correctly decoded
		if r.FormValue("name") != "John Doe" {
			t.Errorf("expected 'John Doe', got '%s'", r.FormValue("name"))
		}
		if r.FormValue("query") != "a=b&c=d" {
			t.Errorf("expected 'a=b&c=d', got '%s'", r.FormValue("query"))
		}
		if r.FormValue("special") != "hello+world" {
			t.Errorf("expected 'hello+world', got '%s'", r.FormValue("special"))
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`ok`))
	}))
	defer server.Close()

	params := map[string]string{
		"name":    "John Doe",       // space
		"query":   "a=b&c=d",        // equals and ampersand
		"special": "hello+world",    // plus sign
	}
	result, err := PostForWWWFormForBytres(server.URL, params, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if string(result) != "ok" {
		t.Errorf("unexpected result: %s", result)
	}
}

func TestPostForWWWFormReXML_URLEncoding(t *testing.T) {
	// Test that special characters are properly URL-encoded
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			t.Fatalf("failed to parse form: %v", err)
		}
		// Verify special characters are correctly decoded
		if r.FormValue("data") != "<xml>&value</xml>" {
			t.Errorf("expected '<xml>&value</xml>', got '%s'", r.FormValue("data"))
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`<xml>ok</xml>`))
	}))
	defer server.Close()

	params := map[string]string{
		"data": "<xml>&value</xml>", // XML special characters
	}
	result, err := PostForWWWFormReXML(server.URL, params, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if string(result) != `<xml>ok</xml>` {
		t.Errorf("unexpected result: %s", result)
	}
}
