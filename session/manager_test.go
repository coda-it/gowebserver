package session

import (
    "testing"
    "reflect"
)

func TestInitializeSessions(t *testing.T) {
    t.Run("Should initialize 'session' singleton", func(t *testing.T) {
        Init()
        expectedSessions := make(map[string]Session)

        isInitialized := reflect.DeepEqual(sessions, expectedSessions)

        if !isInitialized {
            t.Errorf("Sessions array not initialized")
        }
    })
}

func TestCreateSession(t *testing.T) {
    t.Run("Should create session", func(t *testing.T) {
        createdSession := Create("mySessionId")

        isSessionType := reflect.TypeOf(createdSession) == reflect.TypeOf(Session{})

        if !isSessionType {
            t.Errorf("Session object not returned")
        }
    })
}

func TestIsExist(t *testing.T) {
    t.Run("Should return true if user have session cookie which is " +
        "persisted in singleton", func(t *testing.T) {

        isLogged := IsExist("mySessionId")

        if !isLogged {
            t.Errorf("User shouldn be recognised as logged")
        }
    })

    t.Run("Should return false if user doesn't have session cookie",
        func(t *testing.T) {

        isLogged := IsExist("myNotExistingSessionId")

        if isLogged {
            t.Errorf("User shouldn't be recognised as logged")
        }
    })
}
