package session

import (
    "testing"
)

func TestSession(t *testing.T) {
    session :=  &Session{
        make(map[string]interface{}),
    }

    t.Run("Should set and get session variables", func(t *testing.T) {
        session.Set("key1", "value1")
        value1 := session.Get("key1")

        if value1 != "value1" {
            t.Errorf("Variable was not stored properly")
        }

        session.Set("key2", 100)
        value2 := session.Get("key2")

        if value2 != 100 {
            t.Errorf("Variable was not stored properly")
        }
    })
}
