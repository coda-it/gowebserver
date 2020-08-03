package session

import "net/http"

// GetSessionID - get user session ID
func GetSessionID(r *http.Request) (string, error) {
	sessionCookie, err := r.Cookie(SessionKey)

	if err != nil {
		return "", err
	}

	return sessionCookie.Value, nil
}
