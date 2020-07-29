package router

import (
	"bytes"
	"encoding/json"
	"github.com/coda-it/gowebserver/session"
	"github.com/coda-it/gowebserver/store"
	"net/http"
	"net/http/httptest"
	"testing"
)

func handlerCallback(w http.ResponseWriter, r *http.Request, opt URLOptions, sm session.ISessionManager, s store.IStore) {
}

func contentHandler(t *testing.T, content string) ControllerHandler {
	return func(w http.ResponseWriter, r *http.Request, opt URLOptions, sm session.ISessionManager, s store.IStore) {
		_, err := w.Write([]byte(content))

		if err != nil {
			t.Error(err)
		}
	}
}

func TestNew(t *testing.T) {
	t.Run("Should add route", func(t *testing.T) {
		sm := session.New()

		router := New(sm, handlerCallback)

		if len(router.urlRoutes) != 0 {
			t.Errorf("Router should have no routes")
		}

		router.AddRoute("/api/user", "ALL", handlerCallback)

		addedRoute := router.urlRoutes[0]

		if len(router.urlRoutes) != 1 && addedRoute.method == "ALL" && addedRoute.urlRegExp == "^\\/api\\/user$" {
			t.Errorf("Router should have one route after adding one")
		}
	})

	t.Run("Should handle route with ALL method", func(t *testing.T) {
		sm := session.New()

		router := New(sm, handlerCallback)

		if len(router.urlRoutes) != 0 {
			t.Errorf("Router should have no routes")
		}

		content := "handler executed"

		router.AddRoute("/api/user", "ALL", contentHandler(t, content))

		jsonBytes, _ := json.Marshal(struct{}{})
		request, _ := http.NewRequest(http.MethodPost, "/api/user", bytes.NewReader(jsonBytes))
		writer := httptest.NewRecorder()

		router.Route(writer, request)

		if bytes.Compare(writer.Body.Bytes(), []byte(content)) != 0 {
			t.Errorf("Route not handled")
		}
	})

	t.Run("Should handle route with GET method", func(t *testing.T) {
		sm := session.New()

		router := New(sm, handlerCallback)

		if len(router.urlRoutes) != 0 {
			t.Errorf("Router should have no routes")
		}

		content := "handler executed"

		router.AddRoute("/api/user", "GET", contentHandler(t, content))

		jsonBytes, _ := json.Marshal(struct{}{})
		request, _ := http.NewRequest(http.MethodGet, "/api/user", bytes.NewReader(jsonBytes))
		writer := httptest.NewRecorder()

		router.Route(writer, request)

		if bytes.Compare(writer.Body.Bytes(), []byte(content)) != 0 {
			t.Errorf("Route not handled")
		}
	})

	t.Run("Should handle route with POST method", func(t *testing.T) {
		sm := session.New()

		router := New(sm, handlerCallback)

		if len(router.urlRoutes) != 0 {
			t.Errorf("Router should have no routes")
		}

		content := "handler executed"

		router.AddRoute("/api/user", "POST", contentHandler(t, content))

		jsonBytes, _ := json.Marshal(struct{}{})
		request, _ := http.NewRequest(http.MethodPost, "/api/user", bytes.NewReader(jsonBytes))
		writer := httptest.NewRecorder()

		router.Route(writer, request)

		if bytes.Compare(writer.Body.Bytes(), []byte(content)) != 0 {
			t.Errorf("Route not handled")
		}
	})

	t.Run("Should not handle added route when HTTP methods are different", func(t *testing.T) {
		sm := session.New()

		router := New(sm, handlerCallback)

		if len(router.urlRoutes) != 0 {
			t.Errorf("Router should have no routes")
		}

		content := "handler executed"

		router.AddRoute("/api/user", "GET", contentHandler(t, content))

		jsonBytes, _ := json.Marshal(struct{}{})
		request, _ := http.NewRequest(http.MethodPost, "/api/user", bytes.NewReader(jsonBytes))
		writer := httptest.NewRecorder()

		router.Route(writer, request)

		if bytes.Compare(writer.Body.Bytes(), []byte(content)) == 0 {
			t.Errorf("Route shouldn't be handled")
		}
	})
}
