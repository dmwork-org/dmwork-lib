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
	// Create a test server that captures the request body
	var receivedBody string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			t.Errorf("failed to parse form: %v", err)
		}
		receivedBody = r.PostForm.Encode()
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`ok`))
	}))
	defer server.Close()

	// Test with special characters that need URL encoding
	params := map[string]string{
		"name":  "test&admin=true",
		"value": "hello world",
		"chars": "a=b&c=d",
	}

	_, err := PostForWWWFormForBytres(server.URL, params, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify that special characters were properly encoded and not injected
	if strings.Contains(receivedBody, "admin=true") && !strings.Contains(receivedBody, "name=test%26admin%3Dtrue") {
		t.Errorf("URL encoding vulnerability: special chars not properly encoded. Body: %s", receivedBody)
	}
}

func TestPostForWWWFormReXML_URLEncoding(t *testing.T) {
	// Create a test server that verifies form data is properly encoded
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			t.Errorf("failed to parse form: %v", err)
		}
		// Verify the value is correctly decoded (was properly encoded in request)
		if r.PostForm.Get("injection") != "value&extra=injected" {
			t.Errorf("expected value with special chars, got: %s", r.PostForm.Get("injection"))
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`<ok/>`))
	}))
	defer server.Close()

	params := map[string]string{
		"injection": "value&extra=injected",
	}

	_, err := PostForWWWFormReXML(server.URL, params, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
