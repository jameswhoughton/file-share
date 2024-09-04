package web

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSetMessageAddsCookieWithMessageToResponse(t *testing.T) {
	response := httptest.NewRecorder()

	cookieValue := "A message about something"

	setMessage(response, "test-key", cookieValue)

	cookies := response.Result().Cookies()

	if len(cookies) != 1 {
		t.Fatalf("expected 1 cookie, got %d", len(cookies))
	}

	if cookies[0].Value != cookieValue {
		t.Fatalf("expected %s got %s", cookieValue, cookies[0].Value)
	}
}

func TestGetMessageRetrievesMessageAndRemovesCookie(t *testing.T) {
	request := httptest.NewRequest("GET", "/", nil)
	response := httptest.NewRecorder()

	cookieKey := "test-key"
	cookieValue := "A message about something"

	request.AddCookie(&http.Cookie{
		Name:  cookieKey,
		Value: cookieValue,
	})

	message, err := getMessage(response, request, cookieKey)

	if err != nil {
		t.Fatal(err)
	}

	if message != cookieValue {
		t.Fatalf("message: %s, doesn't match expected value: %s", message, cookieValue)
	}

	cookies := response.Result().Cookies()

	if len(cookies) != 1 {
		t.Fatalf("expected 1 cookie, got %d", len(cookies))
	}

	if cookies[0].MaxAge != -1 {
		t.Fatalf("Incorrect max age %d", cookies[0].MaxAge)
	}
}
