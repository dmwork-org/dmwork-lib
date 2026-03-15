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
	// Test that special characters are properly URL encoded
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			t.Errorf("failed to parse form: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Verify the value with special characters is correctly parsed
		// If not properly encoded, "test&admin=true" would be split
		name := r.PostForm.Get("name")
		if name != "test&admin=true" {
			t.Errorf("expected name='test&admin=true', got '%s'", name)
		}

		// Verify no extra "admin" parameter was injected
		admin := r.PostForm.Get("admin")
		if admin != "" {
			t.Errorf("unexpected admin parameter injected: '%s'", admin)
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`ok`))
	}))
	defer server.Close()

	// This value contains special characters that could cause injection
	params := map[string]string{
		"name": "test&admin=true",
	}
	_, err := PostForWWWFormForBytres(server.URL, params, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestPostForWWWFormReXML_URLEncoding(t *testing.T) {
	// Test that special characters are properly URL encoded
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			t.Errorf("failed to parse form: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Verify special characters are preserved correctly
		value := r.PostForm.Get("data")
		expected := "a=b&c=d"
		if value != expected {
			t.Errorf("expected data='%s', got '%s'", expected, value)
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`<ok/>`))
	}))
	defer server.Close()

	params := map[string]string{
		"data": "a=b&c=d",
	}
	_, err := PostForWWWFormReXML(server.URL, params, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
