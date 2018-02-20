package session

var sessions map[string]Session

func Init() {
    sessions = make(map[string]Session)
}

func Create(sessionId string) Session {
    if persistedSession, ok := sessions[sessionId]; ok {
        return persistedSession
    }

    session := SessionFactory()

    sessions[sessionId] = session

    return session
}

func IsExist(sessionId string) bool {
    if _, ok := sessions[sessionId]; ok {
        return true
    }

	return false
}

