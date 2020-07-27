package router

import (
	"github.com/coda-it/gowebserver/session"
	"github.com/coda-it/gowebserver/store"
	"net/http"
	"testing"
)

func handlerCallback(w http.ResponseWriter, r *http.Request, opt URLOptions, sm session.ISessionManager, s store.IStore) {
}

func TestNew(t *testing.T) {
	t.Run("Should initialize session manager", func(t *testing.T) {
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
}
