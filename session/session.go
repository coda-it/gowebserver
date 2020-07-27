package session

// ISession - session interface
type ISession interface {
	Get(string) interface{}
	Set(string, interface{})
}

// Session - session struct
type Session struct {
	Variables map[string]interface{}
}

// Get - gets session
func (s *Session) Get(key string) interface{} {
	return s.Variables[key]
}

// Set - sets new session
func (s *Session) Set(key string, value interface{}) {
	s.Variables[key] = value
}
