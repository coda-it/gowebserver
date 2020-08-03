package session

import (
	"net/http"
	"time"
)

// GetSessionID - get user session ID
func GetSessionID(r *http.Request) (string, error) {
	sessionCookie, err := r.Cookie(SessionKey)

	if err != nil {
		return "", err
	}

	return sessionCookie.Value, nil
}

// ClearSession - remove session cookie
func ClearSession(w http.ResponseWriter) {
	cookie := http.Cookie{
		Path:    "/",
		Name:    SessionKey,
		Expires: time.Now().Add(-100 * time.Hour),
		MaxAge:  -1}
	http.SetCookie(w, &cookie)
}
