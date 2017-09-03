package session

type Session struct {
    Variables map[string]interface{}
}

func (session *Session) Get(key string) interface{} {
    return session.Variables[key]
}

func (session *Session) Set(key string, value interface{}) {
    session.Variables[key] = value
}
