package session

type ISessionManager interface {
    IsExist(string) bool
    Create(string)  Session
}

type SessionManager struct {
    sessions map[string]Session
}

func New() SessionManager {
   return SessionManager{
       make(map[string]Session),
   }
}

func (s SessionManager) Create(sessionId string) Session {
    if persistedSession, ok := s.sessions[sessionId]; ok {
        return persistedSession
    }

    session := Session {
        Variables: make(map[string]interface{}),
    }

    s.sessions[sessionId] = session

    return session
}

func (s SessionManager) Get(sid string) Session {
    return s.sessions[sid]
}

func (s SessionManager) IsExist(sessionId string) bool {
    if _, ok := s.sessions[sessionId]; ok {
        return true
    }

    return false
}

