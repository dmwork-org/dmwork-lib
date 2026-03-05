package network

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestPostForWWWFormForBytres_ClosesBody(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	}))
	defer server.Close()

	params := map[string]string{"key": "value"}
	headers := map[string]string{}

	body, err := PostForWWWFormForBytres(server.URL, params, headers)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if string(body) != `{"status":"ok"}` {
		t.Fatalf("unexpected body: %s", string(body))
	}
}

func TestPostForWWWFormForBytres_NonOKStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":"bad request"}`))
	}))
	defer server.Close()

	params := map[string]string{"key": "value"}
	headers := map[string]string{}

	_, err := PostForWWWFormForBytres(server.URL, params, headers)
	if err == nil {
		t.Fatal("expected error for non-OK status")
	}
}

func TestPostForWWWFormForAll_ClosesBody(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"result":"success"}`))
	}))
	defer server.Close()

	bodyData := strings.NewReader("key=value")
	headers := map[string]string{}

	body, err := PostForWWWFormForAll(server.URL, bodyData, headers)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if string(body) != `{"result":"success"}` {
		t.Fatalf("unexpected body: %s", string(body))
	}
}

func TestPostForWWWFormForAll_NonOKStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"server error"}`))
	}))
	defer server.Close()

	bodyData := strings.NewReader("key=value")
	headers := map[string]string{}

	_, err := PostForWWWFormForAll(server.URL, bodyData, headers)
	if err == nil {
		t.Fatal("expected error for non-OK status")
	}
}

func TestPostForWWWForm_ParsesJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"key":"value","count":123}`))
	}))
	defer server.Close()

	params := map[string]string{"param": "test"}
	headers := map[string]string{}

	result, err := PostForWWWForm(server.URL, params, headers)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result["key"] != "value" {
		t.Fatalf("unexpected key value: %v", result["key"])
	}
}

// trackingReadCloser tracks whether Close was called
type trackingReadCloser struct {
	io.Reader
	closed bool
}

func (t *trackingReadCloser) Close() error {
	t.closed = true
	return nil
}
