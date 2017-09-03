package session

import (
	"net/http"
)

var sessions map[string]Session

func InitializeSessions() {
    sessions = make(map[string]Session)
}

func GetSessions() map[string]Session {
    return sessions
}

func CreateSession(sessionId string) Session {
    if persistedSession, ok := sessions[sessionId]; ok {
        return persistedSession
    }

    session := Session {
        Variables: make(map[string]interface{}),
    }

    sessions[sessionId] = session

    return session
}

func IsLogged(r *http.Request) bool {
	sessionCookie, err := r.Cookie("sid")

    if err == http.ErrNoCookie {
        return false
    } else if err != nil {
        return false
    }

    if _, ok := sessions[sessionCookie.Value]; ok {
        return true
    }

	return false
}

