package render_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/libchaos/litchi/render"

	"github.com/go-chi/chi/v5"
)

func Stream(w http.ResponseWriter, r *http.Request) {
	streamContent := strings.NewReader("streaming content")

	render.Stream(r, streamContent).Write(w)
}

func Example_streaming() {
	r := chi.NewRouter()
	r.Get("/stream", Stream)

	res := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/stream", nil)
	r.ServeHTTP(res, req)

	fmt.Println(res.Body)
	fmt.Println(res.Header())

	res = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/stream", nil)
	req.Header.Set("If-Match", "tag")
	r.ServeHTTP(res, req)

	fmt.Println(res.Body)

	// Output:
	// streaming content
	// map[Accept-Ranges:[bytes] Content-Length:[17] Content-Type:[text/plain; charset=utf-8]]
	//
}
