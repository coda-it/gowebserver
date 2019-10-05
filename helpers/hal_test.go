package helpers

import (
	"testing"
)

func TestServeHal(t *testing.T) {
	t.Run("Should transform data into HAL compliant response", func(t *testing.T) {
		data := struct {
			Prop1 string `json:"prop1"`
		}{
			"value1",
		}

		links := struct{}{}
		embedded := struct{}{}

		res := ServeHal(data, links, embedded)

		if _, ok := res["_links"]; !ok {
			t.Errorf("Response doen't contain _links")
		}

		if _, ok := res["_embedded"]; !ok {
			t.Errorf("Response doen't contain _embedded")
		}

		if _, ok := res["prop1"]; !ok {
			t.Errorf("Response doen't contain regular properties")
		}
	})
}
