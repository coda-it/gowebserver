package session

import (
    "testing"
    "reflect"
    "net/http/httptest"
    "net/http"
)

func TestInitializeSessions(t *testing.T) {
    t.Run("Should initialize 'session' singleton", func(t *testing.T) {
        InitializeSessions()
        sessions := GetSessions()
        expectedSessions := make(map[string]Session)

        isInitialized := reflect.DeepEqual(sessions, expectedSessions)

        if !isInitialized {
            t.Errorf("Sessions array not initialized")
        }
    })
}

func TestCreateSession(t *testing.T) {
    t.Run("Should create session", func(t *testing.T) {
        createdSession := CreateSession("mySessionId")

        isSessionType := reflect.TypeOf(createdSession) == reflect.TypeOf(Session{})

        if !isSessionType {
            t.Errorf("Session object not returned")
        }
    })
}

func TestIsLogged(t *testing.T) {
    t.Run("Should return true if user have session cookie which is " +
        "persisted in singleton", func(t *testing.T) {

        recorder := httptest.NewRecorder()

        http.SetCookie(recorder, &http.Cookie{Name: "sid", Value: "mySessionId"})
        req := &http.Request{
            Header: http.Header{"Cookie": recorder.HeaderMap["Set-Cookie"]},
        }

        isLogged := IsLogged(req)

        if !isLogged {
            t.Errorf("User shouldn be recognised as logged")
        }
    })

    t.Run("Should return false if user doesn't have session cookie",
        func(t *testing.T) {

        req := &http.Request{}
        isLogged := IsLogged(req)

        if isLogged {
            t.Errorf("User shouldn't be recognised as logged")
        }
    })
}
