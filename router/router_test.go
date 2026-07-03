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

func checkerHandler(r *http.Request) bool {
	return true
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

		router := New(sm, handlerCallback, "/login")

		if len(router.urlRoutes) != 0 {
			t.Errorf("Router should have no routes")
		}

		router.AddRoute("/api/user", "ALL", false, handlerCallback, checkerHandler)

		addedRoute := router.urlRoutes[0]

		if len(router.urlRoutes) != 1 && addedRoute.method == "ALL" && addedRoute.urlRegExp == "^\\/api\\/user$" {
			t.Errorf("Router should have one route after adding one")
		}
	})

	t.Run("Should handle route with ALL method", func(t *testing.T) {
		sm := session.New()

		router := New(sm, handlerCallback, "/login")

		if len(router.urlRoutes) != 0 {
			t.Errorf("Router should have no routes")
		}

		content := "handler executed"

		router.AddRoute("/api/user", "ALL", false, contentHandler(t, content), checkerHandler)

		jsonBytes, _ := json.Marshal(struct{}{})

		requestGet, _ := http.NewRequest(http.MethodGet, "/api/user", bytes.NewReader(jsonBytes))
		writerGet := httptest.NewRecorder()
		router.Route(writerGet, requestGet)
		if bytes.Compare(writerGet.Body.Bytes(), []byte(content)) != 0 {
			t.Errorf("Method GET not handled")
		}

		requestPost, _ := http.NewRequest(http.MethodPost, "/api/user", bytes.NewReader(jsonBytes))
		writerPost := httptest.NewRecorder()
		router.Route(writerPost, requestPost)
		if bytes.Compare(writerPost.Body.Bytes(), []byte(content)) != 0 {
			t.Errorf("Method POST not handled")
		}

		requestDelete, _ := http.NewRequest(http.MethodDelete, "/api/user", bytes.NewReader(jsonBytes))
		writerDelete := httptest.NewRecorder()
		router.Route(writerDelete, requestDelete)
		if bytes.Compare(writerDelete.Body.Bytes(), []byte(content)) != 0 {
			t.Errorf("Method DELETE not handled")
		}

		requestPatch, _ := http.NewRequest(http.MethodPatch, "/api/user", bytes.NewReader(jsonBytes))
		writerPatch := httptest.NewRecorder()
		router.Route(writerPatch, requestPatch)
		if bytes.Compare(writerDelete.Body.Bytes(), []byte(content)) != 0 {
			t.Errorf("Method PATCH not handled")
		}
	})

	t.Run("Should handle only request with GET method", func(t *testing.T) {
		sm := session.New()

		router := New(sm, handlerCallback, "/login")

		if len(router.urlRoutes) != 0 {
			t.Errorf("Router should have no routes")
		}

		content := "handler executed"

		router.AddRoute("/api/user", "GET", false, contentHandler(t, content), checkerHandler)

		jsonBytes, _ := json.Marshal(struct{}{})

		request, _ := http.NewRequest(http.MethodGet, "/api/user", bytes.NewReader(jsonBytes))
		writer := httptest.NewRecorder()
		router.Route(writer, request)
		if bytes.Compare(writer.Body.Bytes(), []byte(content)) != 0 {
			t.Errorf("Method GET not handled")
		}

		requestOther, _ := http.NewRequest(http.MethodPost, "/api/user", bytes.NewReader(jsonBytes))
		writerOther := httptest.NewRecorder()
		router.Route(writerOther, requestOther)
		if writer.Body == nil {
			t.Errorf("Method POST should not handled")
		}
	})

	t.Run("Should handle only request with POST method", func(t *testing.T) {
		sm := session.New()

		router := New(sm, handlerCallback, "/login")

		if len(router.urlRoutes) != 0 {
			t.Errorf("Router should have no routes")
		}

		content := "handler executed"

		router.AddRoute("/api/user", "POST", false, contentHandler(t, content), checkerHandler)

		jsonBytes, _ := json.Marshal(struct{}{})

		request, _ := http.NewRequest(http.MethodPost, "/api/user", bytes.NewReader(jsonBytes))
		writer := httptest.NewRecorder()
		router.Route(writer, request)
		if bytes.Compare(writer.Body.Bytes(), []byte(content)) != 0 {
			t.Errorf("Route not handled")
		}

		requestOther, _ := http.NewRequest(http.MethodGet, "/api/user", bytes.NewReader(jsonBytes))
		writerOther := httptest.NewRecorder()
		router.Route(writerOther, requestOther)
		if writer.Body == nil {
			t.Errorf("Method GET should not handled")
		}
	})

	t.Run("Should redirect to the fallback page when session doesn't exist", func(t *testing.T) {
		sm := session.New()

		router := New(sm, handlerCallback, "/login")

		if len(router.urlRoutes) != 0 {
			t.Errorf("Router should have no routes")
		}

		content := "handler executed"

		router.AddRoute("/user", "GET", true, contentHandler(t, content), checkerHandler)

		jsonBytes, _ := json.Marshal(struct{}{})

		request, _ := http.NewRequest(http.MethodGet, "/user", bytes.NewReader(jsonBytes))
		writer := httptest.NewRecorder()
		router.Route(writer, request)

		if writer.Code != 303 && writer.Header().Get("Location") != "/login" {
			t.Errorf("Route not handled")
		}
	})

	t.Run("Should not execute handler after redirect when session is invalid", func(t *testing.T) {
		sm := session.New()

		router := New(sm, handlerCallback, "/login")

		handlerExecuted := false
		protectedHandler := func(w http.ResponseWriter, r *http.Request, opt URLOptions, sm session.ISessionManager, s store.IStore) {
			handlerExecuted = true
		}

		router.AddRoute("/protected", "GET", true, protectedHandler, checkerHandler)

		jsonBytes, _ := json.Marshal(struct{}{})
		request, _ := http.NewRequest(http.MethodGet, "/protected", bytes.NewReader(jsonBytes))
		writer := httptest.NewRecorder()
		router.Route(writer, request)

		if writer.Code != http.StatusSeeOther {
			t.Errorf("Expected redirect status %d, got %d", http.StatusSeeOther, writer.Code)
		}
		if handlerExecuted {
			t.Errorf("Protected handler must not execute after redirect when session is invalid")
		}
	})

	t.Run("Should execute handler when valid session exists", func(t *testing.T) {
		sm := session.New()
		sm.Create("valid-session-id")

		router := New(sm, handlerCallback, "/login")

		content := "protected content"
		router.AddRoute("/protected", "GET", true, contentHandler(t, content), checkerHandler)

		jsonBytes, _ := json.Marshal(struct{}{})
		request, _ := http.NewRequest(http.MethodGet, "/protected", bytes.NewReader(jsonBytes))
		request.AddCookie(&http.Cookie{Name: session.SessionKey, Value: "valid-session-id"})
		writer := httptest.NewRecorder()
		router.Route(writer, request)

		if writer.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, writer.Code)
		}
		if bytes.Compare(writer.Body.Bytes(), []byte(content)) != 0 {
			t.Errorf("Expected handler to execute with valid session, got body: %s", writer.Body.String())
		}
	})

	t.Run("Should replace routes atomically", func(t *testing.T) {
		sm := session.New()
		notFoundContent := "not found handler"
		router := New(sm, contentHandler(t, notFoundContent), "/login")

		oldContent := "old handler"
		newContent := "new handler"

		router.AddRoute("/old", "GET", false, contentHandler(t, oldContent), checkerHandler)

		router.ReplaceRoutes(func(addRoute AddRouteFunc) {
			addRoute("/new", "GET", false, contentHandler(t, newContent), checkerHandler)
		})

		requestNew, _ := http.NewRequest(http.MethodGet, "/new", nil)
		writerNew := httptest.NewRecorder()
		router.Route(writerNew, requestNew)
		if bytes.Compare(writerNew.Body.Bytes(), []byte(newContent)) != 0 {
			t.Errorf("Route added via ReplaceRoutes not handled")
		}

		requestOld, _ := http.NewRequest(http.MethodGet, "/old", nil)
		writerOld := httptest.NewRecorder()
		router.Route(writerOld, requestOld)
		if bytes.Compare(writerOld.Body.Bytes(), []byte(notFoundContent)) != 0 {
			t.Errorf("Route registered before ReplaceRoutes should be removed")
		}
	})

	t.Run("Should handle requests while routes are being replaced", func(t *testing.T) {
		sm := session.New()
		router := New(sm, handlerCallback, "/login")

		router.AddRoute("/api/user", "GET", false, handlerCallback, checkerHandler)

		done := make(chan struct{})

		go func() {
			for i := 0; i < 100; i++ {
				router.ReplaceRoutes(func(addRoute AddRouteFunc) {
					addRoute("/api/user", "GET", false, handlerCallback, checkerHandler)
				})
			}
			close(done)
		}()

		for i := 0; i < 100; i++ {
			request, _ := http.NewRequest(http.MethodGet, "/api/user", nil)
			writer := httptest.NewRecorder()
			router.Route(writer, request)
		}

		<-done
	})

	t.Run("Should render to 'not found' page when checker doesn't allow to access route", func(t *testing.T) {
		checkerHandler := func(r *http.Request) bool {
			return false
		}

		sm := session.New()
		notFoundContent := "not found handler"
		router := New(sm, contentHandler(t, notFoundContent), "/login")
		router.AddRoute("/user", "GET", false, contentHandler(t, "user handler"), checkerHandler)

		jsonBytes, _ := json.Marshal(struct{}{})
		request, _ := http.NewRequest(http.MethodGet, "/user", bytes.NewReader(jsonBytes))
		writer := httptest.NewRecorder()

		router.Route(writer, request)

		if bytes.Compare(writer.Body.Bytes(), []byte(notFoundContent)) != 0 {
			t.Errorf("Handler should not be accessed and blocked by checker")
		}
	})
}
