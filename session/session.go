package session

type ISession interface {
	Get(string) interface{}
	Set(string, interface{})
}

type Session struct {
	Variables map[string]interface{}
}

func (s *Session) Get(key string) interface{} {
	return s.Variables[key]
}

func (s *Session) Set(key string, value interface{}) {
	s.Variables[key] = value
}
