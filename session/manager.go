package session

// ISessionManager - session manager interface
type ISessionManager interface {
	IsExist(string) bool
	Create(string) Session
	Get(string) Session
}

// Manager - session manager struct
type Manager struct {
	sessions map[string]Session
}

// New - factory for session manager
func New() Manager {
	return Manager{
		make(map[string]Session),
	}
}

// Create - creates a session
func (s Manager) Create(sessionID string) Session {
	if persistedSession, ok := s.sessions[sessionID]; ok {
		return persistedSession
	}

	session := Session{
		Variables: make(map[string]interface{}),
	}

	s.sessions[sessionID] = session

	return session
}

// Get - gets session
func (s Manager) Get(sid string) Session {
	return s.sessions[sid]
}

// IsExist - checks whether session exists
func (s Manager) IsExist(sessionID string) bool {
	if _, ok := s.sessions[sessionID]; ok {
		return true
	}

	return false
}
